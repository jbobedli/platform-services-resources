
Este template tiene la posibilidad de incluir varias mesh en cada cluster. (Itera por clusterMesh)
Para cada tipo de elemento se introduce un unico yaml y se centraliza todo en el values.
La estructura del values que genera los yaml esta simplificada y abierta a evolucion, por ejemplo:
* Los destinationRules solo tienen habilitado el caso actual que es el de deshabilitar el TLS de istio (por ejemplo con jaeger-collector)
* Mas opciones (Añadir CA, destinationRules de tipo MESHInternal ...) en TODO

Lo mismo para el resto. Los serviceEntry tambien simplifican al maximo y ni siquiera hace falta explicitar el nombre del se (se forma a partir del host)