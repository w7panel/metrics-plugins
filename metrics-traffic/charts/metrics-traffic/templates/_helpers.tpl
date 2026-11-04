{{- define "metrics-traffic.fullname" -}}metrics-traffic{{- end -}}
{{- define "metrics-traffic.labels" -}}
app.kubernetes.io/name: metrics-traffic
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
