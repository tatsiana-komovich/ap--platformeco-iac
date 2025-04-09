# About modules

## Admission policy
To enable it:
```yaml
admission-policy-engine:
  enabled: true
```

An example of how to add AssignMetadata mutation:
```
mutations:
  metadata:
  - apiGroups: [""]
    kinds: ["Service"]
    scope: Namespaced  # (Namespaced | Cluster)
    type: annotations  # (annotations | labels)
    target:
      yandex.cpi.flant.com/listener-subnet-id: <some_value>
    exclusions:
    - d8-*
    - kube-*
```
