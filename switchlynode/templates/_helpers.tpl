{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "switchlynode.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "switchlynode.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "switchlynode.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "switchlynode.labels" -}}
helm.sh/chart: {{ include "switchlynode.chart" . }}
{{ include "switchlynode.selectorLabels" . }}
app.kubernetes.io/version: {{ include "switchlynode.tag" . | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/net: {{ include "switchlynode.net" . }}
app.kubernetes.io/type: {{ .Values.type }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "switchlynode.selectorLabels" -}}
app.kubernetes.io/name: {{ include "switchlynode.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "switchlynode.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "switchlynode.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Net
*/}}
{{- define "switchlynode.net" -}}
{{- default .Values.net .Values.global.net -}}
{{- end -}}

{{/*
Tag
*/}}
{{- define "switchlynode.tag" -}}
{{- coalesce  .Values.global.tag .Values.image.tag .Chart.AppVersion -}}
{{- end -}}

{{/*
Image
*/}}
{{- define "switchlynode.image" -}}
{{/* A hash is not needed for mocknet, or in the case that a node is not a validator w/ key material and autoupdate is enabled. */}}
{{- if and .Values.autoupdate.enabled (eq .Values.type "fullnode") -}}
{{- .Values.image.repository -}}:{{ include "switchlynode.tag" . }}
{{- else -}}
{{- .Values.image.repository -}}:{{ include "switchlynode.tag" . }}@sha256:{{ coalesce .Values.global.hash .Values.image.hash }}
{{- end -}}
{{- end -}}

{{/*
RPC Port
*/}}
{{- define "switchlynode.rpc" -}}
{{- if eq (include "switchlynode.net" .) "mainnet" -}}
    {{ .Values.service.port.mainnet.rpc}}
{{- else if eq (include "switchlynode.net" .) "stagenet" -}}
    {{ .Values.service.port.stagenet.rpc}}
{{- else -}}
    {{ .Values.service.port.mainnet.rpc}}
{{- end -}}
{{- end -}}

{{/*
GRPC Port
*/}}
{{- define "switchlynode.grpc" -}}
{{- if eq (include "switchlynode.net" .) "mainnet" -}}
    {{ .Values.service.port.mainnet.grpc}}
{{- else if eq (include "switchlynode.net" .) "stagenet" -}}
    {{ .Values.service.port.stagenet.grpc}}
{{- else -}}
    {{ .Values.service.port.mainnet.grpc}}
{{- end -}}
{{- end -}}

{{/*
P2P Port
*/}}
{{- define "switchlynode.p2p" -}}
{{- if eq (include "switchlynode.net" .) "mainnet" -}}
    {{ .Values.service.port.mainnet.p2p}}
{{- else if eq (include "switchlynode.net" .) "stagenet" -}}
    {{ .Values.service.port.stagenet.p2p}}
{{- else -}}
    {{ .Values.service.port.mainnet.p2p}}
{{- end -}}
{{- end -}}

{{/*
chain id
*/}}
{{- define "switchlynode.chainID" -}}
{{- if eq (include "switchlynode.net" .) "mainnet" -}}
    {{ .Values.chainID.mainnet}}
{{- else if eq (include "switchlynode.net" .) "stagenet" -}}
    {{ .Values.chainID.stagenet}}
{{- else -}}
    {{ .Values.chainID.mainnet}}
{{- end -}}
{{- end -}}