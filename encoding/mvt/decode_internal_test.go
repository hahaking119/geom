package mvt

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hahaking119/geom"
	"github.com/hahaking119/geom/cmp"
	vectorTile "github.com/hahaking119/geom/encoding/mvt/vector_tile"
)

func TestDecode(t *testing.T) {
	//file := "/Users/bytedance/Downloads/tile(432016,197887,19).mvt"
	//file = "/Users/bytedance/Downloads/tile(432016,197888,19).mvt"
	//file = "/Users/bytedance/Downloads/tile(432017,197887,19).mvt"
	//file = "/Users/bytedance/Downloads/tile(432017,197888,19).mvt"
	files := []string{"/Users/bytedance/Downloads/tile(432016,197887,19).mvt",
		"/Users/bytedance/Downloads/tile(432016,197888,19).mvt",
		"/Users/bytedance/Downloads/tile(432017,197887,19).mvt",
		"/Users/bytedance/Downloads/tile(432017,197888,19).mvt",
		"/Users/bytedance/Downloads/tile(428737,202157,19).mvt",
		"/Users/bytedance/Downloads/tile(428738,202156,19).mvt",
		"/Users/bytedance/Downloads/tile(428738,202157,19).mvt",
	}
	for _, file := range files {
		f, err := os.Open(file)
		if err == nil {
			ret, err := Decode(f)
			if err == nil {
				for _, layer := range ret.layers {
					if layer.Name != "speeds" {
						continue
					}
					for _, feature := range layer.features {
						//if _, ok := feature.Tags["name"]; !ok {
						//	continue
						//}
						//name := feature.Tags["name"].(string)
						//fmt.Print(name)
						//if !strings.Contains(name, "147024753") || !strings.Contains(name, "147019978") {
						//	continue
						//}
						fmt.Println(strconv.FormatFloat(feature.Tags["weight"].(float64), 'f', 15, 64))
					}
				}
			}
		}
	}

	type tcase struct {
		typ vectorTile.Tile_GeomType
		buf []uint32
		geo geom.Geometry
		err error
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			geo, err := DecodeGeometry(tc.typ, tc.buf)
			if err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error %v", err)
				} else if tc.err.Error() == err.Error() {
					t.Errorf("unexpeced error %v, expected %v",
						err, tc.err)
				}
				return
			}

			if !cmp.GeometryEqual(geo, tc.geo) {
				t.Errorf("incorrect geometry, expected\n\t%v\ngot\n\t%v", tc.geo, geo)
			}
		}
	}

	// TODO(ear7h) error test cases/fuzzing
	testcases := map[string]tcase{
		"point": {
			typ: vectorTile.Tile_POINT,
			buf: []uint32{9, 50, 34},
			geo: geom.Point{25, 17},
		},
		"multi point": {
			typ: vectorTile.Tile_POINT,
			buf: []uint32{17, 10, 14, 3, 9},
			geo: geom.MultiPoint{{5, 7}, {3, 2}},
		},
		"line string": {
			typ: vectorTile.Tile_LINESTRING,
			buf: []uint32{9, 4, 4, 18, 0, 16, 16, 0},
			geo: geom.LineString{{2, 2}, {2, 10}, {10, 10}},
		},
		"multi line string": {
			typ: vectorTile.Tile_LINESTRING,
			buf: []uint32{9, 4, 4, 18, 0, 16, 16, 0, 9, 17, 17, 10, 4, 8},
			geo: geom.MultiLineString{{{2, 2}, {2, 10}, {10, 10}}, {{1, 1}, {3, 5}}},
		},
		"polygon": {
			typ: vectorTile.Tile_POLYGON,
			buf: []uint32{9, 6, 12, 18, 10, 12, 24, 44, 15},
			geo: geom.Polygon{{{3, 6}, {8, 12}, {20, 34}}},
		},
		"multi polygon": {
			typ: vectorTile.Tile_POLYGON,
			buf: []uint32{9, 0, 0, 26, 20, 0, 0, 20, 19, 0, 15, 9, 22, 2, 26, 18, 0, 0, 18, 17, 0, 15, 9, 4, 13, 26, 0, 8, 8, 0, 0, 7, 15},
			geo: geom.MultiPolygon{
				{ // poly 1
					{{0, 0}, {10, 0}, {10, 10}, {0, 10}},
				},
				{ // poly 2
					{{11, 11}, {20, 11}, {20, 20}, {11, 20}},
					{{13, 13}, {13, 17}, {17, 17}, {17, 13}},
				},
			},
		},
	}

	for k, v := range testcases {
		t.Run(k, fn(v))
	}
}
