FROM registry.redhat.io/openshift4/ose-cli-alt-rhel8:v4.10

ENV IMAGETAG invalid

RUN mkdir /tmp/workdir  \
      && cd /tmp/workdir  \
      && chgrp -R 0 /tmp/workdir \
      && chmod -R g=u /tmp/workdir \
      && wget https://github.com/mikefarah/yq/releases/download/v4.33.3/yq_linux_amd64 \
      && chmod a+x /tmp/workdir/yq_linux_amd64

COPY script.sh /tmp/workdir/script.sh
RUN chmod a+x /tmp/workdir/script.sh
ENV PATH $PATH:/tmp/workdir
ENV HOME /tmp/workdir
USER 1001


ENTRYPOINT /tmp/workdir/script.sh