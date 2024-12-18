package mvt

import (
	"testing"

	"github.com/hahaking119/geom"
	"github.com/hahaking119/geom/cmp"
)

func TestPrepareLinestring(t *testing.T) {

	type tcase struct {
		in   geom.LineString
		out  geom.LineString
		tile geom.Extent
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			got := preparelinestr(tc.in, &tc.tile, float64(DefaultExtent))

			if len(got) != len(tc.out) {
				t.Errorf("expected %v got %v", tc.out, got)
			}

			for i := range got {
				if !cmp.PointEqual(tc.out[i], got[i]) {
					t.Errorf("expected (%d) %v got %v", i, tc.out, got)
				}
			}
		}
	}

	tests := map[string]tcase{
		"duplicate pt simple line": {
			in:   geom.LineString{{9.0, 4090.0}, {9.0, 4090.0}},
			out:  geom.LineString{},
			tile: geom.Extent{0.0, 0.0, 4096.0, 4096.0},
		},
		"triplicate pt simple line": {
			in:   geom.LineString{{9.0, 4090.0}, {9.0, 4090.0}, {9.0, 4090.0}},
			out:  geom.LineString{},
			tile: geom.Extent{0.0, 0.0, 4096.0, 4096.0},
		},
		"triplicate pt simple line1": {
			in:   geom.LineString{{9.0, 4090.0}, {9.0, 4090.0}, {9.0, 4090.0}, {11.0, 4091.0}},
			out:  geom.LineString{{9.0, 6.0}, {11.0, 5.0}},
			tile: geom.Extent{0.0, 0.0, 4096.0, 4096.0},
		},
		"simple line": {
			in:   geom.LineString{{9.0, 4090.0}, {11.0, 4091.0}},
			out:  geom.LineString{{9.0, 6.0}, {11.0, 5.0}},
			tile: geom.Extent{0.0, 0.0, 4096.0, 4096.0},
		},
		"edge line": {
			in:   geom.LineString{{0.0, 0.0}, {4096.0, 20.0}},
			out:  geom.LineString{{0.0, 4096.0}, {4096.0, 4076.0}},
			tile: geom.Extent{0.0, 0.0, 4096.0, 4096.0},
		},
		"simple line 3pt": {
			in:   geom.LineString{{9.0, 4090.0}, {11.0, 4090.0}, {11.0, 4076.0}},
			out:  geom.LineString{{9.0, 6.0}, {11.0, 6.0}, {11.0, 20.0}},
			tile: geom.Extent{0.0, 0.0, 4096.0, 4096.0},
		},
		"scale": {
			in:   geom.LineString{{100.0, 100.0}, {300.0, 300.0}},
			out:  geom.LineString{{1024.0, 3072.0}, {3072.0, 1024.0}},
			tile: geom.Extent{0.0, 0.0, 400.0, 400.0},
		},
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}
