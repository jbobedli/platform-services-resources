# Require Request and Limits

Requiere que los pods tengan declarado unos valores de request y de limit.

## Descripción

Sin importar los valores aportados.

## Parametros

No recibe parámetros

## Recursos

Por defecto actua sobre Pods unicamente

## Filtros

Podemos pasarle una etiqueta para filtrar pods a los que NO aplicar la norma. Ejemplo los pods de tipo build de Openshift tienen la etiqueta 'openshift.io/build.name: "*"' y no suelen tener probes porque solo se usan para construcción.

