package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"net"
	"strconv"
	"strings"

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
	i := 0
	last_bat := "NAN"
	for {
		background := color.RGBA{R: 20, G: 20, B: 20, A: 255}
		draw.Draw(img, img.Bounds(), image.NewUniform(background), image.Point{}, draw.Src)

		cursor := freetype.Pt(20, 40)
		drop := func(y int) {
			cursor = cursor.Add(freetype.Pt(0, y))
		}
		right := func() {
			cursor = cursor.Add(freetype.Pt(280, 0))
		}
		left := func() {
			cursor = cursor.Sub(freetype.Pt(280, 0))
		}

		now := time.Now()
		if i%5 == 0 {
			bat, err := get_bat()
			if err != nil {
				fmt.Println(err)
			}
			last_bat = strconv.Itoa(bat) + "%"
		}

		trimlead := func(disp string) string {
			if disp[0] == '0' {
				return " " + disp[1:]
			}
			return disp
		}
		tdisp := func(loc *time.Location) string {
			disp := zonetimes.GetTimeWithLoc(now, loc, tick)
			if loc == zones.BRZ {
				disp = strings.Replace(disp, "-03", "BRZ", 1)
			}
			return trimlead(disp)
		}

		zone, time_string := zonetimes.GetDisplayTime(now, zones.NZT, tick)
		medium_font.DrawString(zone+" time:", cursor)
		drop(80)
		large_font.DrawString(trimlead(time_string), cursor)

		drop(90)

		zone, time_string = zonetimes.GetDisplayTime(now, zones.HKT, tick)
		medium_font.DrawString(zone+" time:", cursor)
		drop(80)
		large_font.DrawString(trimlead(time_string), cursor)

		drop(90)

		medium_font.DrawString(fmt.Sprintf("Epoch time: %d", now.Unix()), cursor)
		drop(60)
		medium_font.DrawString(tdisp(zones.CET), cursor)
		drop(60)
		medium_font.DrawString(tdisp(zones.JPN), cursor)
		drop(60)
		medium_font.DrawString(tdisp(zones.UTC), cursor)
		right()
		medium_font.DrawString("BAT:", cursor)
		left()
		drop(60)
		medium_font.DrawString(tdisp(zones.USP), cursor)
		right()
		medium_font.DrawString(last_bat, cursor)
		left()
		drop(60)
		medium_font.DrawString(tdisp(zones.USE), cursor)
		drop(60)
		medium_font.DrawString(tdisp(zones.BRZ), cursor)
		drop(60)

		rotated := rotate_90(img)
		display.Blit(rotated)

		time.Sleep(time.Second)
		tick = !tick
		i++
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

func rotate_90(src *image.RGBA) *image.RGBA {
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

func get_bat() (int, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8423")
	if err != nil {
		return 0, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("get battery\n"))
	if err != nil {
		return 0, fmt.Errorf("failed to send command: %w", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var level int
	_, err = fmt.Sscanf(response, "battery: %d", &level)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return level, nil
}
