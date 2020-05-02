package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	_ "image/png"

	m "github.com/arrafiv/img-txt-combiner"
)

func main() {
	tempImg1, err := downloadMainImage("http://www.pngmart.com/files/4/Aloe-PNG-Image.png")
	if err != nil {
		log.Println(err)
		return
	}

	tempImg2, err := downloadMainImage("https://www.gstatic.com/webp/gallery3/1.png")
	if err != nil {
		log.Println(err)
		return
	}

	imgs := []m.ImageLayer{
		m.ImageLayer{
			Image: tempImg1,
			XPos:  200,
			YPos:  -100,
		},
		m.ImageLayer{
			Image: tempImg2,
			XPos:  200,
			YPos:  100,
		},
	}

	bg := m.BgProperty{
		Width:   500,
		Length:  380,
		BgColor: color.RGBA{227, 221, 221, 1},
	}

	//add label
	labels := []m.Label{
		m.Label{
			FontPath: "../../golang/freetype/testdata/",
			FontSize: 48,
			FontType: "BebasNeue-Regular.ttf",
			Text:     "Tumbuhan &",
			XPos:     10,
			YPos:     0,
		},
		m.Label{
			FontPath: "../../golang/freetype/testdata/",
			FontSize: 48,
			FontType: "BebasNeue-Regular.ttf",
			Text:     "Tanaman",
			XPos:     10,
			YPos:     50,
		},
		m.Label{
			FontPath: "../../golang/freetype/testdata/",
			FontSize: 32,
			FontType: "BebasNeue-Light.ttf",
			Text:     "di bawah",
			XPos:     10,
			YPos:     290,
		},
		m.Label{
			FontPath: "../../golang/freetype/testdata/",
			FontSize: 48,
			FontType: "BebasNeue-Bold.ttf",
			Text:     "Rp 90rb",
			XPos:     10,
			YPos:     320,
		},
	}

	res, err := m.GenerateBanner(imgs, labels, bg)
	if err != nil {
		log.Printf("Error generating banner: %+v\n", err)
	}

	out, err := os.Create("./output.jpg")
	if err != nil {
		log.Printf("Error creating image file: %+v\n", err)
		return
	}

	var opt jpeg.Options
	opt.Quality = 80

	err = jpeg.Encode(out, res, &opt)
	if err != nil {
		log.Printf("Error creating image file: %+v\n", err)
		return
	}

	log.Println("Image Generated")
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
