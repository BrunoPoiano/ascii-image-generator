# ASCII Image Generator with Golang & WebAssembly

Convert images into ASCII art using a web interface. Built using Go (Golang) and WebAssembly (WASM) to execute in the client's browser, providing a fast and efficient way to generate ASCII art from an image.

live version at [Demo](ascii-image-generator-two.vercel.app/).

![Screenshot of the App.](/src/assets/demo.png)

## Installation

Both options run on port 3000 | `localhost:3000`

### Container
```bash
git clone git@github.com:BrunoPoiano/ascii-image-generator.git
cd ascii-image-generator
```
```bash
docker compose up -d 
```

### Node
Clone the project

```bash
git clone git@github.com:BrunoPoiano/ascii-image-generator.git
cd ascii-image-generator
```

**Install dependencies**
```bash
npm i
```
**Start the Server**
```bash
npm run start
```

### Go Code Location

```bash
src/go/main.go
```

**To compile the Go code into WebAssembly, run:**
```bash
GOOS=js GOARCH=wasm go build -o main.wasm
```


## Ref
 - [Bild](https://github.com/anthonynsimon/bild)

## Dependecies
**Node version:** v20.0

**Go version:** v1.23.5 


