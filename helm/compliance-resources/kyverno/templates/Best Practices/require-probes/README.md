# Require probes

Requiere que los pods tengan liveness y readiness.

## Descripción
Es buena practica que los pods tengan declaradas las pruebas de readiness y liveness, para permitir a Kubernetes juzgar su estado.

## Parametros
No recibe parámetros

## Recursos
Por defecto actua sobre Pods unicamente

## Filtros
Podemos pasarle una etiqueta para filtrar pods a los que no aplicar la norma. Ejemplo los pods de tipo build de Openshift tienen la etiqueta 'openshift.io/build.name: "*"' y no suelen tener probes porque solo se usan para construcción.

