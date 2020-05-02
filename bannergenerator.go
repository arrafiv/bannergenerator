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
		Size     float64
		Color    image.Image
		DPI      float64
		Spacing  float64
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

	//set the background color
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{bgProperty.BgColor}, image.ZP, draw.Src)

	//looping image layer, higher array index = upper layer
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
	//initialize the context
	c := freetype.NewContext()

	for _, label := range labels {
		//read font data
		fontBytes, err := ioutil.ReadFile(label.FontPath + label.FontType)
		if err != nil {
			return nil, err
		}
		f, err := freetype.ParseFont(fontBytes)
		if err != nil {
			return nil, err
		}

		//set label configuration
		c.SetDPI(label.DPI)
		c.SetFont(f)
		c.SetFontSize(label.Size)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(label.Color)

		//positioning the label
		pt := freetype.Pt(label.XPos, label.YPos+int(c.PointToFixed(label.Size)>>6))

		//draw the label on image
		_, err = c.DrawString(label.Text, pt)
		if err != nil {
			log.Println(err)
			return img, nil
		}
		pt.Y += c.PointToFixed(label.Size * label.Spacing)
	}

	return img, nil
}
