// Resize and reencode image.
package rimage

import (
	"bytes"
	"image"
	"strconv"

	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

const (
	NameResizeJpeg = "img-fetch-resize-jpeg"

	ExtensionJpeg = ".jpg"
)

type Image struct {
	// With and height as string because use will use it in attribute.
	Height string
	Width  string

	JPEG []byte `json:"-"`
}

func New(t *tool.Tool, url string) *Image {
	if url == "" {
		return nil
	}

	jpeg := t.LongTask(NameResizeJpeg, url, []byte(url))
	config, _, err := image.DecodeConfig(bytes.NewReader(jpeg))
	if err != nil {
		return nil
	}

	return &Image{
		Height: strconv.Itoa(config.Height),
		Width:  strconv.Itoa(config.Width),

		JPEG: jpeg,
	}
}

func (img *Image) Save(t *tool.Tool, path string) {
	if img == nil || len(img.JPEG) == 0 {
		return
	}
	t.WriteFile(path+ExtensionJpeg, img.JPEG)
}

func (img *Image) Render(base, title string) render.Node {
	if img == nil {
		return render.Z
	}

	return render.Na("img.logo", "loading", "lazy").
		A("src", base+ExtensionJpeg).
		A("width", img.Width).
		A("height", img.Height).
		A("title", title).
		N()
}
