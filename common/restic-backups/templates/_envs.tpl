{{ define "envs" }}
- name: ETCDCTL_API
  value: "3"
- name: KUBECONFIG
  value: "$KUBECONFIG:/etc/kubernetes/admin.conf"
- name: AWS_ACCESS_KEY_ID
  value: '<AWS_ACCESS_KEY_ID>'
- name: RESTIC_REPOSITORY
  value: '<RESTIC_REPOSITORY>'
- name: AWS_DEFAULT_REGION
  value: ru-central1
- name: RESTIC_PASSWORD
  value: '<RESTIC_PASSWORD>'
- name: AWS_SECRET_ACCESS_KEY
  value: '<AWS_SECRET_ACCESS_KEY>'
- name: FILES_TO_BACKUP
  value: /tmp/files_to_backup.txt
- name: BUCKET
  value: '<AWS_BUCKET>'
- name: RESTIC_EXPORTER_START_CHECK
  value: '0 */1 * * *'
- name: RESTIC_EXPORTER_METRICS_PREFIX
  value: flant_backups_restic
{{ end }}
