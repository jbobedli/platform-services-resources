
#manual install en cicd

#!/bin/bash
oc project cicd && \
helm repo add oteemocharts https://oteemo.github.io/charts && \
helm install sonarqube oteemocharts/sonarqube -f values.yaml -n cicd
