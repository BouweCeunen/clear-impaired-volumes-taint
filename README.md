# Clear Impaired Volumes Taint

There exists a node taint in Kubernetes, called `NodeWithImpairedVolumes=true:NoSchedule`.

This taint is triggered when volumes are in a pending state for a long time. A closed issue exists on GitHub but is not active (https://github.com/kubernetes/kubernetes/issues/55946). Volume limits exist on AWS nodes (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/volume_limits.html), but the number of attached EBS volumes on my Kubernetes nodes don't even come close to these limits.

This is a deployment which runs in your Kubernetes cluster and removes the taint from your nodes when it occurs. If anyone has more information on why this taint is being set on my nodes, despite no EBS volumes are in an unattached state, feel free to let me know.

## Usage

In order to add this deployment to your cluster, apply the kubernetes folder.
```bash
kubectl apply -f kubernetes/
```

To build executable locally set your GOPATH.
```bash
GOPATH=$(pwd)
go build -o ./clear-impaired-volumes-taint .
```

To manually add taints to your nodes for testing execute following patch to your node.
```bash
kubectl patch node <node> -p '{
        "spec": {
                "taints": [{
                        "effect": "NoSchedule",
                        "key": "NodeWithImpairedVolumes",
                        "value": "true"
                }]
        }
}'
```