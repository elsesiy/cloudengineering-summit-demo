package main

import (
	"encoding/base64"
	"fmt"
	containerservice "github.com/pulumi/pulumi-azure-nextgen/sdk/go/azure/containerservice/latest"
	"github.com/pulumi/pulumi-tls/sdk/v2/go/tls"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create stack reference to read outputs from the remote stack
		slug := fmt.Sprintf("elsesiy/network-layer/%v", ctx.Stack())
		stackRef, err := pulumi.NewStackReference(ctx, "stackRef", &pulumi.StackReferenceArgs{
			Name: pulumi.String(slug),
		})
		if err != nil {
			return err
		}

		// Read remote stack properties
		resourceGroupName := stackRef.GetStringOutput(pulumi.String("resourceGroupName"))
		resourceGroupLocation := stackRef.GetStringOutput(pulumi.String("resourceGroupLocation"))
		subnetID := stackRef.GetStringOutput(pulumi.String("subnetID"))

		// Generate an SSH key pair for the worker nodes
		sshArgs := tls.PrivateKeyArgs{
			Algorithm: pulumi.String("RSA"),
			RsaBits:   pulumi.Int(4096),
		}
		sshKey, err := tls.NewPrivateKey(ctx, "ssh-key", &sshArgs)
		if err != nil {
			return err
		}

		// Create managed aks cluster attached to the previously created subnet
		kubernetes, err := containerservice.NewManagedCluster(ctx, "cluster",
			&containerservice.ManagedClusterArgs{
				DnsPrefix:         pulumi.String("demo-k8s"),
				ResourceName:      pulumi.String("demo-k8s"),
				Location:          resourceGroupLocation,
				ResourceGroupName: resourceGroupName,
				ApiServerAccessProfile: containerservice.ManagedClusterAPIServerAccessProfileArgs{
					AuthorizedIPRanges: pulumi.StringArray{
						pulumi.String(GetCurrentIP()),
					},
				},
				AgentPoolProfiles: containerservice.ManagedClusterAgentPoolProfileArray{
					&containerservice.ManagedClusterAgentPoolProfileArgs{
						Name:         pulumi.String("agentpool"),
						Mode:         pulumi.String("System"),
						Count:        pulumi.Int(3),
						VmSize:       pulumi.String("Standard_DS2_v2"),
						OsType:       pulumi.String("Linux"),
						VnetSubnetID: subnetID,
					},
				},
				NetworkProfile: containerservice.ContainerServiceNetworkProfileArgs{
					LoadBalancerSku: pulumi.String("standard"),
					NetworkPlugin:   pulumi.String("azure"),
					NetworkPolicy:   pulumi.String("calico"),
				},
				LinuxProfile: &containerservice.ContainerServiceLinuxProfileArgs{
					AdminUsername: pulumi.String("demouser"),
					Ssh: containerservice.ContainerServiceSshConfigurationArgs{
						PublicKeys: containerservice.ContainerServiceSshPublicKeyArray{
							containerservice.ContainerServiceSshPublicKeyArgs{
								KeyData: sshKey.PublicKeyOpenssh,
							},
						},
					},
				},
				Identity: containerservice.ManagedClusterIdentityArgs{
					Type: pulumi.String("SystemAssigned"),
				},
				EnableRBAC:              pulumi.Bool(true),
				EnablePodSecurityPolicy: pulumi.Bool(true),
				KubernetesVersion:       pulumi.String("1.16.13"),
			})
		if err != nil {
			return err
		}

		// Read kube config after cluster creation
		ctx.Export("kubeConfig", pulumi.ToSecret(pulumi.All(kubernetes.Name, resourceGroupName).ApplyString(
			func(args interface{}) (string, error) {
				clusterName := args.([]interface{})[0].(string)
				rgName := args.([]interface{})[1].(string)
				credentials, err := containerservice.ListManagedClusterUserCredentials(ctx, &containerservice.ListManagedClusterUserCredentialsArgs{
					ResourceGroupName: rgName,
					ResourceName:      clusterName,
				})
				if err != nil {
					return "", err
				}
				encoded := credentials.Kubeconfigs[0].Value
				kubeConfig, err := base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					return "", err
				}
				return string(kubeConfig), nil
			})))

		return nil
	})
}
