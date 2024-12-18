package subdivision

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/hahaking119/geom"
	"github.com/hahaking119/geom/encoding/wkt"
	"github.com/hahaking119/geom/internal/debugger"
	"github.com/hahaking119/geom/planar/triangulate/delaunay/quadedge"
)

/*
const (
	debug = false
)
*/

// For when we have debug as a var
var debug = false

func ToggleDebug() {
	debug = !debug
}

func init() {
	debugger.DefaultOutputDir = "output"
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Debug %v", debug)
}

// ErrAssumptionFailed is an assert of when our assumptions fails, in debug mode it will return and error. In
// non debug mode it will panic
func ErrAssumptionFailed() error {
	str := fmt.Sprintf("Assumption failed at: %v ", debugger.FFL(0))
	if debug {
		return errors.New(str)
	}
	panic(str)
}

// DumpSubdivision will print each edge in the subdivision
func DumpSubdivision(sd *Subdivision) {
	var str strings.Builder
	DumpSubdivisionW(&str, sd)
	log.Print(str.String())
}

// DumpSubdivisionW will write each edge in the subdivision to w
func DumpSubdivisionW(w io.Writer, sd *Subdivision) {
	fmt.Fprintf(w, "Frame: %#v\n", sd.frame)
	var edges geom.MultiLineString

	_ = sd.WalkAllEdges(func(e *quadedge.Edge) error {
		/*
			if IsFrameEdge(sd.frame,e) {
				return nil
			}
		*/
		ln := e.AsLine()
		edges = append(edges, ln[:])
		return nil
	})
	fmt.Fprintf(w, "Edges:\n%v", wkt.MustEncode(edges))
}
