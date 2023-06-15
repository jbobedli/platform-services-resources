## Usar una imagen de Red Hat Enterprise Linux como base
FROM registry.access.redhat.com/ubi8/go-toolset:1.18.10 as build
COPY . .

RUN  go build -o my_app

FROM ubi8/ubi-micro

RUN mkdir /.sigstore && mkdir /opt/my_app && chgrp -R 0 /opt/my_app  && chmod -R g=u /opt/my_app

COPY --from=build /opt/app-root/src/my_app /opt/my_app
WORKDIR /opt/my_app
CMD /opt/my_app/my_app