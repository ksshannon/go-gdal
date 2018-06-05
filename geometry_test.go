// Copyright 2017 go-gdal. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdal

import "testing"

func TestGeomFromWKT(t *testing.T) {
	g, err := CreateFromWKT("POINT(-116.2261 43.6067)", SpatialReference{})
	if err != nil {
		t.Fatal(err)
	}
	if !g.IsValid() {
		t.Error("geom is invalid")
	}
	if g.Type() != GT_Point {
		t.Error("geom is not a point")
	}
	if g.PointCount() != 1 {
		t.Error("geom does not have 1 point")
	}
	if g.X(0) != -116.2261 || g.Y(0) != 43.6067 {
		t.Error("x or y invalid")
	}
}

func TestGeomRead(t *testing.T) {
	t.Skip("FIXME:(kyle)")
	fname := "test/poly.shp"
	ds, err := OpenEx(fname, VectorDrivers|ReadOnly, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	lyr, err := ds.Layer(0)
	if err != nil {
		t.Error(err)
	}
	feat := lyr.NextFeature()
	if feat == nil {
		t.Error("nil feature")
	}
	g := feat.Geometry()
	if g.Type() != GT_Polygon {
		t.Errorf("invalid geometry type: %d", g.Type())
	}
	if g.PointCount() < 1 {
		t.Errorf("point count: %d", g.PointCount())
	}
}
