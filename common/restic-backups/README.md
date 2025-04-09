## Деплой restic для бэкапа etcd в кластерах LM с помощью argocd.

TLDR:
Бэкапим таким образом только конфиги куба, никаких данных.
  
Краткий порядок деплоя в новый кластер:
- создаём бакет и сервисный аккаунт к нему
- заводим секреты в vault
- добавляем values-файл с названием кластера
- пушим в мастер


## Про деплой.
- Деплой осуществляется из [репозитория](https://github.lmru.tech/adeo/lmru--devops--argocd-apps/tree/master/common/restic-backups).
- Контролировать ход деплоя можно из [argo, раздел applications](https://argocd-devops.apps.lmru.tech/applications?proj=&sync=&health=&namespace=&cluster=&labels=&search=restic).
- Список кластеров, куда сейчас осуществляется деплой - в директории values, например values/customer-stage-cluster.yaml. Если нужно внести какие-то кастомные для кластера переменные, вносите их в yaml этого кластера. Общие переменные - в values.yaml.
При необходимости добавить новый кластер для деплоя, следует посмотреть его правильное имя в [argo clusters](https://argocd-devops.apps.lmru.tech/settings/clusters) и назвать файл аналогично этому имени. Если это совсем новый кластер, нужно его добавить ещё и в список [тут](https://github.lmru.tech/adeo/lmru--devops--argocd/blob/ceed328e45cc9afc32a8603c7f82e579a9646562/argocd-crd/values.yaml#L1930)
- Секреты для своего кластера следует разместить в [vault](https://vault.lmru.tech/ui/vault/secrets/restic_backups/list?namespace=flant), namespace flant (namespace можно выбрать перед логином в vault).
- Сборка image осуществляется из [другого репозитория](https://github.lmru.tech/adeo/lmru--devops--k8s--backup-restic/tree/flant/build-for-argo). Если необходимо внести изменения в image, соберите его там и принесите в values.yaml новый тег.

## После деплоя. 
Возможно появление двух алертов: 
1. K8sResticBackupsMalfunctioning. Различные ошибки при выполнении бэкапа. Например, неправильные параметры подключения к s3. Иногда при изменении секретов в vault они не подтягиваются в кластер с помощью argo (берутся старые), тогда помогает нажатие кнопки hard refresh в интерфейсе argo. Что у вас за ошибка, можно определить сделав restic check --no-lock

2. K8sResticBackupsNotScheduling. Алерт означает, что не было бэкапов последние 24 часа. В кластере, где бэкапов до этого не делалось, точно будет активен до момента выполнения первого бэкапа (и после этого ещё $RESTIC_EXPORTER_START_CHECK времени)

## Диагностика
- Проверить работоспособность рестика вручную можно из пода restore-restic.
- Командой restic check --no-lock можно проверить корректность используемых кредов и наличие ошибок в бэкапах.
- Командой restic snapshots --no-lock можно посмотреть список сделанных бэкапов.
- Командой restic prune можно воспользоваться, если check выдаёт ошибки.
- Командой restic unlock можно разлочить, если вдруг залочилось.
- Сделать бэкап вручную возможно запуском скрипта из restore-restic: ./restic_backup.sh
- Также из пода restore-restic можно воспользоваться командой aws --endpoint-url=https://storage.yandexcloud.net s3 ls s3://${BUCKET} чтобы без участия рестика посмотреть, что у вас хранится в s3.


