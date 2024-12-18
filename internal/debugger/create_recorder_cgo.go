//go:build cgo
// +build cgo

package debugger

import (
	"fmt"

	rcdr "github.com/hahaking119/geom/internal/debugger/recorder"
	"github.com/hahaking119/geom/internal/debugger/recorder/gpkg"
)

func NewRecorder(dir, filename string) (rcdr.Interface, string, error) {
	r, fn, err := gpkg.New(dir, filename, 0)
	if err != nil {
		return nil, fn, fmt.Errorf("gpkg error: %v", err)
	}
	return r, fn, nil
}
