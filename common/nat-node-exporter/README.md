# nat-node-exporter
## Table of Contents
1. [Description](#description)
2. [Manual for setup to nat instance](#example)

## Description <a name="description"></a>
**What**: Раскатывает по кластреам ЯО service, который следит за nat кластерами.

**How**: Конфигурация из ArgoCD применяется для каждого кластера индивидуально, в зависимости от названия файла в директории `values/`. Отдельно в каждом кластере конфигуарцию считывает Kubernetes-оператор.

## Manual for setup to nat instance <a name="example"></a>
1. Заходим на nat сервер через яндекс кластер, например "lm-techno.p-yc-techno-ks-master-0" по ключу tfadm.
2. Устанавливаем prometheus-node-exporter
```
apt-get install prometheus-node-exporter
cd /tmp/
wget https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz
tar -xf node_exporter-0.18.1.linux-amd64.tar.gz
systemctl stop prometheus-node-exporter.service
cp node_exporter-0.18.1.linux-amd64/node_exporter /usr/bin/prometheus-node-exporter

cat << EOF > /etc/default/prometheus-node-exporter
ARGS="--collector.diskstats.ignored-devices=^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$ \
--collector.textfile.directory=/var/lib/prometheus/node-exporter \
--web.listen-address=`ip r s | grep eth1 | grep src | awk '{print $9}'`:9100 \
--collector.ntp \
--no-collector.wifi \
--collector.ntp.server-is-local \
--collector.processes \
--collector.filesystem.ignored-mount-points (^/(dev|proc|sys|run)($|/)) \
--collector.netclass.ignored-devices=^(veth.*|lxc.*|[a-f0-9]{15})$ \
--collector.netdev.ignored-devices=^(veth.*|lxc.*|[a-f0-9]{15})$ \
--collector.filesystem.ignored-fs-types ^(autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|fuse\.lxcfs|hugetlbfs|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs|squashfs)$ \
--collector.netstat.fields ^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_.*|Tcp_(ActiveOpens|InSegs|OutSegs|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts))$"
EOF

systemctl restart prometheus-node-exporter.service
systemctl status prometheus-node-exporter.service

```
3. проверяем метрики локально
```
curl `ip r s | grep eth1 | grep src | awk '{print $9}'`:9100/metrics
```
4. проверяем метрики на master `curl ip:9100/metrics`
5. Добавляем в argo в values новый кластер nat и прокатываем его
