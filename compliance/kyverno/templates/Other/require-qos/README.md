# Require QoS

Obliga a que la request de memory sea igual que el limit

## Descripción

Los Pods de tipo QoS son los que tienen prioridad de alojamiento en los Nodos, por parte del default scheduler de Kubernetes.

## Parametros

No recibe parámetros

## Recursos

Lista de recursos a los que afectará la regla. Se admiten "Job", "StatefulSet", "Cronjob", "Pod", "Deployment" y "Daemonset"

## Filtros

Patron para filtrar a qué namespaces afecta la norma. Por ejemplo ^(des-|pre-).* Para los namespaces que empiezan por des- y pre- 

