# NsCleaner Operator

How it works:
You need to add your regex to existed file with cluster name or create new one in values folder. Format:
```
<clusterName>-stage-cluster.yaml
```

Example config:
```
---
nsPatterns:
  - techgate:
    pattern: ^techgate-.*
    nsTTL: 3
    protectedNs:
      - techgate-preprod
```
In this example ns cleaner will check every 6 hours (not changeable) all namespaces and if ns matches the regular expression and its not protected ns, check the latest release in this ns. If the latest release is older than today by 3 days, cleaner adds annotation to this ns with deletition date (today + 1 day). On deletetion day, if there are no new releases in this ns, cleaner removes ns.

[Ns delete motification](https://lemanapro.loop.ru/leroymerlin/channels/lmru-devops-namespace-deletion)
