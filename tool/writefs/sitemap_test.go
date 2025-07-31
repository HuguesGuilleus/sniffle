package writefs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSitemap(t *testing.T) {
	m := Memory()
	mcreator := Creator(m)
	sitemap := Sitemap(&mcreator)

	// Create
	w, err := sitemap.Create("/a.html")
	assert.NoError(t, err)
	_, err = w.Write([]byte("some html ..."))
	assert.NoError(t, err)
	assert.NoError(t, w.Close())

	data, err := ReadAll(m, "/a.html")
	assert.NoError(t, err)
	assert.Equal(t, "some html ...", string(data))

	// Simple write
	assert.NoError(t, sitemap.WriteFile("/b.html", []byte("other html")))

	data, err = ReadAll(m, "/b.html")
	assert.NoError(t, err)
	assert.Equal(t, "other html", string(data))

	// Output
	assert.Equal(t, ""+
		"https://example/a.html\n"+
		"https://example/b.html\n",
		string(sitemap.Sitemap("https://example")),
	)
}
