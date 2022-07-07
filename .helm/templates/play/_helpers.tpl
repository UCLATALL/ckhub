{{/* Returns the fully qualified container image. */}}
{{- define "play.image" -}}
{{- $tag := default .Chart.AppVersion .Values.play.image.tag -}}
{{- printf "%s/%s:%s" .Values.play.image.registry .Values.play.image.repository $tag -}}
{{- end -}}

{{/* Returns labels of the play. */}}
{{- define "play.labels" -}}
{{- include "chart.labels" . }}
app.kubernetes.io/component: "play"
{{- end -}}

{{/* Returns the name of the play. */}}
{{- define "play.name" -}}
{{- include "sanitize" (cat (include "chart.name" .) "play") -}}
{{- end -}}

{{/* Returns the name of the image pull secret. */}}
{{- define "play.pullSecret" -}}
{{- $name := include "sanitize" (cat (include "play.name" .) "registry") -}}
{{- default $name .Values.play.image.pullSecret.name -}}
{{- end -}}

{{/* Returns selector labels of the play deployment. */}}
{{- define "play.selector" -}}
{{- include "chart.selector" . }}
app.kubernetes.io/component: "play"
{{- end -}}

{{/* Returns the name of the service account. */}}
{{- define "play.serviceAccount" -}}
{{- default (include "play.name" .) .Values.play.serviceAccount.name -}}
{{- end -}}
