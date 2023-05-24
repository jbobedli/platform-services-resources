# Limit replicaCount

Regla para limitar el numero de replicas máximo de un Deployment.

## Descripción

Evita que por error se elija un numero de replicas muy grande, u obliga a que las aplicaciones HA tengan un mínimo de replicas.

## Parametros

Lista con un solo término, de tipo string. Por ejemplo: ">3"

## Recursos

Por defecto actua sobre Deployments unicamente

## Filtros

El valor de la etiqueta kyverno.application.type para filtrar los Deployments. Por ejemplo podemos indicar que las aplicaciones que tienen por kyverno.application.type: "HA" han de tener ">3" replicas. En ese caso el valor de filters: "HA"

