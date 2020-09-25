package main

import (
	"context"
	"fmt"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optup"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	stackName = "demo"
	outPath   = "/tmp/kubeconfig"
)

func main() {
	ctx := context.Background()

	fmt.Printf("==> Shipping a New Multi-Cloud, Multi-Region Kubernetes Platform at ❄️\n")
	fmt.Printf("==> Demo by @elsesiy\n\n")

	// Specify dirs to include
	workingDirs := []string{"network", "kubernetes"}
	for i, dir := range workingDirs {
		fmt.Printf("==> Processing directory '%s' (%d/%d)\n", dir, i+1, len(workingDirs))

		// Use the Passphrase secrets provider which allows us to encrypt the state with user-provided password, other
		// supported options: default, awskms, azurekeyvault, gcpkms, hashivault
		opts := auto.SecretsProvider("passphrase")

		// Create or select a stack from the local workspace using dir & stackName as inputs
		stack, err := auto.UpsertStackLocalSource(ctx, stackName, filepath.Join(".", dir), opts)
		if err != nil {
			fmt.Printf("==> Failed to create/select stack '%s': %v\n", stackName, err)
			os.Exit(1)
		}

		fmt.Printf("==> Created/Selected stack '%s'\n", stack.Name())

		// Attach stdout to progress stream
		stdoutStreamer := optup.ProgressStreams(os.Stdout)

		// Run pulumi up programmatically
		res, err := stack.Up(ctx, stdoutStreamer)
		if err != nil {
			fmt.Printf("Failed to deploy stack '%s': %v\n", stack.Name(), err)
			os.Exit(1)
		}

		// Get kube config after cluster creation
		if dir == "kubernetes" {
			kubeConfig, ok := res.Outputs["kubeConfig"].Value.(string)
			if !ok {
				fmt.Println("==> Failed to unmarshall kubeConfig stack output")
				os.Exit(1)
			}

			fmt.Printf("==> Writing kubeconfig to '%s'\n", outPath)
			err := ioutil.WriteFile(outPath, []byte(kubeConfig), 0644)
			if err != nil {
				fmt.Printf("==> Failed to write kubeconfig to '%s': %v", outPath, err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("==> The end 😎")
}
