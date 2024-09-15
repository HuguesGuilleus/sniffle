//go:build !cgo

package resize0

import "errors"

func Resize(data []byte) ([]byte, error) {
	return nil, errors.New("need cgo")
}
