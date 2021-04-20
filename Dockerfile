FROM node:lts-buster

WORKDIR /code

COPY package.json /code/package.json
COPY package-lock.json /code/package-lock.json

RUN npm install

# OPA
COPY opa_capabilities.json /code/opa_capabilities.json
RUN wget https://github.com/open-policy-agent/opa/releases/download/v0.27.1/opa_linux_amd64
RUN cp opa_linux_amd64 /code/opa
RUN chmod a+x /code/opa

COPY . /code
RUN chmod a+x /code/amf-opa-validator-wrapper

ENTRYPOINT ["/code/amf-opa-validator-wrapper"]