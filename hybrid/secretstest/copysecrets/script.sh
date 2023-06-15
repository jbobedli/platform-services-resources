#!/usr/bin/env bash
#ORIGIN_NS=namespace-a
#TARGET_NS=namespace-b namespace-b
echo "alias"
YQ=/tmp/workdir/yq_linux_amd64
PATH=$PATH:/tmp/workdir/
cd /tmp/workdir/
if [[ -n "$ORIGIN_NS" && -n "$TARGET_NS" ]]; then
  for i in $(oc get secret -n $ORIGIN_NS -l sync-to=true -o name); do
    echo $i
    echo "copy to"
    for j in $TARGET_NS; do
      oc get $i -n $ORIGIN_NS -o yaml | \
          $YQ -e 'del(.metadata.ownerReferences)' |\
          $YQ -e ".metadata.namespace=\"$j\"" > currentsecret.yaml
      oc delete -f currentsecret.yaml
      oc create -f currentsecret.yaml
      rm currentsecret.yaml
    done
  done
else
  echo "need ORIGIN_NS and TARGET_NS"
fi
