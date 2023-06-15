# EventRouter
Pod que recopila los eventos del cluster (oc get events -A) y los escribe por consola. Si tenemos el stack OpenshiftLogging configurado, se enviarán al destino que esté dado de alta en la instancia del ClusterLogForwarder.
https://docs.openshift.com/container-platform/4.12/logging/cluster-logging-eventrouter.html

Hay que ser clusterAdmin.
# Instalación:
## A nivel de cluster  #
* Instalamos en el namespace openshift-logging en cuyo caso se recojen eventos a nivel de cluster. Todos los recursos se crean en el ns 'openshift-logging'.
Usamos este values.yaml.
```
namespaces:
  common: openshift-logging
  eventforward:
  - openshift-logging
.
.
.
```
* Los logs se recojen por el stack OpenshiftLogging al destino que tengamos configurado en el ClusterLogForwarder, en concreto al espacio de nombres "infraestructure".
## A nivel de namespaces específicos #
* Instalamos en los 'namespace' de origen desde el que queremos recojer los eventos
Usamos este values.yaml, por ejemplo para recojer eventos de los ns 'des' y 'pre'.
Algunos recursos como el CR se crearán en el ns común 'openshift-logging'.
```
namespaces:
  common: openshift-logging
  eventforward:
  - des
  - pre
.
.
.
```
* Los logs se recojen por el stack OpenshiftLogging al destino que tengamos configurado en el ClusterLogForwarder, en concreto al espacio de nombres "application".
* Hay que revisar que no haya un filtrado de namespaces de usuario en el ClusterLogForwarder que excluya el namespace del que queremos obtener los eventos, en cuyo caso no llegará nada.
Por ejemplo, siguiendo el ejemplo de antes:
```
apiVersion: "logging.openshift.io/v1"
kind: ClusterLogForwarder
metadata:
  name: instance 
  namespace: openshift-logging 
spec:
  outputs:
   - #rellenar con outputs
  inputs: 
   - name: my-app-logs
     application:
        namespaces:
        - des
        - pre
  pipelines:
   - #rellenar con pipelines
```
O bien sin incluir bloque "inputs", se envian logs de todos los namespaces.
```
apiVersion: "logging.openshift.io/v1"
kind: ClusterLogForwarder
metadata:
  name: instance 
  namespace: openshift-logging 
spec:
  outputs:
   - #rellenar con outputs
  pipelines:
   - #rellenar con pipelines
```