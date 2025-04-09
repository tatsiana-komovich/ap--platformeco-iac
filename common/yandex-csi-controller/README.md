# cloud provisioning

```bash
export YANDEX_SA_DATA=$(yc iam service-account create --name k8s-csi-controller --format json )
export YANDEX_SA_ID=$(echo $YANDEX_SA_DATA | jq .id | sed 's\"\\g')
export YANDEX_SA_FOLDER=$(echo $YANDEX_SA_DATA | jq .folder_id | sed 's\"\\g' | tr -d '\n' )
export YANDEX_SA_FOLDER_B64=$(echo -n $YANDEX_SA_FOLDER | base64 )


yc resource-manager folder add-access-binding ${YANDEX_SA_FOLDER} \
  --role admin \
  --subject serviceAccount:${YANDEX_SA_ID}

yc iam key create --service-account-name k8s-csi-controller  --output k8s-csi-controller-key.json

export serviceAccountJSON=$( base64 k8s-csi-controller-key.json | tr -d '\n')

export NAMESPACE=kube-fraima-csi


```bash
helm upgrade yandex-csi-controller . \
--install \
--create-namespace \
--namespace=cloud-provider-yandex \
--set=serviceAccountJSON=$( base64 ../../terraform-p-kafka-ha-yc.json  | tr -d '\n') \
--set=folderID=$(echo -n b1gac70j8rrak03f84ai | base64)
```
