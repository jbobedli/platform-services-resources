FROM registry.redhat.io/ubi8/nodejs-16:latest
WORKDIR /app
RUN mkdir carpetaauxiliar
COPY package.json .
RUN npm i
COPY . . 
CMD ["node","start"]