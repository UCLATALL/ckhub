{{/* Returns the fully qualified name of the chart. */}}
{{- define "chart.name" -}}
{{- $chart := include "sanitize" .Chart.Name -}}
{{- $release := include "sanitize" .Release.Name -}}
{{- if contains $chart $release -}}
{{- $release -}}
{{- else -}}
{{- include "sanitize" (cat $chart $release) -}}
{{- end -}}
{{- end -}}

{{/* Returns generic labels of the chart. */}}
{{- define "chart.labels" -}}
app.kubernetes.io/name: {{ include "sanitize" .Chart.Name | quote }}
app.kubernetes.io/instance: {{ include "chart.name" . | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
helm.sh/chart: {{ printf "%s-%s" (include "sanitize" .Chart.Name) .Chart.Version }}
{{- end -}}

{{/* Returns selector labels of the chart. */}}
{{- define "chart.selector" -}}
app.kubernetes.io/name: {{ include "sanitize" .Chart.Name | quote }}
app.kubernetes.io/instance: {{ include "chart.name" . | quote }}
{{- end -}}

{{/**************************************************************************/}}

{{/* Returns the appropriate deployment api version. */}}
{{- define "kubernetes.deployment.version" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.Version -}}
{{- print "extensions/v1beta1" -}}
{{- else -}}
{{- print "apps/v1" -}}
{{- end -}}
{{- end -}}

{{/* Returns the appropriate horizontal pod autoscaler api version. */}}
{{- define "kubernetes.hpa.version" -}}
{{- if semverCompare "<1.23-0" .Capabilities.KubeVersion.Version -}}
{{- print "autoscaling/v2beta1" -}}
{{- else -}}
{{- print "autoscaling/v2" -}}
{{- end -}}
{{- end -}}

{{/* Returns the appropriate ingress api version. */}}
{{- define "kubernetes.ingress.version" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.Version -}}
{{- print "extensions/v1beta1" -}}
{{- else if semverCompare "<1.19-0" .Capabilities.KubeVersion.Version -}}
{{- print "networking.k8s.io/v1beta1" -}}
{{- else -}}
{{- print "networking.k8s.io/v1" -}}
{{- end -}}
{{- end -}}

{{/* Returns the appropriate pod disruption budget api version. */}}
{{- define "kubernetes.pdb.version" -}}
{{- if semverCompare "<1.21-0" .Capabilities.KubeVersion.Version -}}
{{- print "policy/v1beta1" -}}
{{- else -}}
{{- print "policy/v1" -}}
{{- end -}}
{{- end -}}

{{/**************************************************************************/}}

{{/* Returns the container resources. */}}
{{- define "resources" -}}
{{- if or (hasKey . "requests") (hasKey . "limits") -}}
{{- toYaml . -}}
{{- else -}}
requests: {{- toYaml . | nindent 2 }}
limits: {{- toYaml . | nindent 2 }}
{{- end -}}
{{- end -}}

{{/* Sanitizes the given resource name. */}}
{{- define "sanitize" -}}
{{- $name := regexReplaceAll "[[:^alnum:]]" . "-" -}}
{{- regexReplaceAll "-+" $name "-" | lower | trunc 63 | trimAll "-" -}}
{{- end -}}
