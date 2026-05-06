package main

import (
	_ "embed"
	"fmt"

	"image"
	"image/color"
	"image/draw"
	"log"
	"time"

	"github.com/amnonbc/pidisp"
	"github.com/chrispritchard/onit/internal/zonetimes"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed roboto-mono.ttf
var font_bytes []byte

func main() {
	display, err := pidisp.Open(pidisp.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer display.Close()

	font_face, err := freetype.ParseFont(font_bytes)
	if err != nil {
		log.Fatal("Could not parse font:", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, display.Height(), display.Width())) // rotated dims

	medium_font := get_font(font_face, 32, color.White, img)
	large_font := get_font(font_face, 72, color.White, img)

	zones := zonetimes.GetZones()

	tick := true
	for {
		background := color.RGBA{R: 20, G: 20, B: 20, A: 255}
		draw.Draw(img, img.Bounds(), image.NewUniform(background), image.Point{}, draw.Src)

		cursor := freetype.Pt(20, 40)
		drop := func(y int) {
			cursor = cursor.Add(freetype.Pt(0, y))
		}

		now := time.Now()

		zone, time_string := zonetimes.GetDisplayTime(now, zones.NZT, tick)
		medium_font.DrawString(zone+" time:", cursor)
		drop(80)
		large_font.DrawString(time_string, cursor)

		drop(90)

		zone, time_string = zonetimes.GetDisplayTime(now, zones.HKT, tick)
		medium_font.DrawString(zone+" time:", cursor)
		drop(80)
		large_font.DrawString(time_string, cursor)

		drop(90)

		medium_font.DrawString(fmt.Sprintf("Epoch time: %d", now.Unix()), cursor)
		drop(60)
		medium_font.DrawString(zonetimes.GetTimeWithLoc(now, zones.UTC, tick), cursor)
		drop(60)
		medium_font.DrawString(zonetimes.GetTimeWithLoc(now, zones.CET, tick), cursor)
		drop(60)
		medium_font.DrawString(zonetimes.GetTimeWithLoc(now, zones.USE, tick), cursor)
		drop(60)
		medium_font.DrawString(zonetimes.GetTimeWithLoc(now, zones.USP, tick), cursor)
		drop(60)
		medium_font.DrawString(zonetimes.GetTimeWithLoc(now, zones.BRZ, tick), cursor)
		drop(60)
		medium_font.DrawString(zonetimes.GetTimeWithLoc(now, zones.JPN, tick), cursor)
		drop(60)

		rotated := rotate90Clockwise(img)
		display.Blit(rotated)

		time.Sleep(time.Second)
		tick = !tick
	}
}

func get_font(font_face *truetype.Font, font_size float64, colour color.Color, target *image.RGBA) *freetype.Context {
	context := freetype.NewContext()
	context.SetDPI(72)
	context.SetFont(font_face)
	context.SetFontSize(font_size)
	context.SetClip(target.Bounds())
	context.SetDst(target)
	context.SetSrc(image.NewUniform(colour)) // White text
	context.SetHinting(font.HintingNone)
	return context
}

func rotate90Clockwise(src *image.RGBA) *image.RGBA {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcH, srcW))

	for y := range srcH {
		for x := range srcW {
			dst.Set(srcH-1-y, x, src.At(x, y))
		}
	}
	return dst
}
