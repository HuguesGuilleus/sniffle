//go:build !cgo

package resize0

import (
	"errors"
	"sniffle/tool"
)

func Resize(_ *tool.Tool, data []byte) ([]byte, error) {
	return nil, errors.New("need cgo")
}
