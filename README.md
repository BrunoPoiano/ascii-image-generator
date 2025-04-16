# ASCII Image Generator with Golang & WebAssembly

Image to ASCII Art Converter is a web-based tool that transforms images into ASCII art directly in the browser. Built with Go (Golang) and WebAssembly (WASM), it runs on the client side, ensuring fast processing without relying on a server. This approach enhances performance and privacy while providing a seamless user experience for generating ASCII art from any image.

live version at [Demo](ascii-image-generator-two.vercel.app/).

[Docker image](https://hub.docker.com/repository/docker/brunopoiano/ascii-img-generator/general)

![Screenshot of the App.](/src/assets/demo.png)

## On the browser
`localhost:3000`

## Deployment
### Using image
```bash
docker run --name ascii-img-gen -p 3000:3000 -d brunopoiano/ascii-img-generator
#or
podman run --name ascii-img-gen -p 3000:3000 -d docker.io/brunopoiano/ascii-img-generator
```

### from source code
```bash
git clone git@github.com:BrunoPoiano/ascii-image-generator.git
cd ascii-image-generator
```
**Docker**
```bash
docker compose up -d
```

**Dev**

Compile the Go code into WebAssembly
```bash
cd src/go
GOOS=js GOARCH=wasm go build -o main.wasm
```
Install dependencies and start the server
```bash
npm i
npm run start
```

## Ref
 - [Effects](https://github.com/BrunoPoiano/imgeffects)

## Dependecies
- Node version: v20.0
- Go version: v1.24
