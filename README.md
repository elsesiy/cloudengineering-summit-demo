# Cloud Engineering Summit (Demo)

This repo hosts the code used during my portion of the talk at the Cloud Engineering Summit
by Pulumi: [Shipping a New Multi-Cloud, Multi-Region Kubernetes Platform at Snowflake](https://cloudengineering.heysummit.com/talks/shipping-a-new-multi-cloud-multi-region-kubernetes-platform-at-snowflake/) :snowflake:

In this demo we're going to deploy a managed Kubernetes Cluster on Azure using
Pulumi's [next-gen Azure provider](https://www.pulumi.com/blog/announcing-nextgen-azure-provider/) and 
the new [Automation API (alpha)](https://github.com/pulumi/pulumi/issues/3901).

The code has been structured in so-called micro-stacks, which is equivalent to microservices only in project and stack form.
In addition, we'll be using a custom secrets provider to ensure we're in control of the encryption key for our sensitive information.
If you don't know what any of this means, 
you can read up on it [here](https://www.pulumi.com/docs/intro/concepts/programming-model/#program-structure).

## Contents

This demo has two projects, the `network-layer` and the `platform-layer`.
The former creates a `VirtualNetwork` & `Subnet` to be used by the latter.
There's some information passed between these two stacks using a `StackReference`.

The `platform-layer` creates a managed Azure Kubernetes Service (AKS) with restricted access 
to the API server (using your current IP) and related resources.

After the resources have been deployed, the stack exports the `kubeConfig` file needed 
to connect to the Kubernetes cluster.

Instead of manually creating the stacks using the `pulumi` CLI, we're using a regular Go program.

### Pre-requisites

- Go 1.15
- Pulumi CLI
- Azure CLI

## Usage

Make sure that...  
...you're logged into the [Pulumi Console](https://app.pulumi.com), if not run `pulumi login`  
...you're logged into your Azure CLI, if not run `az login`

To deploy both stacks, run the following snippet:

```
export PULUMI_CONFIG_PASSPHRASE="<my-secret-passphrase>"
go run main.go
```

After some :clock1:, you'll have your resources deployed and the `kubeconfig` ready to be used located in `/tmp/kubeconfig`.

That's it :sunglasses:

You can now use the `kubeconfig` file to connect to the cluster, e.g.

```
export KUBECONFIG=/tmp/kubeconfig
kubectl get nodes -o wide
```

### Testimonials

Thank you to my friends at Pulumi for giving me the opportunity to speak at their event & 
my employer for using their resources for the purposes of this demo.
