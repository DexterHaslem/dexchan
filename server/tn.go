// dexChan copyright Dexter Haslem <dmh@fastmail.com> 2018
// This file is part of dexChan
//
// dexChan is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dexChan is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dexChan. If not, see <http://www.gnu.org/licenses/>.


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
