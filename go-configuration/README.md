# ArgoCD common Go-configuration

### Как пользоваться

Бинарники лежат в папке в корневой папке и называются argo-linux и argo-mac.
Для своей ОС используйте нужный бинарник, винда идет лесом =)

## Все аргументы:
- -**create**
- -**file**
- -**delete**
- -**backup**

Чтобы запустить скрипт, необходимо:
- -скопировать себе на компьютер (я кидаю в папку /tmp) и заполнить файл из папки example(usecustomcrowdapp: true - если необходимо использовать кастомный апп в crowd, по умолч. используется общий - k8s-argocd).
- -добавить переменную окружения VAULT_ADDR=https://vault.lmru.tech
- -добавить переменные окружения **VAULT_RESTIC_ROLE_ID** и **VAULT_RESTIC_SECRET_ID**, которые можно [найти по ссылке](https://vault.lmru.tech/ui/vault/secrets/service-accounts/show/approles/approle-sa-s3-kubernetes-backup?namespace=it-devops)

После заполнения файла можно запускать скрипт с разными параметрами:

```
# команда *create* создает необходимые файлы с конфигами
# а так же создает Bucket в S3 и делает необходимые записи в Vault
# по умолчанию файл с конфигами будет искаться по пути /tmp/clusterInfo.yaml

argo-mac -create

# или можно указать расположение файла

argo-linux -create -file /path/to/file

# Если вдруг для какого-либо кластера необходимо сделать *Bucket* и запись в *Vault*
# то можно использовать отдельную команду -backup
# в том числе будет создан файл в restic-backups для ArgoCD

argo-mac -backup
argo-mac -backup -file /path/to/file

# Если нужно удалить конфиги какого-то кластера, например создали с неправильным названием
# чтобы не удалять руками, можно запустить аргумент -delete
# ПРИМЕЧАНИЕ: файл все равно нужно законфигурить (/tmp/clusterInfo.yaml)

argo-mac -delete

```
#### ВАЖНО!! Запускать скрипт из корневой папки со скриптом (**go-configuration**), т.к пути в пробной версии захардкожены внутри кода.

В папке *folders.yml* находится описание сущностей, для которых будет составлен файл конфигураций.

```
folders:
  common:
    - upmeter
    - clusters-permissions
    - ingress-nginx-controller
    - kube-blender-exporter
    - prometheus-api-ingress
    - prometheus-remote-write
    - restic-backups
  nodes:
    - instanceclasses
    - node_groups
    - disable_smtp
    - Chart
    - values
  monitoring:
    - api-manager
  istio:
    - federation
  alertmanager:
    - config
    - values
```

***common*** говорит о том, что файлы будут расположены в папке common корневого каталога **argocd**.
Дальше в списке идет перечисление папок (в точности как называются папки), в которых будет создан новый файл с конфигом.

В папке **templates** лежат шаблоны, на основе которых создаются файлы конфигураций.

### сборка под разные системы
Для linux - env GOOS=linux GOARCH=amd64 go build -o argo-linux main.go
Для mac m1 - env GOOS=darwin GOARCH=arm64 go build -o argo-mac-m1 main.go
Для mac - env GOOS=darwin GOARCH=amd64 go build -o argo-mac main.go
