package gdal

import "testing"

func getLayer(t *testing.T) *Layer {
	fname := "test/poly.shp"
	ds, err := OpenEx(fname, VectorDrivers|ReadOnly, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	lyr, err := ds.Layer(0)
	if err != nil {
		t.Fatal(err)
	}
	return &lyr
}

func TestLayerName(t *testing.T) {
	lyr := getLayer(t)
	if lyr.Name() != "poly" {
		t.Errorf("invalid layer name: %s", lyr.Name())
	}
}

func TestLayerGeomType(t *testing.T) {
	lyr := getLayer(t)
	if lyr.GeomType() != GT_Polygon {
		t.Errorf("invalid geometry type: %+v", lyr.GeomType())
	}
}

func TestFeatureCount(t *testing.T) {
	lyr := getLayer(t)
	if n, ok := lyr.FeatureCount(true); !ok || n != 10 {
		t.Errorf("invalid feature count: %d", n)
	}
}

func TestNextFeature(t *testing.T) {
	lyr := getLayer(t)
	n, _ := lyr.FeatureCount(true)
	i := 0
	for feat := lyr.NextFeature(); feat != nil; feat = lyr.NextFeature() {
		i++
	}
	if i != n {
		t.Errorf("did not get %d features, only %d", n, i)
	}
}

func TestSpatialFilters(t *testing.T) {
	// Extent: (478315.531250, 4762880.500000) - (481645.312500, 4765610.500000)
	lyr := getLayer(t)
	if lyr.SpatialFilter() != nil {
		t.Errorf("spatial filter not nil")
	}
	// Create a spatial filter that has no features
	geom, err := CreateFromWKT("POLYGON((0 0, 1 0, 1 1, 0 1, 0 0))", *lyr.SpatialRef())
	if err != nil {
		t.Fatal("could not create filter geometry")
	}
	lyr.SetSpatialFilter(geom)
	lyr.ResetReading()
	if lyr.SpatialFilter() == nil {
		t.Error("spatial filter not set")
	}
	if n, ok := lyr.FeatureCount(true); ok != true || n != 0 {
		t.Error("spatial filter not set")
	}
	lyr.SetSpatialFilter(nil)
	lyr.ResetReading()
	if lyr.SpatialFilter() != nil {
		t.Error("spatial filter set")
	}
	if n, ok := lyr.FeatureCount(true); ok != true || n != 10 {
		t.Error("spatial filter set")
	}
	lyr.SetSpatialFilterRect(478903, 4764757, 478904, 4764756)
	lyr.ResetReading()
	if lyr.SpatialFilter() == nil {
		t.Error("spatial filter not set")
	}
	n, ok := lyr.FeatureCount(true)
	if ok != true {
		t.Error("failed to get feature count")
	}
	if n != 1 {
		t.Error("spatial filter does not include expected feature")
	}
	feat := lyr.NextFeature()
	if feat == nil || feat.FieldAsInteger(feat.FieldIndex("EAS_ID")) != 173 {
		t.Error("failed to get feature with proper fid")
	}
}

func TestAttributeFilter(t *testing.T) {
	lyr := getLayer(t)
	lyr.SetAttributeFilter("EAS_ID=173")
	lyr.ResetReading()
	if n, ok := lyr.FeatureCount(true); !ok || n != 1 {
		t.Error("failed to set attribute filter")
	}
}

func TestResetReading(t *testing.T) {
	lyr := getLayer(t)
	lyr.NextFeature()
	lyr.ResetReading()
	n, _ := lyr.FeatureCount(true)
	i := 0
	for feat := lyr.NextFeature(); feat != nil; feat = lyr.NextFeature() {
		i++
	}
	if i != n {
		t.Error("failed to reset reading")
	}
}

func TestIgnoreFields(t *testing.T) {
	t.Skip("find a format that supports ignoring fields")
	lyr := getLayer(t)
	err := lyr.SetIgnoredFields([]string{"EAS_ID"})
	if err != nil {
		t.Fatal(err)
	}
	feat := lyr.NextFeature()
	if feat.FieldIndex("EAS_ID") != -1 {
		t.Error("ignored fields not ignored")
		t.Log(feat.FieldIndex("EAS_ID"))
	}
}
