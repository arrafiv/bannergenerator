package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "image/png"

	"github.com/golang/freetype"
)

type (
	label struct {
		text     string
		fontType string
		fontSize float64
		xPos     int
		yPos     int
	}
)

func main() {
	mainImg, err := downloadMainImage("https://www.gstatic.com/webp/gallery3/1.png")
	if err != nil {
		log.Println(err)
		return
	}

	side := 500

	//create image's background
	bgImg := image.NewRGBA(image.Rect(0, 0, side, 380))
	white := color.RGBA{227, 221, 221, 1}

	//set the color to white
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	//set the main image to center
	offsetSide := (side - mainImg.Bounds().Max.X) / 2
	offset := image.Pt(200, offsetSide)

	//combine the image
	draw.Draw(bgImg, mainImg.Bounds().Add(offset), mainImg, image.ZP, draw.Over)

	//add label
	labels := []label{
		label{
			fontSize: 48,
			fontType: "BebasNeue-Regular.ttf",
			text:     "Tumbuhan &",
			xPos:     10,
			yPos:     0,
		},
		label{
			fontSize: 48,
			fontType: "BebasNeue-Regular.ttf",
			text:     "Tanaman",
			xPos:     10,
			yPos:     50,
		},
		label{
			fontSize: 32,
			fontType: "BebasNeue-Light.ttf",
			text:     "di bawah",
			xPos:     10,
			yPos:     290,
		},
		label{
			fontSize: 48,
			fontType: "BebasNeue-Bold.ttf",
			text:     "Rp 90rb",
			xPos:     10,
			yPos:     320,
		},
	}
	bgImg, err = addLabel(bgImg, labels)
	if err != nil {
		log.Println(err)
		return
	}

	out, err := os.Create("./output.jpg")
	if err != nil {
		log.Println(err)
		return
	}

	var opt jpeg.Options
	opt.Quality = 80

	jpeg.Encode(out, bgImg, &opt)

}

func downloadMainImage(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return m, err
}

func addLabel(img *image.RGBA, labels []label) (*image.RGBA, error) {
	for _, label := range labels {
		// Read font data
		fontBytes, err := ioutil.ReadFile("../../../golang/freetype/testdata/" + label.fontType)
		if err != nil {
			return nil, err
		}
		f, err := freetype.ParseFont(fontBytes)
		if err != nil {
			return nil, err
		}

		// Initialize the context
		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(f)
		c.SetFontSize(label.fontSize)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.Black)

		pt := freetype.Pt(label.xPos, label.yPos+int(c.PointToFixed(label.fontSize)>>6))

		_, err = c.DrawString(label.text, pt)
		if err != nil {
			log.Println(err)
			return img, nil
		}
		pt.Y += c.PointToFixed(label.fontSize * 1.5)
	}

	return img, nil
}
