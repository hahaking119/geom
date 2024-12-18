package geom_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/hahaking119/geom"
)

func TestLineStringZSetter(t *testing.T) {
	type tcase struct {
		points   [][3]float64
		setter   geom.LineStringZSetter
		expected geom.LineStringZSetter
		err      error
	}
	fn := func(t *testing.T, tc tcase) {
		err := tc.setter.SetVertices(tc.points)
		if tc.err == nil && err != nil {
			t.Errorf("error, expected nil got %v", err)
			return
		}
		if tc.err != nil {
			if tc.err.Error() != err.Error() {
				t.Errorf("error, expected %v got %v", tc.err, err)
			}
			return
		}
		// compare the results
		if !reflect.DeepEqual(tc.expected, tc.setter) {
			t.Errorf("setter, expected %v got %v", tc.expected, tc.setter)
		}
		lsm := tc.setter.Vertices()
		if !reflect.DeepEqual(tc.points, lsm) {
			t.Errorf("Vertices, expected %v got %v", tc.points, lsm)
		}
	}
	tests := []tcase{
		{
			points: [][3]float64{
				{15, 20, 30},
				{35, 40, 30},
				{-15, -5, 12},
			},
			setter: &geom.LineStringZ{
				{10, 20, 30},
				{30, 40, 30},
				{-10, -5, -2},
			},
			expected: &geom.LineStringZ{
				{15, 20, 30},
				{35, 40, 30},
				{-15, -5, 12},
			},
		},
		{
			setter: (*geom.LineStringZ)(nil),
			err:    geom.ErrNilLineStringZ,
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) { fn(t, tc) })
	}
}
