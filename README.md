# ASCII Image Generator with Golang & WebAssembly

Convert images into ASCII art using a web interface. Built using Go (Golang) and WebAssembly (WASM) to execute in the client's browser, providing a fast and efficient way to generate ASCII art from an image.

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
src/go/main.go
GOOS=js GOARCH=wasm go build -o main.wasm
```
Install dependencies and start the server
```bash
npm i
npm run start
```

## Ref
 - [Bild](https://github.com/anthonynsimon/bild)

## Dependecies
- Node version: v20.0
- Go version: v1.24
