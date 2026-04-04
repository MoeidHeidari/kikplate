{{/*
Expand the name of the chart.
*/}}
{{- define "kikplate.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "kikplate.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Chart name and version label used by the chart resource selector.
*/}}
{{- define "kikplate.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels applied to every resource.
*/}}
{{- define "kikplate.labels" -}}
helm.sh/chart: {{ include "kikplate.chart" . }}
app.kubernetes.io/name: {{ include "kikplate.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Namespace helper – falls back to .Release.Namespace when not set in values.
*/}}
{{- define "kikplate.namespace" -}}
{{- default .Release.Namespace .Values.namespace }}
{{- end }}
