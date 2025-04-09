{{/*
Generates image name.
If the image_tag exists use it as the tag,
otherwise use default image tag.
*/}}
{{- define "image.name" -}}
{{- $defaultImageTag := "main-0fa9589" -}}
{{- if .Values.globals.image_tag }}
{{- printf "%s:%s" .Values.image_name .Values.globals.image_tag }}
{{- else -}}
{{- printf "%s:%s" .Values.image_name $defaultImageTag -}}
{{- end -}}
{{- end -}}
