# Add ImagePull Secret

Regla para añadir un imagepullsecret a los deploys.

## Descripción

Añade a los deployments que descarguen una imagen de un registry concreto, tipicamente el registry del codigo fuente del proyecto, el secreto necesario para hacer pull.

## Parametros

Lista con un solo término, que coincide con el secreto con las credenciales del registry.

## Recursos

Por defecto actua sobre Pods unicamente

## Filtros

String con un wildcard para indicar el patron del registry. filter: "registryconcreto/*"  

