FROM docker.io/library/golang:1.24 as builder
WORKDIR /ascii-image-generator
COPY ./ ./
RUN cd ./src/go && go mod init ascii-generator && go mod tidy && GOOS=js GOARCH=wasm go build -o main.wasm main.go

FROM node:20-alpine
WORKDIR /ascii-image-generator
COPY ./ ./
RUN npm install

EXPOSE 3000

CMD ["npx", "serve", "src", "-l", "3000"]
