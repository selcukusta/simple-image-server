package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	i "github.com/disintegration/imaging"
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

	if err != nil || img == nil {
		return nil, fmt.Sprintf("Image (%s) decode process is failed", contentType), err
	}

	width, err := strconv.Atoi(params["width"])
	if err != nil {
		return nil, "Invalid width param", err
	}

	height, err := strconv.Atoi(params["height"])
	if err != nil {
		return nil, "Invalid height param", err
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if (width == w && height == h) || (width == 0 && height == 0) {
		return imageProcess(params, img, contentType)
	}

	if !strings.ContainsAny("tc", params["option"]) || width == 0 || height == 0 {
		img = i.Resize(img, width, height, i.Lanczos)
		return imageProcess(params, img, contentType)
	}

	if strings.ContainsAny("c", params["option"]) {
		img = i.CropCenter(img, width, height)
		return imageProcess(params, img, contentType)
	}

	if strings.ContainsAny("t", params["option"]) {
		img = i.Thumbnail(img, width, height, i.Lanczos)
		return imageProcess(params, img, contentType)
	}

	return nil, "Unknown parameters or sizes", err
}

func imageProcess(params map[string]string, img image.Image, contentType string) ([]byte, string, error) {

	if strings.ContainsAny("g", params["option"]) {
		img = i.Grayscale(img)
	}

	buf := new(bytes.Buffer)
	switch contentType {
	case "image/png":
		err := png.Encode(buf, img)
		if err != nil {
			return nil, "Invalid image/png encoding operation", err
		}
	case "image/jpeg":
		quality, _ := strconv.Atoi(params["quality"])
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, "Invalid image/jpeg encoding operation", err
		}
	}

	return buf.Bytes(), "", nil
}
