# CKHub

## Quick Start

You need following tools:

- Configured Kubernetes cluster
- [Helm](https://helm.sh/)
- [Helmfile](https://github.com/helmfile/helmfile)

Follow next steps to deploy ckhub to your cluster:

1. Clone repository.

```bash
git clone https://github.com/UCLATALL/ckhub.git
```

2. Deploy Helm charts to Kubernetes.

```bash
helmfile sync
```

## Jupyter Kernels

The [Helm chart](./.helm) contains an [option](./.helm/values.yaml#L136) that
allows configure the available Jupyter kernels.

| Parameter | Description                                                 | Example                           |
| --------- | ----------------------------------------------------------- | --------------------------------- |
| name      | The external name of the kernel.                            | ir                                |
| init      | Path to the bootstrap script, stored in the `.helm` folder. | [scripts/init.R](./.helm/scripts) |
| kernel    | The internal name of the kernel (as it called in Jupyter).  | ir                                |
| min       | The minimum number of the kernel replicas (clusterwide).    | 5                                 |
| max       | The maximum number of the kernel replicas (clusterwide).    | 50                                |

You can change kernels settings in the [helmfile.yaml](./helmfile.yaml#L64).

## Development

The project contains the [Development Container](.devcontainer) configuration
with well-defined tools and its prerequisites. Using this configuration along
with [Visual Studio Code](https://code.visualstudio.com) and its
[Remote Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
extension is a preferred way to set up a local development environment.

## Contributing

We will do our best to keep [main branch](../../tree/main) in good shape,
with tests passing at all times.

If you intend to make breaking changes, we recommend [filling an issue](../../issues).
If you're only fixing a bug, it's fine to submit a merge request right away but
we still recommend to fill an issue detailing what you're fixing. This is helpful
in case we don't accept that specific fix but want to keep track of the issue.
