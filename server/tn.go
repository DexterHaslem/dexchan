package server

import (
	"path/filepath"
	"image"
	"image/jpeg"
	"os"
	"github.com/nfnt/resize"
	"image/png"
	"bytes"
)

const (
	// beware: resize docs show using height 0 to keep aspect,
	// but thumbnail doesnt work like that, it trys to resize both bounds
	// with keeping of aspect ratio

	MaxWidth  = 350
	MaxHeight = 1000
)

// createThumbnail will create thumbnail of file already open
func createThumbnail(f *os.File) ([]byte, error) {
	ext := filepath.Ext(f.Name())

	if ext == ".webm" {
		return nil, nil
	}

	// we have to provide a decoded image to resizer. woopie
	var img image.Image
	var err error

	isJpeg := false
	switch ext {
	case ".jpg":
		fallthrough
	case ".jpeg":
		img, err = jpeg.Decode(f)
		isJpeg = true
	case ".png":
		img, err = png.Decode(f)
	}

	if err != nil {
		return nil, err
	}
	//t := resize.Thumbnail(MaxWidth, MaxHeight, img, resize.Lanczos3)
	t := resize.Resize(MaxWidth, 0, img, resize.Lanczos3)

	bb := &bytes.Buffer{}
	if isJpeg {
		err = jpeg.Encode(bb, t, nil)
	} else {
		err = png.Encode(bb, t)
	}

	return bb.Bytes(), err
}
