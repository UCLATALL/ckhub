{{/* Returns the name of the ingress TLS secret. */}}
{{- define "ingress.secret" -}}
{{- include "sanitize" (cat (include "chart.name" .) "tls") -}}
{{- end -}}
