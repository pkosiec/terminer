FROM node:11-alpine

WORKDIR /app

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache git openssh

# Install lerna-changelog
RUN npm install -g lerna-changelog@0.8.2

COPY ./generate.sh /app/generate.sh
COPY ./package.json /app/package.json

ENTRYPOINT ["sh", "-c", "/app/generate.sh"]
