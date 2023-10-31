package delaunay

import (
	pkg "github.com/hahaking119/geom/cmp"
)

var cmp = pkg.HiCMP

var oldCmp = pkg.SetDefault(pkg.HiCMP)
