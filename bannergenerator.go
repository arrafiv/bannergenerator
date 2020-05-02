package bannergenerator

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
)

type (
	//ImageLayer is a struct
	ImageLayer struct {
		Image image.Image
		XPos  int
		YPos  int
	}

	//Label is a struct
	Label struct {
		Text     string
		FontPath string
		FontType string
		FontSize float64
		XPos     int
		YPos     int
	}

	//BgProperty is background property struct
	BgProperty struct {
		Width   int
		Length  int
		BgColor color.Color
	}
)

// GenerateBanner is a function that combine images and texts into one image
func GenerateBanner(imgs []ImageLayer, labels []Label, bgProperty BgProperty) (*image.RGBA, error) {
	//create image's background
	bgImg := image.NewRGBA(image.Rect(0, 0, bgProperty.Width, bgProperty.Length))

	//set the color to white
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{bgProperty.BgColor}, image.ZP, draw.Src)

	for _, img := range imgs {
		//set image offset
		offset := image.Pt(img.XPos, img.YPos)

		//combine the image
		draw.Draw(bgImg, img.Image.Bounds().Add(offset), img.Image, image.ZP, draw.Over)
	}

	//add label(s)
	bgImg, err := addLabel(bgImg, labels)
	if err != nil {
		return nil, err
	}

	return bgImg, nil
}

func addLabel(img *image.RGBA, labels []Label) (*image.RGBA, error) {
	for _, label := range labels {
		// Read font data
		fontBytes, err := ioutil.ReadFile(label.FontPath + label.FontType)
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
		c.SetFontSize(label.FontSize)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.Black)

		pt := freetype.Pt(label.XPos, label.YPos+int(c.PointToFixed(label.FontSize)>>6))

		_, err = c.DrawString(label.Text, pt)
		if err != nil {
			log.Println(err)
			return img, nil
		}
		pt.Y += c.PointToFixed(label.FontSize * 1.5)
	}

	return img, nil
}
