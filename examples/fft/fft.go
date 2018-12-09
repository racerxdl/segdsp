package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/dsp/fft"
	"github.com/racerxdl/segdsp/tools"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
)

const sampleRate = 1e6
const fftSize = 4096
const height = 4096
const gridSteps = 8
const fftOffset = -20
const fftScale = 20
const fftFilterAlpha = 0.4 // 0.4

var loadedFont *truetype.Font

func drawGrid(gc *draw2dimg.GraphicContext, fftOffset, fftScale, width int) {
	gc.Save()
	gc.SetLineWidth(4)
	gc.SetFontSize(128)
	gc.SetStrokeColor(color.RGBA{R: 127, G: 127, B: 127, A: 255})
	gc.SetFillColor(color.RGBA{R: 127, G: 127, B: 127, A: 255})
	// region Draw dB Scale Grid
	for i := 0; i < gridSteps; i++ {
		var y = float64(i) * (float64(height) / float64(gridSteps))
		var dB = int(float64(fftOffset) - float64(y)/float64(fftScale))
		gc.MoveTo(0, y)
		gc.LineTo(float64(width), y)
		gc.FillStringAt(fmt.Sprintf("%d dB", dB), 20, y)
	}
	// endregion
	// region Draw Frequency Scale Grid

	// endregion
	gc.Close()
	gc.Stroke()
	gc.Restore()
}

func main() {
	// region Load Font
	fontBytes, err := ioutil.ReadFile("FreeMono.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	loadedFont, err = truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	draw2d.RegisterFont(draw2d.FontData{
		Name:   "FreeMono",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	}, loadedFont)
	// endregion
	// region Generate IQ Sample
	var samples = make([]complex64, 1024)
	var interp = dsp.MakeInterpolator(20)
	var lowPass = dsp.MakeLowPassFixed(1, sampleRate, sampleRate/2-5e3, 63)
	var frequencyShift = dsp.MakeFrequencyTranslator(1, -100e3, sampleRate, lowPass)

	for i := 0; i < len(samples); i++ {
		samples[i] = complex((rand.Float32()-1)*0.5, 0)
	}

	samples = interp.Work(samples)

	// Generate some background noise
	for i := 0; i < len(samples); i++ {
		samples[i] += complex((rand.Float32()-1)*1e-4, 0)
	}

	samples = frequencyShift.Work(samples)
	// Should have a 50 kHz signal shifted to -100 kHz
	log.Printf("Generated %d samples.\n", len(samples))
	// Uncomment for saving IQ Samples
	//f , _ := os.Create("test.iq")
	//for i := 0; i < 16; i++ {
	//	for j := 0; j < len(samples); j++ {
	//		_ = binary.Write(f, binary.LittleEndian, samples[j])
	//	}
	//}
	//f.Close()
	// endregion
	// region Compute FFT
	window := dsp.BlackmanHarris(fftSize, 61)
	fftSamples := make([]complex64, fftSize)
	copy(fftSamples, samples)

	for j := 0; j < fftSize; j++ {
		var s = fftSamples[j]
		var r = real(s) * float32(window[j])
		var i = imag(s) * float32(window[j])
		fftSamples[j] = complex(r, i)
	}

	fftResult := fft.FFT(fftSamples)
	fftReal := make([]float32, len(fftResult))
	for i := 0; i < len(fftResult); i++ {
		// Convert FFT to Power in dB
		var v = tools.ComplexAbsSquared(fftResult[i]) * (1.0 / sampleRate)
		fftReal[i] = float32(10 * math.Log10(float64(v)))
	}

	var maxVal = float32(-999999999)
	var minVal = float32(999999999)
	for i := 0; i < len(fftReal); i++ {
		maxVal = float32(math.Max(float64(maxVal), float64(fftReal[i])))
		minVal = float32(math.Min(float64(minVal), float64(fftReal[i])))
	}
	var delta = maxVal - minVal
	log.Println("Max Val: ", maxVal)
	log.Println("Min Val: ", minVal)
	log.Println("Delta: ", delta)

	// Filter FFT
	for i := 1; i < len(fftReal); i++ {
		fftReal[i] = fftReal[i-1]*fftFilterAlpha + fftReal[i]*(1-fftFilterAlpha)
	}

	// endregion
	// region Draw Image and Save
	img := image.NewRGBA(image.Rect(0, 0, len(fftReal), height))

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFontData(draw2d.FontData{
		Name:   "FreeMono",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	})

	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{R: 255, A: 255})
	gc.SetFillColor(color.RGBA{R: 255, A: 255})

	gc.SetFontSize(256)
	gc.FillStringAt("FFT", 100, 256)

	var startV float64

	for i := 0; i < len(fftReal); i++ {
		var iPos = (i + len(fftReal)/2) % len(fftReal)
		var s = float64(fftReal[iPos])
		var v = (float64(fftOffset) - s) * float64(fftScale)
		var x = float64(i)
		if i == 0 {
			startV = v
			gc.MoveTo(x, v)
		} else {
			gc.LineTo(x, v)
		}
	}

	gc.LineTo(float64(len(samples)), float64(height))
	gc.LineTo(0, float64(height))
	gc.LineTo(0, startV)
	gc.Close()

	gc.FillStroke()
	gc.Fill()

	drawGrid(gc, fftOffset, fftScale, len(fftReal))

	f, err := os.Create("test.jpg")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f, img, nil)
	if err != nil {
		panic(err)
	}
	// endregion
}
