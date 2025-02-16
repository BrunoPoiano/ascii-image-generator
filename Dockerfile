
from docker.io/library/golang:1.23.5 as builder 
workdir /ascii-image-generator
copy ./ ./
run cd ./src/go && go mod init ascii-generator && go mod tidy && GOOS=js GOARCH=wasm go build -o main.wasm main.go

from node:20-alpine
workdir /ascii-image-generator
copy ./ ./
run npm install

expose 3000

cmd ["npm", "run", "start"]

