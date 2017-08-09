// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

import (
	"log"
	"net/http"
	"time"
)

var (
	Green   = color.RGBA{0, 185, 0, 255}
	Black   = color.Black
	palette = []color.Color{Black, Green}
)

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w)
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles    = 10              // number of complete x oscillator revolutions
		res       = 0.001           // angular resolution
		nframes   = 256             // number of animation frames
		delay     = 8               // delay between frames in 10ms units
		sizeScale = 3               // size multiplier
		size      = sizeScale * 100 // image canvas covers [-size..+size]
		nshapes   = 8               // Number of shape changes
	)

	anim := gif.GIF{LoopCount: nframes}

	for i := 0; i < nshapes; i++ {
		freq := rand.Float64() * 3.0 // relative frequency of y oscillator
		phase := 0.0                 // phase difference
		for i := 0; i < nframes/nshapes; i++ {
			rect := image.Rect(0, 0, 2*size+1, 2*size+1)
			img := image.NewPaletted(rect, palette)
			for t := 0.0; t < cycles*2*math.Pi; t += res {
				x := math.Sin(t)
				y := math.Sin(t*freq + phase)
				img.SetColorIndex(
					size+int(x*size+0.5),
					size+int(y*size+0.5),
					blackIndex)
			}
			phase += 0.1
			anim.Delay = append(anim.Delay, delay)
			anim.Image = append(anim.Image, img)
		}
	}
	err := gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
	if err != nil {
		log.Fatal(err)
	}
}
