//go:build cgo

package resize0

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gen2brain/avif"
	"github.com/nfnt/resize"
)

// Resize image and save it as AVIF format.
func Resize(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	out := resize.Thumbnail(Size, Size, img, resize.Lanczos3)

	buff := bytes.Buffer{}
	avif.Encode(&buff, out, avif.Options{
		Quality: 55,
		Speed:   0,
	})

	return buff.Bytes(), nil
}
