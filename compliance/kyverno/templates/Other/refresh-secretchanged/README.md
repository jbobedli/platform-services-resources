# Refresh Secret after changes

Reinicia Deploys que usen un secret que ha sido modificado

## Descripción

Esta escuchando cambios en el secret y hace un rollout de los Deploys que usen dicho secret para cargar variables de entorno.
Los secrets tienen que tener la label: 'kyverno.application.type: $filter'
Usar cuidadosamente, no recomendable para produccion.

## Parametros

No recibe parámetros

## Recursos

Por defecto actua sobre Deployments unicamente

## Filtros

El valor de la etiqueta "kyverno.application.type" por ejemplo lo podemos setear a "refresh" y solo se aplicará a estos deployments.

