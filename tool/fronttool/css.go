// Tools collection to minify and prepare assets.
package fronttool

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"github.com/tdewolff/minify/v2/css"
)

// Merge and minify all **.css files from fsys.
// The paths are sorted before the merge.
//
// Replace all `(_[\w\d]+)` with value of map m.
//
// It panics on error.
func CSS(fsys fs.FS, m map[string]string) []byte {
	paths := make([]string, 0)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".css") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	sort.Strings(paths)

	r := regexp.MustCompile(`(_[\w\d]+)`)
	src := bytes.Buffer{}
	for _, p := range paths {
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			panic(err)
		}
		src.WriteString(r.ReplaceAllStringFunc(string(data), func(s string) string {
			out, ok := m[s]
			if !ok {
				panic(fmt.Sprintf("Not found %q in file %q", s, p))
			}
			return out
		}))
	}

	out := bytes.Buffer{}
	if err := css.Minify(nil, &out, &src, nil); err != nil {
		panic(err.Error())
	}

	return out.Bytes()
}

// FileSum hashs with sha256 the data.
//   - shortHash return is the hex of the hash.
//   - integrity is the integrity use as html attribute.
func FileSum(data []byte) (shortHash string, integrity string) {
	hash := sha256.Sum256(data)
	shortHash = hex.EncodeToString(hash[:4])
	integrity = "sha256-" + base64.StdEncoding.EncodeToString(hash[:])
	return
}
