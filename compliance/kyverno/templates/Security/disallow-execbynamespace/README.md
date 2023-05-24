# Disallow Exec Command by Namespace

Evita que los pods de un namespace puedan ejecutarse manualmente.

## Descripción

No permite abrir terminal con los pods de un namespace concreto. Se puede usar para bloquear esta capacidad de abrir consola con los pods en un namespace entero en entornos productivos. Si se quiere abrir temporalmente se puede poner en modo "audit".

## Parametros

Lista con solamente un término, que coincidirá con el nombre del nampespace.

## Recursos

Por defecto actua sobre Pods unicamente

## Filtros

No recibe filtros