package imageprocessor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	r "github.com/nfnt/resize"
)

//ImageProcess will be used to resize the image, create the thumbnail from the image and change the color mode of image to the grayscale.
func ImageProcess(params map[string]string, imageAsByte []byte, contentType string) ([]byte, string, error) {
	var (
		img image.Image
		err error
	)
	switch contentType {
	case "image/png":
		img, err = png.Decode(bytes.NewReader(imageAsByte))
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(imageAsByte))
	}

	if err != nil {
		return nil, fmt.Sprintf("Image (%s) decode process is failed", contentType), err
	}
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()

	width, err := strconv.Atoi(params["width"])
	if err != nil {
		return nil, "Invalid width param", err
	}

	height, err := strconv.Atoi(params["height"])
	if err != nil {
		return nil, "Invalid height param", err
	}

	if (width != x || height != y) && (width != 0 || height != 0) {
		if strings.ContainsAny("t", params["option"]) {
			img = thumbImage(img, width, height)
		} else {
			img = resizeImage(img, width, height, strings.ContainsAny("s", params["option"]))
		}
	}

	if strings.ContainsAny("g", params["option"]) {
		img = grayscaleImage(img)
	}

	buf := new(bytes.Buffer)
	switch contentType {
	case "image/png":
		err = png.Encode(buf, img)
		if err != nil {
			return nil, "Invalid image/png encoding operation", err
		}
	case "image/jpeg":
		quality, _ := strconv.Atoi(params["quality"])
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, "Invalid image/jpeg encoding operation", err
		}
	}

	return buf.Bytes(), "", nil
}

func thumbImage(img image.Image, w int, h int) image.Image {
	max := w
	if w < h {
		max = h
	}
	return r.Thumbnail(uint(max), uint(max), img, r.Lanczos3)
}

func resizeImage(img image.Image, w int, h int, smart bool) image.Image {
	resizer := nfnt.NewResizer(r.Lanczos3)
	if smart {
		analyzer := smartcrop.NewAnalyzer(resizer)
		topCrop, _ := analyzer.FindBestCrop(img, w, h)
		type SubImager interface {
			SubImage(r image.Rectangle) image.Image
		}
		img = img.(SubImager).SubImage(topCrop)
	}
	return resizer.Resize(img, uint(w), uint(h))
}

func grayscaleImage(img image.Image) image.Image {
	grayscale := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < grayscale.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayscale.Set(x, y, img.At(x, y))
		}
	}
	return grayscale
}
