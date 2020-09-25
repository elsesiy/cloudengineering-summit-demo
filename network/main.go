package main

import (
	network "github.com/pulumi/pulumi-azure-nextgen/sdk/go/azure/network/latest"
	resources "github.com/pulumi/pulumi-azure-nextgen/sdk/go/azure/resources/latest"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		resourceGroup, err := resources.NewResourceGroup(ctx, "rg", &resources.ResourceGroupArgs{
			ResourceGroupName: pulumi.String("cloud-engineering-summit"),
			Location:          pulumi.String("EastUS2"),
			Tags: pulumi.StringMap{
				"Env":   pulumi.String("Demo"),
				"Event": pulumi.String("Cloud Engineering Summit"),
			},
		})
		if err != nil {
			return err
		}

		vnet, err := network.NewVirtualNetwork(ctx, "vnet", &network.VirtualNetworkArgs{
			VirtualNetworkName: pulumi.String("demo-vnet"),
			AddressSpace: network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.4.0.0/18"),
				},
			},
			ResourceGroupName: resourceGroup.Name,
			Location:          resourceGroup.Location,
			Tags:              resourceGroup.Tags,
		})
		if err != nil {
			return err
		}

		subnet, err := network.NewSubnet(ctx, "aks", &network.SubnetArgs{
			AddressPrefix:     pulumi.String("10.4.0.0/19"),
			SubnetName:        pulumi.String("aks"),
			ResourceGroupName: resourceGroup.Name,
			VirtualNetworkName: vnet.Name,
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceGroupName", resourceGroup.Name)
		ctx.Export("resourceGroupLocation", resourceGroup.Location)
		ctx.Export("virtualNetwork", vnet.Name)
		ctx.Export("subnetID", subnet.ID())

		return nil
	})
}
