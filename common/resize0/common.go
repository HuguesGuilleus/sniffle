// resize0 is long task function to resize and encode image in AVIF.
package resize0

const Name = "resize0"
const Extension = ".avif"
const MIME = "image/avif"
const Size uint = 200

// NewDimension calculate the width and heigh max by [Size].
func NewDimension(width, height int) (int, int) {
	const size = int(Size)

	if size >= width && size >= height {
		return width, height
	}

	if width >= size {
		height = max(1, size*height/width)
		width = size
	}

	if height >= size {
		width = max(1, size*width/height)
		height = size
	}

	return width, height
}
