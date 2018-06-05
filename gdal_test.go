// Copyright 2011 go-gdal. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdal

import (
	"testing"
)

func TestTiffDriver(t *testing.T) {
	_, err := GetDriverByName("GTiff")
	if err != nil {
		t.Error(err)
	}
}

func TestShapeDriver(t *testing.T) {
	drv := OGRDriverByName("ESRI Shapefile")
	drv.TestCapability("NA")
}

func TestInvalidDriver(t *testing.T) {
	drv, err := GetDriverByName("FOO")
	if err == nil {
		t.Error("fetched invalid driver")
	}
	if drv != nil {
		_ = drv.ShortName()
	}
}

func TestInvalidCreate(t *testing.T) {
	drv, err := GetDriverByName("AAIGrid")
	if err != nil {
		t.Error(err)
	}
	PushQuietHandler()
	ds := drv.Create("/vsimem/foo.asc", 10, 10, 1, Float64, nil)
	PopHandler()
	if ds != nil {
		t.Error("created valid dataset without Create() support")
	}
}

func TestShortName(t *testing.T) {
	format := "GTiff"
	drv, err := GetDriverByName(format)
	if err != nil {
		t.Error(err)
	}
	if drv.ShortName() != format {
		t.Errorf("invalid short name: %s", drv.ShortName())
	}
}

func TestLongName(t *testing.T) {
	drv, err := GetDriverByName("GTiff")
	if err != nil {
		t.Error(err)
	}
	if drv.LongName() != "GeoTIFF" {
		t.Errorf("invalid long name: %s", drv.LongName())
	}
}

func TestOpen(t *testing.T) {
	ds, err := Open("test/small_world.tif", ReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	ds.Close()
}

func TestOpenEx(t *testing.T) {
	ds, err := OpenEx("test/small_world.tif",
		ReadOnly|RasterDrivers, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	ds.Close()
}

// Test open with flags we know should fail
func TestOpenExMismatchFlags(t *testing.T) {
	tests := []struct {
		fname   string
		flags   Access
		drivers []string
	}{
		{
			"test/small_world.tif",
			ReadOnly | RasterDrivers,
			[]string{"JPEG", "PNG"},
		},
		{
			"/vsicurl/http://download.osgeo.org/gdal/data/gtiff/small_world.tif",
			Update | RasterDrivers,
			nil,
		},
		{
			"test/small_world.tif",
			ReadOnly | VectorDrivers,
			nil,
		},
	}

	PushQuietHandler()
	defer PopHandler()
	for _, test := range tests {
		_, err := OpenEx(test.fname, test.flags, test.drivers, nil, nil)
		if err == nil {
			t.Errorf("driver flag test failed:%+v", test)
		}
	}
}

func TestRasterSize(t *testing.T) {
	ds, err := Open("test/small_world.tif", ReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	x := ds.RasterXSize()
	y := ds.RasterYSize()
	if x != 400 {
		t.Errorf("invalid x size: %d", x)
	}
	if y != 200 {
		t.Errorf("invalid y size: %d", y)
	}
}

func TestDataType(t *testing.T) {
	ds, err := Open("test/small_world.tif", ReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ds.RasterBand(1)
	if err != nil {
		t.Fatal(err)
	}
	dt := b.RasterDataType()
	if dt != Byte {
		t.Errorf("invalid data type: %s", dt.Name())
	}
}

func TestHistogram(t *testing.T) {
	drv, err := GetDriverByName("MEM")
	if err != nil {
		t.Error(err)
	}
	ds := drv.Create("/vsimem/tmp", 10, 10, 1, Byte, nil)
	defer ds.Close()
	band, err := ds.RasterBand(1)
	if err != nil {
		t.Error(err)
	}
	data := make([]uint8, 100)
	cs := 0
	for i := uint8(0); i < 100; i++ {
		data[i] = i
		cs += int(i)
	}
	err = band.IO(Write, 0, 0, 10, 10, data, 10, 10, 0, 0)
	if err != nil {
		t.Error(err)
	}
	band.FlushCache()
	hist, err := band.Histogram(0, 100, 10, 1, 0, DummyProgress, nil)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		if hist[i] != 10 {
			t.Errorf("failed to compute histogram. got: %+v, expected: 10\n", hist[i])
		}
	}
}

func TestInvalidBand(t *testing.T) {
	drv, err := GetDriverByName("MEM")
	if err != nil {
		t.Error(err)
	}
	ds := drv.Create("/vsimem/tmp", 10, 10, 1, Byte, nil)
	defer ds.Close()
	PushQuietHandler()
	band, err := ds.RasterBand(0)
	PopHandler()
	if err == nil {
		t.Error("failed to return error")
	}
	if band != nil && band.cval != nil {
		t.Error("C null deref")
	}
}

func TestGetLayer(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	fname := "test/poly.shp"
	ds, err := OpenEx(fname, ReadOnly|VectorDrivers, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ds.Layer(0)
	if err != nil {
		t.Error(err)
	}
	_, err = ds.LayerByName("poly")
	if err != nil {
		t.Error(err)
	}
}

func TestExecuteSQL(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	fname := "test/poly.shp"
	ds, err := OpenEx(fname, ReadOnly|VectorDrivers, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	lyr, err := ds.ExecuteSQL("SELECT * FROM poly", Geometry{}, "")
	if err != nil {
		t.Error(err)
	}
	if n, ok := lyr.FeatureCount(true); !ok || n < 1 {
		t.Error("failed to get a valid layer")
	}
}

func TestConfigOption(t *testing.T) {
	k, v := "GDAL_GO_TEST", "ON"
	SetConfigOption(k, v)
	if ConfigOption(k, "") != v {
		t.Errorf("Invalid value: %s\n", ConfigOption(k, ""))
	}
}

func BenchmarkRasterIO(b *testing.B) {
	ds, err := Open("test/small_world.tif", ReadOnly)
	if err != nil {
		b.Fatal(err)
	}
	nx := ds.RasterXSize()
	ny := ds.RasterYSize()
	buf := make([]byte, nx*ny)
	band, err := ds.RasterBand(1)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		band.IO(Read, 0, 0, nx, ny, buf, nx, ny, 0, 0)
	}
	ds.Close()
}
