
Esta carpeta contiene dos implementaciones de external providers
Son servicios que responden a una interfaz especifica (external_data) desde una constraint de Opa-Gatekeeper

El caso concreto es de un proveedor al que se le pasa una imagen y se verifica su firma contra un sigstore.
Esta constraint se aplica en el momento de creacion de un pod.

Requiere distintos secrets, clave publica para validar la firma y generacion de claves (server.crt y server.key) para el
sigstore (Es obligado TSL)