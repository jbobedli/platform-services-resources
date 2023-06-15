Otros recursos relacionados con compliance especificos de Openshift


basicamente consiste en modificar a nivel de cluster (no valido en ROSA ARO) configuracion de maquina mediante un MachineConfig

kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: worker
  name: worker-custom-registry-trust
spec:
  config:
    ignition:
      config: {}
      security:
        tls: {}
      timeouts: {}
      version: 2.2.0
    networkd: {}
    passwd: {}
    storage:
      files:
        - contents:
            source: data:text/plain;charset=utf-8;base64,${DOCKER_REG}
            verification: {}
          filesystem: root
          mode: 420
          path: /etc/containers/registries.d/docker.io.yaml
        - contents:
            source: data:text/plain;charset=utf-8;base64,${POLICY_CONFIG}
            verification: {}
          filesystem: root
          mode: 420
          path: /etc/containers/policy.json
        - contents:
            source: data:text/plain;charset=utf-8;base64,${SIGNER_KEY}
            verification: {}
          filesystem: root
          mode: 420
          path: /etc/pki/developers/harbor-signer-key.pub
  osImageURL: ""



  Cada elemento de tipo variable es un fichero codificado (hay ejemplos) en base64
  Es importante tener la clave publica del par de claves generadas inicialmente para firmar imagenes con Cosign