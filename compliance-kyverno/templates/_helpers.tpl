{{/*
Expand the name of the chart.
*/}}
{{- define "kyverno-hector.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kyverno-hector.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "getNamespace" -}}
{{- $ns := "namespace" -}}
{{- $policy := .policy -}}
{{- $namespace := .policy.namespace -}}
{{- if $namespace }}
{{- printf "\n  %s: %s" $ns $namespace }}

{{- end }}
{{- end }}


{{- define "evaluateNamespace" -}}
    {{- $ns := "namespace" -}}
    {{- $policy := .policy -}}
    {{- $namespace := .policy.namespace -}}
    {{- if $namespace }}
        {{- printf $namespace }}
    {{- else}}
        {{- printf "{{request.namespace}}" | toString }}
    {{- end }}
{{- end }}

{{- define "evaluateNamespace2" -}}
    {{- $ns := "namespace" -}}
    {{- $policy := .policy -}}
    {{- $namespace := .policy.namespace -}}
    {{- $paramns := {{- first $policy.parameters | indent 1 }} }}
    {{- if $namespace }}
        {{- printf $namespace }}
    {{- else}}
        {{- printf $paramns }}
    {{- end }}
{{- end }}


{{- define "whitelistsparser" -}}
    {{- $policy := .policy -}}
    {{/*Creamos lista vacia*/}}
    {{ $lista := list }}
    {{/*Recorremos las ips que vienen en el values.yaml*/}}
    {{- range $ip := $policy.parameters}}
        {{- $quotedip := $ip | quote }}
        {{- $lista = append $lista $quotedip }}
    {{- end}}
    {{- $quotedlist := $lista | join ","}}
          =(externalIPs): [{{ $quotedlist }}]
{{- end }}


{{- define "clusterOrNot" -}}
{{- $policy := .policy -}}
{{- $namespace := .policy.namespace -}}
{{- if $namespace }}
{{- printf " Policy" }}
{{- else }}
{{- printf " ClusterPolicy" }}
{{- end }}
{{- end }}



{{/*
Variable para el nombre de la propiedad namespace de las politicas del values
*/}}
{{- define "ns" -}}
{{- printf "namespace" }}
{{- end}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kyverno-hector.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kyverno-hector.labels" -}}
helm.sh/chart: {{ include "kyverno-hector.chart" . }}
{{ include "kyverno-hector.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kyverno-hector.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kyverno-hector.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kyverno-hector.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kyverno-hector.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
