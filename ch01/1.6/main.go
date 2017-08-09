// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	Black   = color.RGBA{0x0E, 0x0E, 0x0E, 192} // tansparency: 75%
	palette = []color.Color{Black}              // set background
)

func init() {
	// Set our gradient colors ("rainbow")
	for _, clr := range []color.Color{
		color.RGBA{0xce, 0x2c, 0xa1, 0xff},
		color.RGBA{0xd1, 0x50, 0x2f, 0xff},
		color.RGBA{0xbf, 0xd5, 0x32, 0xff},
		color.RGBA{0x36, 0xd8, 0x40, 0xff},
		color.RGBA{0x3a, 0xdc, 0xdb, 0xff},
		color.RGBA{0x3e, 0x49, 0xe0, 0xff},
	} {
		palette = append(palette, clr)
	}
}

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/gif")
			lissajous(w)
		}
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8000", nil))
		return
	}
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles    = 10              // number of complete x oscillator revolutions
		res       = 0.001           // angular resolution
		nframes   = 64              // number of animation frames
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
				colorIdx := uint8(i%(len(palette)-1) + 1)
				img.SetColorIndex(
					size+int(x*size+0.5),
					size+int(y*size+0.5),
					colorIdx)
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
