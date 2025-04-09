{{/*
Generates multicluster metadata name by its metadataEndpoint
*/}}
{{- define "mutlicluster.name" -}}
{{- $cluster := . -}}
{{- trimSuffix ".apps.lmru.tech/metadata/" (trimPrefix "https://istio-" $cluster) -}}
{{- end -}}
