# Cinder CSI volume provisioner

Deploys a Cinder csi provisioner to your cluster, with the appropriate storageClass.

## How To install
- Enable deployment of storageclasses using `storageClass.enabled`
- Tag the retain or delete class as default class using `storageClass.delete.isDefault` in your value yaml
- Set `storageClass.<reclaim-policy>.allowVolumeExpansion` to `true` or `false`

First add the repo:

    helm repo add cpo https://kubernetes.github.io/cloud-provider-openstack
    helm repo update

If you are using Helm v3:
```
helm upgrade openstack-cinder-csi . --install  --create-namespace --namespace=openstack-cinder-csi
```

