package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const (
	volumes = 100
	count   = 10000 // images per volume
	size    = 16 
	zoom    = 8
)

var bits = map[byte][4]byte{
	'0': {0, 0, 0, 0},
	'1': {0, 0, 0, 1},
	'2': {0, 0, 1, 0},
	'3': {0, 0, 1, 1},
	'4': {0, 1, 0, 0},
	'5': {0, 1, 0, 1},
	'6': {0, 1, 1, 0},
	'7': {0, 1, 1, 1},
	'8': {1, 0, 0, 0},
	'9': {1, 0, 0, 1},
	'a': {1, 0, 1, 0},
	'b': {1, 0, 1, 1},
	'c': {1, 1, 0, 0},
	'd': {1, 1, 0, 1},
	'e': {1, 1, 1, 0},
	'f': {1, 1, 1, 1},
}

var paletteBW = color.Palette{
	color.RGBA{0x00, 0x00, 0x00, 255},
	color.RGBA{0xff, 0xff, 0xff, 255},
}

func buildVolume(in *bufio.Reader, vol int) {
	id := 0
	img := image.NewPaletted(image.Rect(0, 0, size*zoom, size*zoom), paletteBW)
	for n := 0; n < count; n++ {
		x := 0
		y := 0
		for {
			digit, err := in.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			if digit == '.' {
				continue
			}
			for _, color := range bits[digit] {
				for px := 0; px < zoom; px++ {
					for py := 0; py < zoom; py++ {
						img.SetColorIndex(x+px, y+py, color)
					}
				}
				x += zoom
			}
			if x >= size*zoom {
				x = 0
				y += zoom
				if y >= size*zoom {
					break
				}
			}
		}
		fileName := fmt.Sprintf("images/pi-volume-%02d/pi-%02d-%04d.png", vol, vol, id)
		fmt.Println(fileName)
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		if err := png.Encode(f, img); err != nil {
			f.Close()
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		id++
	}
}

func main() {
	fin, err := os.OpenFile("pi/pi_hex_1b.txt", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	in := bufio.NewReader(fin)

	for vol := 0; vol < volumes; vol++ {
		volDir := fmt.Sprintf("images/pi-volume-%02d", vol)
		if err := os.MkdirAll(volDir, 0755); err != nil {
			log.Fatal(err)
		}
		buildVolume(in, vol)
	}
}
