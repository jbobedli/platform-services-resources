# Allowed Signed Images

Policy para evitar que se desplieguen pods que hacen pull de imagenes sin firmar.
## Descripción
Puedes obligar a que las imagenes que mantienen paquetes de código fuente en repositorios de artefactos privados, estén previamente firmadas con Cosign.
## Parametros
No recibe parametros
## Recursos
Por defecto actua sobre Pods unicamente
## Filtros
El substring del registry de imagenes al que afectar. Por ejemplo: registry.cliente.dominio/projects/preproduccion/*
## Consideraciones :star:
Para utilizar esta Policy debemos tener un secret de nombre "cosign-public-key" en nuestro cluster, con una key "cosign.pub" y el contenido de la clave pública en base64. 

**Podemos crear el fichero manualmente:**
1. Creamos el secret por el medio que escojamos.
2. En el values mencionamos en qué namespace lo hemos creado. Si lo hemos creado en el namespace "desarrollo".
```yaml 
verifyImages:  
  secret_namespace: desarrollo
  enabled: false  #No queremos crear ningun fichero
.
.
#Añadir el resto del fichero values según se necesite
```
   3. Si lo hemos creado en el namespace "kyverno"
```yaml 
verifyImages:  
  enabled: false  #No queremos crear ningun fichero
  # no especificamos la propiedad secret_namespace y se usará por defecto "kyverno".
.
.
#Añadir el resto del fichero values según se necesite
```
**Podemos crear el fichero automaticamente:**
1. Creamos un fichero llamado 'cosignpub.txt' con el contenido de la clave púglica en plano, en la ruta "secrets/allowed-signed_images/" (*Ver ejemplo en la carpeta*).
2. En el values mencionamos en qué namespace lo queremos crear. Si lo queremos crear en el namespace "desarrollo".
```yaml 
verifyImages:  
  secret_namespace: desarrollo
  enabled: true  #Queremos que se cree automáticamente
.
#Añadir el resto del fichero values según se necesite
```
3. Si lo queremos crear en el namespace "kyverno"
```yaml 
verifyImages:  
  enabled: true  #Queremos que se cree automáticamente
  # no especificamos la propiedad secret_namespace y se usará por defecto "kyverno".
.
.
#Añadir el resto del fichero values según se necesite
```
4. Helm creará un secret, en el namespace indicado, con el nombre "cosign-public-key" y lo proporcionará a la Policy "allowed-signed_images" como clave pública para la verificación.
```yaml 
.
.
verifyImages:
    - attestors:
        - count: 1
            entries:
            - keys:
                publicKeys: "k8s://{{ $secret_ns }}/cosign-public-key"
                signatureAlgorithm: sha256
        imageReferences:
.
.
```