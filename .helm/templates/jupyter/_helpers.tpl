{{/* Returns the fully qualified container image. */}}
{{- define "jupyter.image" -}}
{{- $tag := default .Chart.AppVersion .Values.jupyter.image.tag -}}
{{- printf "%s/%s:%s" .Values.jupyter.image.registry .Values.jupyter.image.repository $tag -}}
{{- end -}}

{{/* Returns labels of the jupyter. */}}
{{- define "jupyter.labels" -}}
{{- include "chart.labels" . }}
app.kubernetes.io/component: "jupyter"
{{- end -}}

{{/* Returns the name of the jupyter. */}}
{{- define "jupyter.name" -}}
{{- include "sanitize" (cat (include "chart.name" .) "jupyter") -}}
{{- end -}}

{{/* Returns the name of the image pull secret. */}}
{{- define "jupyter.pullSecret" -}}
{{- $name := include "sanitize" (cat (include "jupyter.name" .) "registry") -}}
{{- default $name .Values.jupyter.image.pullSecret.name -}}
{{- end -}}

{{/* Returns selector labels of the jupyter deployment. */}}
{{- define "jupyter.selector" -}}
{{- include "chart.selector" . }}
app.kubernetes.io/component: "jupyter"
{{- end -}}

{{/* Returns the name of the service account. */}}
{{- define "jupyter.serviceAccount" -}}
{{- default (include "jupyter.name" .) .Values.jupyter.serviceAccount.name -}}
{{- end -}}
