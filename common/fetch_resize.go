// Common operations for services.
package common

import (
	"bytes"
	"image"
	"sniffle/common/resize0"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"strconv"
)

type ResizedImage struct {
	Raw     Image
	Resized *Image
}

type Image struct {
	// File extension like ".png", ".jpg" ...
	Extension string
	// With and height as string because use will use it in attribute.
	Height string
	Width  string
	// Raw image data
	Data []byte `json:"-"`
}

func FetchImage(t *tool.Tool, request *fetch.Request) *ResizedImage {
	data := tool.FetchAll(t, request)
	if len(data) == 0 {
		return nil
	}

	img := Image{Data: data}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width == 0 || config.Height == 0 {
		return nil
	}
	switch format {
	case "png":
		img.Extension = ".png"
	case "jpeg":
		img.Extension = ".jpg"
	default:
		t.Warn("fetchImage", "err", "unknown format", "format", format, "url", request.URL.String())
		return nil
	}
	img.Height = strconv.Itoa(config.Height)
	img.Width = strconv.Itoa(config.Width)

	ri := &ResizedImage{img, nil}
	if resized := t.LongTask(resize0.Name, request.URL.String(), img.Data); len(resized) != 0 {
		width, height := resize0.NewDimension(config.Width, config.Height)
		ri.Resized = &Image{
			Extension: resize0.Extension,
			Width:     strconv.Itoa(width),
			Height:    strconv.Itoa(height),
			Data:      resized,
		}
	}

	return ri
}

// Render an image as <img.logo> or <picture.logo><source ...><img.logo></picture>.
// If img is nil, just return [render.Z].
func (img *ResizedImage) Render(base, title string) render.Node {
	if img == nil {
		return render.Z
	}

	raw := img.Raw
	imageElement := render.Na("img.logo", "loading", "lazy").
		A("src", base+raw.Extension).
		A("width", raw.Width).
		A("height", raw.Height).
		A("title", title).
		N()

	if res := img.Resized; res != nil {
		return render.N("picture.logo",
			render.Na("source", "type", resize0.MIME).
				A("srcset", base+res.Extension).
				N(),
			imageElement,
		)
	} else {
		return imageElement
	}
}
