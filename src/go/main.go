package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"strconv"
	"syscall/js"
	"time"

	"github.com/BrunoPoiano/imgeffects/blur"
	"github.com/BrunoPoiano/imgeffects/dithering"
	"github.com/BrunoPoiano/imgeffects/flip"
	"github.com/BrunoPoiano/imgeffects/hsl"
	"github.com/BrunoPoiano/imgeffects/kuwahara"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/noise"
	"github.com/anthonynsimon/bild/segment"
	"github.com/anthonynsimon/bild/transform"
)

type model struct {
	imageWidth      int
	effectSelected  string
	effectRange     float64
	checkAscii      bool
	checkColor      bool
	imageSelected   js.Value
	global          js.Value
	document        js.Value
	effectsRateMap  map[string]bool
	asciiChars      string
	execChangeImage func()
	lineHeight      int
	fontSize        int
}

func main() {
	g := js.Global()
	m := &model{
		asciiChars:     "░▒▓█",
		imageWidth:     100,
		effectRange:    0,
		effectSelected: "",
		checkAscii:     false,
		effectsRateMap: effectsRateMapFunc(),
		global:         g,
		document:       g.Get("document"),
		checkColor:     true,
		lineHeight:     9,
		fontSize:       12,
	}

	//Adding debounce to changeImage func
	m.execChangeImage = debounce(200*time.Millisecond, func() {
		m.changeImage()
	})

	//getting elements
	inputCheckboxAscii := m.document.Call("getElementById", "input-checkbox-ascii")
	inputCheckboxColor := m.document.Call("getElementById", "input-checkbox-color")
	inputEffectRange := m.document.Call("getElementById", "input-effect-range")
	inputTextAscii := m.document.Call("getElementById", "input-text-ascii")
	inputZoomRange := m.document.Call("getElementById", "input-zoom-range")
	selectEffect := m.document.Call("getElementById", "select-effect")
	selectAscii := m.document.Call("getElementById", "select-ascii")
	inputFile := m.document.Call("getElementById", "input-file")

	//adding reactivity
	inputCheckboxAscii.Call("addEventListener", "input", js.FuncOf(m.inputAsciiCheckboxChange))
	inputCheckboxColor.Call("addEventListener", "input", js.FuncOf(m.inputAsciiCheckboxColor))
	inputEffectRange.Call("addEventListener", "input", js.FuncOf(m.inputEffectRangeChange))
	inputTextAscii.Call("addEventListener", "input", js.FuncOf(m.inputTextAsciiChange))
	inputZoomRange.Call("addEventListener", "input", js.FuncOf(m.inputZoomRangeChange))
	selectEffect.Call("addEventListener", "input", js.FuncOf(m.effectChange))
	selectAscii.Call("addEventListener", "input", js.FuncOf(m.selectAsciiChange))
	inputFile.Call("addEventListener", "input", js.FuncOf(m.fileChange))

	//setting default value
	inputCheckboxAscii.Set("checked", m.checkAscii)
	inputCheckboxColor.Set("checked", m.checkColor)
	inputEffectRange.Set("value", m.effectRange)
	inputTextAscii.Set("value", m.asciiChars)
	inputZoomRange.Set("value", m.imageWidth)
	selectEffect.Set("value", m.effectSelected)
	selectAscii.Set("value", m.asciiChars)

	select {}
}

func (m *model) inputAsciiCheckboxColor(this js.Value, args []js.Value) interface{} {
	m.checkColor = this.Get("checked").Bool()
	m.execChangeImage()
	return nil
}

func (m *model) inputAsciiCheckboxChange(this js.Value, args []js.Value) interface{} {
	m.checkAscii = this.Get("checked").Bool()
	m.execChangeImage()
	return nil
}

func (m *model) fileChange(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}

	files := args[0].Get("target").Get("files")

	if files.Length() > 0 {
		m.imageSelected = files.Index(0)
		m.changeImage()
	}
	return nil
}

func (m *model) inputTextAsciiChange(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}
	m.asciiChars = args[0].Get("target").Get("value").String()
	m.execChangeImage()
	return nil
}

func (m *model) selectAsciiChange(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}
	m.asciiChars = args[0].Get("target").Get("value").String()
	m.document.Call("getElementById", "input-text-ascii").Set("value", m.asciiChars)
	m.changeImage()
	return nil
}

func (m *model) effectChange(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}

	m.effectSelected = args[0].Get("target").Get("value").String()
	m.effectRange = 0

	if m.effectsRateMap[m.effectSelected] {
		m.updateEffectRange("0", "10", "1")
	} else {
		switch m.effectSelected {
		case "brightness", "contrast", "saturation":
			m.updateEffectRange("-1", "1", "0.1")
		case "hue":
			m.updateEffectRange("0", "360", "1")
		case "gamma":
			m.updateEffectRange("1", "5", "0.2")
			m.effectRange = 1
		case "threshold":
			m.updateEffectRange("0", "200", "1")
		case "noisePerlin":
			m.updateEffectRange("0", "1", "0.01")
		case "shearH", "shearV":
			m.updateEffectRange("0", "180", "1")
		case "dithering":
			m.updateEffectRange("2", "10", "1")
		case "gaussianBlur":
			m.updateEffectRange("0", "20", "1")
			m.effectRange = 1
		case "kuwahara":
			m.updateEffectRange("1", "20", "2")
			m.effectRange = 1
		default:
			inputRateRangeDiv := m.document.Call("getElementById", "input-rate-range-div")
			changeAttribute(inputRateRangeDiv, "data-visible", "false")
		}
	}

	m.changeImage()
	return nil
}

func (m model) updateEffectRange(min string, max string, step string) {
	inputRange := m.document.Call("getElementById", "input-effect-range")
	inputRange.Set("value", m.effectRange)
	inputRange.Set("min", min)
	inputRange.Set("max", max)
	inputRange.Set("step", step)

	inputRateRangeDiv := m.document.Call("getElementById", "input-rate-range-div")
	changeAttribute(inputRateRangeDiv, "data-min", min)
	changeAttribute(inputRateRangeDiv, "data-max", max)
	changeAttribute(inputRateRangeDiv, "data-visible", "true")

}

func (m *model) inputEffectRangeChange(this js.Value, args []js.Value) interface{} {

	if len(args) == 0 {
		return nil
	}

	value := args[0].Get("target").Get("value").String()
	var err error
	m.effectRange, err = strconv.ParseFloat(value, 64)
	if err != nil {
		println("Error inputEffectRange")
		return nil
	}
	m.execChangeImage()
	return nil
}

func (m *model) inputZoomRangeChange(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}

	value := args[0].Get("target").Get("value").String()
	var err error
	m.imageWidth, err = strconv.Atoi(value)
	if err != nil {
		println("Error inputZoomRange")
		return nil
	}
	m.execChangeImage()
	return nil
}

func (m *model) changeImage() {
	if m.imageSelected.IsUndefined() || m.imageSelected.IsNull() {
		return
	}

	if m.asciiChars == "" || m.asciiChars == " " {
		return
	}

	fileReader := m.global.Get("FileReader").New()

	var onLoad js.Func
	onLoad = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		uint8Array := m.global.Get("Uint8Array").New(this.Get("result"))
		input := make([]byte, uint8Array.Length())
		js.CopyBytesToGo(input, uint8Array)

		var img image.Image = nil
		var err error = nil

		if m.imageSelected.Get("type").String() == "image/jpeg" {
			img, err = jpeg.Decode(bytes.NewReader(input))
		} else {
			img, _, err = image.Decode(bytes.NewReader(input))
		}
		if err != nil {
			m.global.Call("alert", "Image not supportaded")
			return nil
		}

		imgWithEffects := m.applyEffects(img)

		if m.checkAscii {
			m.asciiGenerator(imgWithEffects)
		} else {
			m.imageEffectGenerator(imgWithEffects)
		}

		onLoad.Release()

		inputZoomRangeDiv := m.document.Call("getElementById", "ascii-div")
		changeAttribute(inputZoomRangeDiv, "data-visible", strconv.FormatBool(m.checkAscii))

		return nil
	})

	fileReader.Set("onload", onLoad)
	fileReader.Call("readAsArrayBuffer", m.imageSelected)
}

func (m *model) imageEffectGenerator(img image.Image) {

	var buf bytes.Buffer
	png.Encode(&buf, img)
	data := buf.Bytes()
	uint8Array := m.global.Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(uint8Array, data)

	array := m.global.Get("Array").New(1)
	array.SetIndex(0, uint8Array)

	blobOpt := m.global.Get("Object").New()
	blobOpt.Set("type", "image/png")
	blob := m.global.Get("Blob").New(array, blobOpt)

	image := m.global.Get("Image").New()
	image.Set("src", m.global.Get("URL").Call("createObjectURL", blob))

	image.Call("addEventListener", "load", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		canvas := m.document.Call("getElementById", "canvas")
		drawImage := canvas.Call("getContext", "2d")

		var imgWidth = image.Get("width")
		var imgHeight = image.Get("height")

		if !imgWidth.IsUndefined() && !imgWidth.IsNull() && !imgHeight.IsUndefined() && !imgHeight.IsNull() {

			changeAttribute(canvas, "width", fmt.Sprintf("%d", imgWidth.Int()))
			changeAttribute(canvas, "height", fmt.Sprintf("%d", imgHeight.Int()))

			drawImage.Call("drawImage", image, 0, 0)
		}
		return nil
	}))
}

func (m *model) asciiGenerator(img image.Image) {
	density := []rune(m.asciiChars)

	width, height := resizeImg(img, m.imageWidth)
	resul := transform.Resize(img, width, height, transform.Linear)
	fontSize, lineHeight := m.resizeAscii()
	bounds := resul.Bounds()

	asciiImage := m.global.Get("Array").New()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {

		line := m.global.Get("Array").New()
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := resul.At(x, y)
			gr := color.GrayModel.Convert(px)
			gray := gr.(color.Gray)

			intensity := float64(gray.Y) / 255.0
			charIndex := math.Floor(float64(len(density)-1) * intensity)

			jsImgInfo := m.global.Get("Object").New()
			jsImgInfo.Set("char", string(density[int(charIndex)]))
			jsImgInfo.Set("fontSize", fontSize)
			jsImgInfo.Set("lineHeight", lineHeight)

			if m.checkColor {
				r, g, b, _ := px.RGBA()
				colorCSS := fmt.Sprintf("rgb(%d,%d,%d)", r>>8, g>>8, b>>8)
				jsImgInfo.Set("color", colorCSS)
			} else {
				jsImgInfo.Set("color", "#fff")
			}

			line.Call("push", jsImgInfo)
		}
		asciiImage.Call("push", line)
	}
	m.global.Call("generateCanva", asciiImage)
}

func (m model) applyEffects(img image.Image) image.Image {
	var result image.Image = img

	switch m.effectSelected {
	case "kuwahara":
		result = kuwahara.KuwaharaFilter(result, int(m.effectRange))
	case "dithering":
		result = dithering.OrderedDithering(result, int(m.effectRange), 2)
	case "hue":
		result = hsl.Hue(result, int(m.effectRange))
	case "saturation":
		result = hsl.Saturation(result, float64(m.effectRange))
	case "brightness":
		result = hsl.Luminance(result, float64(m.effectRange))
	case "flipH":
		result = flip.FlipHorizontal(result)
	case "flipV":
		result = flip.FlipVertical(result)
	case "gaussianBlur":
		result = blur.GaussianBlur(result, int(m.effectRange))

	case "Dilate":
		result = effect.Dilate(result, float64(m.effectRange))
	case "edgeDetection":
		result = effect.EdgeDetection(result, float64(m.effectRange))
	case "emboss":
		result = effect.Emboss(result)
	case "erode":
		result = effect.Erode(result, float64(m.effectRange))
	case "grayscale":
		result = effect.Grayscale(result)
	case "invert":
		result = effect.Invert(result)
	case "median":
		result = effect.Median(result, float64(m.effectRange))
	case "sepia":
		result = effect.Sepia(result)
	case "sharpen":
		result = effect.Sharpen(result)
	case "sobale":
		result = effect.Sobel(result)
	case "contrast":
		result = adjust.Contrast(result, float64(m.effectRange))
	case "gamma":
		result = adjust.Gamma(result, float64(m.effectRange))
	case "threshold":
		result = segment.Threshold(result, uint8(m.effectRange))

	case "shearH":
		result = transform.ShearH(result, float64(m.effectRange))
	case "shearV":
		result = transform.ShearV(result, float64(m.effectRange))

	case "noiseUniformColored":
		imgBounds := result.Bounds()
		noise := noise.Generate(imgBounds.Dx(), imgBounds.Dy(), &noise.Options{Monochrome: false, NoiseFn: noise.Uniform})
		result = blend.Overlay(noise, result)
	case "noiseBinaryMonochrome":
		imgBounds := result.Bounds()
		noise := noise.Generate(imgBounds.Dx(), imgBounds.Dy(), &noise.Options{Monochrome: true, NoiseFn: noise.Binary})
		result = blend.Opacity(noise, result, 0.5)
	case "noiseGaussianMonochrome":
		imgBounds := result.Bounds()
		noise := noise.Generate(imgBounds.Dx(), imgBounds.Dy(), &noise.Options{Monochrome: true, NoiseFn: noise.Gaussian})
		result = blend.Overlay(noise, result)
	case "noisePerlin":
		imgBounds := result.Bounds()
		noise := noise.GeneratePerlin(imgBounds.Dx(), imgBounds.Dy(), m.effectRange)
		result = blend.Overlay(noise, result)
	}

	return result
}

func (m model) resizeAscii() (int, int) {

	imageDefaultSize := 100

	lineHeightRatio := float64(m.lineHeight) / float64(m.fontSize)
	inverseFactor := float64(imageDefaultSize) / float64(m.imageWidth)
	newFontSize := float64(m.fontSize) * inverseFactor
	newLineHeight := newFontSize * lineHeightRatio

	newFontSize = math.Max(newFontSize, float64(m.fontSize))
	newLineHeight = math.Max(newLineHeight, float64(m.lineHeight))

	return int(newFontSize), int(newLineHeight)
}

func imageSize(img image.Image) (int, int) {
	imgBounds := img.Bounds()
	return imgBounds.Dx(), imgBounds.Dy()
}

func resizeImg(img image.Image, newWidth int) (int, int) {
	imgBounds := img.Bounds()
	aspectRatio := float64(newWidth) / float64(imgBounds.Dx())
	newHeight := int(float64(imgBounds.Dy()) * aspectRatio)

	return newWidth, newHeight
}

func changeAttribute(content js.Value, attribute string, value string) {
	content.Call("setAttribute", attribute, value)
}

func effectsRateMapFunc() map[string]bool {
	effects := []string{"blur", "Dilate", "edgeDetection", "erode", "median"}
	effectsMap := make(map[string]bool)
	for _, effect := range effects {
		effectsMap[effect] = true
	}

	return effectsMap
}

func debounce(duration time.Duration, fn func()) func() {
	var timer *time.Timer

	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(duration, fn)
	}
}
