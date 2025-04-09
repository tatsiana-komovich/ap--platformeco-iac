# OpenSearch
## Чеклист по заведению нового кластера OLK

 1. Выбрать канал обновления
 2. Создать в `apps/opensearch/<канал>/values/` каталог с именем кластера
 3. Скопировать в него файлы `curator.yaml` `filebeat.yaml` `logstash.yaml` `tenants.yaml` `vector.yaml` `values.yaml`, скорректировав содержимое под новый кластер (имена индексов для ротации куратором, конфигурацию коллектора, конфигурацию парсера и т.п.)
 4. Обязательно поменять `admin.madison_key` и `admin.restore_madison_key` на уникальные, сгенерированные для этого проекта.

## Перенос OLK из одного канала в другой
0. !!!ВАЖНО!!! проверяем, что все PV, связанные с ns=infra-elklogs, имеют Reclaim Policy = Retain и только после этого переходим к след. пунктам
1. В `lmru--devops--argocd-apps/apps/opensearch` делаем `mv from_channel/values/cluster_name to_channel/values/cluster_name`
2. В `lmru--devops--argocd/argocd-crd/values.yaml` переносим cluster_name из `opensearch-flant-from_channel.clusterlist` в  `opensearch-flant-to_channel.clusterlist`
3. Ждем, ждем, ждем пока argocd удалит старое приложение `from_channel.cluster_name` и создаст новое `to_channel.cluster_name`. Все поды и ресурсы в `ns=infra-elklogs` будут пересозданы
