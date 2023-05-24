<a href="https://kyverno.io" rel="kyverno.io">![logo](readmeimg/Kyverno_Horizontal.png)</a>
# **Kyverno Policy Library**
Repositorio que contiene Policies de Kyverno a modo de librería.
## MENU
* Descripción
* Instalación
  * Prerequisitos
  * Instalación
* Configuración
  * Values
  * Example values
* Features
  * Policy o ClusterPolicy
  * Custom Name
  * Verificación de imagenes firmadas
    * Creación de secret con fichero
    * Reutilización de secret existente
* Inventario
  * Best Practices
    * disallow-latesttag
    * disallow-nodeport
    * require-probes
    * require-requestlimit
  * Images
    * add-imagepullsecret
    * allowed-registries
    * allowed-signedimages
  * Networking
    * add-nstomemberroll
    * disallow-externalips
    * require-loadbalancer
  * Other
    * disallow-deprecatedapis
    * limit-replicacount
    * refresh-secretchanged
    * require-qos
  * Security
    * disallow-execbynamespace
    * disallow-privilegedpod
  * Tekton
    * disallow-manualtaskrun
    * limit-pipelineruns
# Descripción
Se ofrecen una serie de reglas, a modo de template de Helm para ser utilizada como fuente de datos en ArgoCD u otra herramienta de GitOps.
# Instalación
## Prerequisitos
Necesitamos dos prerequisitos básicos:
1. Un servidor de ArgoCD instalado y accesible.
2. Un cluster de Kubernetes con Kyverno 1.6.0 instalado
3. Credenciales necesarias para crear recursos en ese cluster utilizando la instancia de ArgoCD.
## Instalación
Debe crearse una aplicación de ArgoCD con este repositorio Git como origen, y como destino un cluster de Kubernetes (Openshift, AKS, GKC etc...) con Kyverno instalado.
ArgoCD utilizará helm como motor para resolver las plantillas y creará las Policy o ClusterPolicy que el usuario le indique a través del fichero values.yaml
# Configuración :pencil:
## Values
Necesita un fichero values.yaml con los parámetros necesarios para la instanciación de los recursos.
**values.yaml**
```yaml 
verifyImages:  #Valores para la funcionalidad de verificar imagenes firmadas.
  secret_namespace: String  #OPTIONAL Namespace en el que crear el secret con la clave publica de cosign "cosign-public-key". Por defecto "kyverno".
  enabled: Boolean  #Define si queremos usar la funcionalidad de verificar imagenes, crea un secret con la clave publica. Necesario para instanciar allowed-signed_images Policy.
defaults:  #Valores por defecto
  validation: String  #Sistema de validación por defecto (audit || enforce)
  background: Boolean  #Ejecución en segundo plano por defecto<span>
policies:
  - name: String #Nombre de la Policy o ClusterPolicy a dar de alta
    background: Boolean #Sobreescribe defaults.background para esa Policy
    filter: String #Filtrar a qué recursos afecta la policy
    namespace: String #Namespace sobre el que aplica la Policy, si no se indica se creará de tipo ClusterPolicy
    validation: String #Sobreescribe defaults.validation para esa Policy (enforce | audit).
    parameters: 
      - string
```
## Example Values
En la carpeta de cada Policy se ofrece un values de ejemplo, con valores reales que pueden tomar los parámetros.
Indicar que cada uno de estos ficheros de ejemplo incluyen solamente la inicialización la de Policy que documentan, ofreciendo el bloque "defaults" y dando de alta una Policy en el bloque "policies".
# Features
## Policy o ClusterPolicy :lock:
Cada Policy de la librería actua como template, pudiéndose instanciarse como Policy o ClusterPolicy. Si en la lista de "policies" indicamos un namespace, ésta se creará como Policy dentro del namespace indicado, si no como ClusterPolicy
**values.yaml**
```yaml
verifyImages: 
  enabled: true
  secret_namespace: kyverno 
defaults:
  validation: enforce
  background: false
policies:
  - name: require-qos
    custom_name: require-qos-in-development
    namespace: des              #---------QoS Policy en el namespace "des"
  - name: add-imagepullsecret   #---------Add ImagePullSecret como ClusterPolicy
    validation: enforce
    parameters:
      - secretoharbor
    filter: "registryconcreto/*"    
```
## Custom Name
Si queremos instanciar varias veces una misma Policy. Por ejemplo:
El caso de "require_qos", queremos que se aplique en los namespaces "des" y "pre".
**values.yaml**
```yaml
verifyImages: 
  enabled: true
  secret_namespace: kyverno 
defaults:
  validation: enforce
  background: false
policies:
  - name: require-qos
    custom_name: require-qos-in-development
    namespace: des
  - name: require-qos
    custom_name: require-qos-in-testing
    namespace: pre
```
## Verificación de imagenes firmadas :white_check_mark:
La Policy a dar de alta es "allowed-signed_images" de la carpeta "Images" de esta librería de politicas de Kyverno.
Ejemplo de values.yaml para dar de alta dicha norma.
***values.yaml***
```yaml
verifyImages: 
  enabled: true #Necesario para que se cree el secreto con la clave publica.
  secret_namespace: desarrollo #OPTIONAL si no se establece se usará el namespace "kyverno" 
defaults:
  validation: enforce
  background: false
policies:
  - name: allowed-signed_images
    validation: audit #Aunque por defecto para el resto de normas se use enforce, para esta puedo establecer audit si lo deseo.
    namespace: desarrollo #La norma se aplicará solamente en este namespace
    parameters:
      - test.harbor.cudc.aws.comunidad.madrid* #Revisamos todas las imagenes de este repositorio de artefactos. Importante el *
```
### Creación de secret con clave publica:
Se tiene que añadir la clave publica a la ruta 'publicKeys/allowed-signed_images/cosignpub.txt' es importante mantener el nombre y extensión del fichero.
Helm creará un secret llamado "cosign-public-key" en el namespace "desarrollo". Luego la Policy utiliza dicho secreto para verificar las firmas.
Aqui vemos cómo queda la Policy y cómo se enlaza con el secret.
***template.yaml***
```yaml
apiVersion: kyverno.io/v1
kind: Policy
metadata:
  namespace: desarrollo        #---------Namespace donde hemos creado la Policy
  name: allowed-signed_images
spec:
  background: false                 #---------Background por defecto
  failurePolicy: Fail
  rules:
    - match:
        any:
              kinds:
                - Pod
      name: check-image
      verifyImages:
        - attestors:
            - count: 1
                - keys:
                    publicKeys: "k8s://desarrollo/cosign-public-key"   #---------Enlazamos el secret. Namespace indicado para crear el secret 'desarrollo'
                    signatureAlgorithm: sha256
          imageReferences:
            - test.harbor.cudc.aws.comunidad.madrid*      #---------Registry proporcionado por parametro
          mutateDigest: true
          verifyDigest: true
  validationFailureAction: enforce            #---------Validación por defecto
  webhookTimeoutSeconds: 30
```
### Reutilización de secret existente:
Si no queremos crear un secret desde este chart de Helm y queremos usar uno creado manualmente, el fichero values anterior varía un poco:
***values.yaml***
```yaml
verifyImages: 
  enabled: false #Indicamos que no se cree un secreto nuevo
  secret_namespace: desarrollo #Indicamos en qué namespace reside el secret creado a mano 
.
.
.
# El resto del bloque es igual que en el ejemplo con creación de secret
```
Por supuesto no es necesario incluir el fichero 'secret/allowed-signed_images/cosignpub.txt' porque no queremos crear ningun secreto.
# Inventario

## Best practices
### disallow-latesttag
* **Descripción**: No permite el uso de la tag :latest
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: disallow-latesttag
    validation: enforce
    namespace: test 
```
### disallow-nodeport
* **Descripción**: No permite el uso nodeport
* **Aplica sobre**: Services
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: disallow-nodeport
    validation: enforce
    namespace: test 
```
### require-probes

* **Descripción**: Exige readiness y liveness
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: require-probes
    validation: enforce
    namespace: test 
```
### require-requestlimit

* **Descripción**: Exige la definición de request y limit
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: require-requestlimit
    validation: enforce
    namespace: test 
```
## Images
### add-imagepullsecret

* **Descripción**: Añade un secreto de pull de imagenes
* **Aplica sobre**: Pods
* **Filtros**:Si
* **Parametros**: Si
```yaml
policies:
  - name: add-imagepullsecret
    validation: enforce
    namespace: test 
    filter: "substringdelregistry/*"
    parameters:
    - nombredelsecreto
```
### allowed-registries

* **Descripción**: Permite usar una lista determinada de registries de imagenes
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: Si
```yaml
policies:
  - name: allowed-registries
    validation: enforce
    namespace: test 
    parameters:
    - substringdelregistry/
    - otroregistry/
```
### allowed-signedimages

* **Descripción**: Verifica la firma de imagenes, de registries concretos.
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: Si
```yaml
policies:
  - name: allowed-signedimages
    validation: enforce
    namespace: test 
    parameters:
    - substringdelregistry/
    - otroregistry/
```
## Networking
### add-nstomemberroll

* **Descripción**: Añade los namespaces nuevos automaticamente al MemberRoll de Istio, filtrando namespaces por un prefijo.
* **Aplica sobre**: Namespaces
* **Filtros**:Si
* **Parametros**: No
```yaml
policies:
  - name: add-nstomemberroll
    filter: ^(des-|pre-).*
    namespace: NO ADMITIDO!!! 
```
### disallow-externalips

* **Descripción**: No permite el uso de ips externas, salvo una whitelist si se le indica.
* **Aplica sobre**: Services
* **Filtros**:No
* **Parametros**: Si
```yaml
policies:
  - name: disallow-externalips
    parameters: 
      - 192.168.1.1
      - 192.168.1.2
```
### disallow-loadbalancer

* **Descripción**: No permite la creación de svc de tipo LoadBalancer
* **Aplica sobre**: Services
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: disallow-latesttag
    validation: enforce
    namespace: test 
```
## Other
### disallow-deprecatedapis

* **Descripción**: No permite el uso apis concretas, previene el uso de apis obsoletas.
* **Aplica sobre**: Resources
* **Filtros**:No
* **Parametros**: Si
```yaml
policies:
  - name: disallow-deprecatedapis
    validation: audit
    parameters:
      - authorization.openshift.io/v1/ClusterRole  
```
### limit-replicacount

* **Descripción**: Establece un numero maximo o minimo de replicas en base a una etiqueta "kyverno.application.type".
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: limit-replicacount
    namespace: test
    parameters:
      - ">3"   #Cantidad minima o máxima de replicas.
    filter: "HA" #Valor de la etiquete kyverno.application.type
```
### refresh-secretchanged

* **Descripción**: Recarga Deploys cuando detecta un cambio en un secreto que está enlazado como ENV y que tienen la etiqueta "kyverno.io/watch: 'true'".
* **Aplica sobre**: Deployments, Pods
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: refresh-secretchanged
    namespace: test 
```
### require-qos

* **Descripción**: Requiere que los limits y los requests sean iguales, en los pods con una etiqueta concreta. "kyverno.application/qos: 'true'"
* **Aplica sobre**: Deployments,DeamonSets,Pods,Statefulsets,Cronjobs,Jobs.
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: require-qos
    namespace: test 
```
## Security
### disallow-execbynamespace

* **Descripción**: No permite ejecutar consola en un pod, en un namespace entero.
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: Si
```yaml
policies:
  - name: disallow-execbynamespace
    custom_name: disallow-exec-on-pre
    parameters: 
      - pre-apps  #ClusterPolicy que afecta a los pods del namespace 'pre-apps'
  - name: disallow-execbynamespace
    namespace: desarrollo #Si es namespaced Policy no hemos de indicar parametros.
    custom_name: disallow-exec-on-desarrollo
```
### disallow-privilegedpod

* **Descripción**: No permite la instanciación de Pods con privilegios
* **Aplica sobre**: Pods
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: disallow-privilegedpod
    namespace: test 
```
## Tekton
### disallow-manualtaskrun

* **Descripción**: No permite la creación manual de taskrun, la pipeline de Tekton deberá instanciarlas.
* **Aplica sobre**: TaskRun
* **Filtros**:No
* **Parametros**: No
```yaml
policies:
  - name: disallow-manualtaskrun
    validation: enforce
    namespace: test 
```
### limit-pipelineruns

* **Descripción**: Limita el numero de ejecuciones simultaneas de Pipelines de Tekton en el namespace en que es creada. Previene de sobrecargas de entornos por continuos despliegues.
* **Aplica sobre**: PipelineRuns
* **Filtros**:No
* **Parametros**: Si
```yaml
policies:
  - name: limit-pipelineruns
    namespace: pre #Solo limita el numero de pipelines en ese namespace.
    validation: enforce
    parameters:
      - 3   #Solo 3 pipelines a la vez.
```
