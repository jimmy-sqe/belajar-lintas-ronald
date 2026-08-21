{{- define "preview.namespace" -}}
{{- printf "preview-%s-pr-%s" .Values.repoName (.Values.pr.number | toString) }}
{{- end }}

{{- define "preview.labels" -}}
app.kubernetes.io/managed-by: Helm
app.kubernetes.io/instance: {{ .Release.Name }}
preview/repo: {{ .Values.repoName | quote }}
preview/pr-number: {{ .Values.pr.number | quote }}
{{- end }}

{{- define "preview.backendImage" -}}
{{- printf "%s/%s/backend:%s" .Values.registry .Values.repoName .Values.pr.sha }}
{{- end }}

{{- define "preview.frontendImage" -}}
{{- printf "%s/%s/frontend:%s" .Values.registry .Values.repoName .Values.pr.sha }}
{{- end }}
