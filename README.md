# CKHub

## Quick Start

You need following tools:

- Configured Kubernetes cluster
- [Docker](https://docker.com)
- [Helm](https://helm.sh/)
- [Helmfile](https://github.com/helmfile/helmfile)
- [Golang](https://go.dev/)

Follow next steps to deploy ckhub to your cluster:

1. Build container image.

```bash
make docker
```

2. Publish container image to your registry.

```bash
docker tag \
us-central1-docker.pkg.dev/ckhub-proto1/ckhub/play:unknown-linux-amd64
us-central1-docker.pkg.dev/ckhub-proto1/ckhub/play:latest
```

```bash
docker push us-central1-docker.pkg.dev/ckhub-proto1/ckhub/play:latest
```
   
3. Deploy Helm charts to Kubernetes.

```bash
helmfile sync
```

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
