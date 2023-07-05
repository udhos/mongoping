# Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

    helm repo add mongoping https://udhos.github.io/mongoping

Update files from repo:

    helm repo update

Search mongoping:

    helm search repo mongoping -l --version ">=0.0.0"
    NAME               	CHART VERSION	APP VERSION	DESCRIPTION
    mongoping/mongoping	0.1.0        	0.1.0      	Install mongoping helm chart into kubernetes.

To install the charts:

    helm install my-mongoping mongoping/mongoping
    #            ^            ^         ^
    #            |            |          \_______ chart
    #            |            |
    #            |             \_________________ repo
    #            |
    #             \______________________________ release (chart instance installed in cluster)

To uninstall the charts:

    helm uninstall my-mongoping

# Source

<https://github.com/udhos/mongoping>
