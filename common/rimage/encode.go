package rimage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"sniffle/tool"
	"sniffle/tool/fetch"

	"github.com/nfnt/resize"
)

func FetchResizeJpeg(t *tool.Tool, url []byte) ([]byte, error) {
	out, err := fetchAndResize(t, url)
	if err != nil {
		return nil, err
	}

	buff := bytes.Buffer{}
	jpeg.Encode(&buff, out, &jpeg.Options{
		Quality: jpeg.DefaultQuality,
	})
	return buff.Bytes(), nil
}

func fetchAndResize(t *tool.Tool, url []byte) (image.Image, error) {
	r := t.Fetch(fetch.URL(string(url)))
	if r == nil {
		return nil, fmt.Errorf("cannot fetch image %q", url)
	}
	defer r.Body.Close()

	img, _, err := image.Decode(r.Body)
	if err != nil {
		return nil, fmt.Errorf("decode image %q: %w", url, err)
	}

	size := uint(300)
	return resize.Thumbnail(size, size, img, resize.Lanczos3), nil
}
