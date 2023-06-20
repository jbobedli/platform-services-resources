CONFIGURACIÓN DE LA PLATAFORMA DEVOPS
--------------------------------------

1. argocd
  - aplicar el operador (gitops-operator.yaml)
  - aplicar el custom resource (argocd-cr.yaml) "cuidado con la variable en las policy (harcodear)"
  - aplicar el proyecto de plataforma (platform-gitops-project.yaml) "añadir el repo y credenciales desde ui para añadirlo a scoped"
2. permisos
  - aplicar el cluster role (argocd-application-controller-cluster-role.yaml)
  - aplicar el cluster role binding (argocd-application-controller-clusterrolebinding.yaml)
3. operadores
  - aplicar el appset de plataforma (plataforma.yaml) "este appset usa el init"
    - en este paso se genera el appset "bootstrap-appset" que lanza los appset de "operators" y "platform-apps"
4. sincronizar appsets
  - desde el ui de argocd sincronizar los appsets creados anteriormente
5. aplicativos
  - Habilitar el project.enabled del values en "helm/argocd" 
    - esto creará para el namespace habilitado el proyecto y el appset (del que pueden derivar más) de los aplicativos

