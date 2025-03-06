document.addEventListener("DOMContentLoaded", () => {
  if (WebAssembly) {
    const go = new Go();

    WebAssembly.instantiateStreaming(fetch("./go/main.wasm"), go.importObject)
      .then((result) => {
        go.run(result.instance);
      })
      .catch((error) => {
        console.error("Failed to load WebAssembly:", error);
      });
  }

  document.getElementById("download-button").addEventListener("click", () => {
    const canvas = document.getElementById("canvas");
    const image = canvas.toDataURL();
    const downloadLink = document.createElement("a");
    downloadLink.download = "ascii_image.png";
    downloadLink.href = image;
    downloadLink.click();
  });
});

function generateCanva(imgInfo) {
  const canvas = document.getElementById("canvas");
  const fontSizeSpace = 0.6;

  const fontSize = parseInt(imgInfo[0][0].fontSize);
  const canvasWidth = imgInfo[0].length * (fontSize * fontSizeSpace);

  const lineHeight = parseInt(imgInfo[0][0].lineHeight);
  const canvasHeight = imgInfo.length * lineHeight;

  canvas.width = canvasWidth;
  canvas.height = canvasHeight;

  const ctx = canvas.getContext("2d");
  ctx.fillStyle = "black";

  let x = 0;
  let y = 0;
  for (const item of imgInfo) {
    x = 0;
    y += parseInt(item[0].lineHeight);
    for (const line of item) {
      const fontSize = parseInt(line.fontSize);

      ctx.font = `${fontSize}px monospace`;
      ctx.fillStyle = line.color;
      ctx.fillText(line.char, x, y);

      x += fontSize * fontSizeSpace;
    }
  }
}
