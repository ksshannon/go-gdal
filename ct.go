package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
#include "cpl_http.h"

#cgo linux  pkg-config: gdal
#cgo darwin pkg-config: gdal
#cgo windows LDFLAGS: -Lc:/gdal/release-1600-x64/lib -lgdal_i
#cgo windows CFLAGS: -IC:/gdal/release-1600-x64/include
*/
import "C"

/* ==================================================================== */
/*      Color tables.                                                   */
/* ==================================================================== */

// Types of color interpretations for a GDALColorTable.
type PaletteInterp int

const (
	// Grayscale (in GDALColorEntry.c1)
	PI_Gray = PaletteInterp(C.GPI_Gray)
	// Red, Green, Blue and Alpha in (in c1, c2, c3 and c4)
	PI_RGB = PaletteInterp(C.GPI_RGB)
	// Cyan, Magenta, Yellow and Black (in c1, c2, c3 and c4)
	PI_CMYK = PaletteInterp(C.GPI_CMYK)
	// Hue, Lightness and Saturation (in c1, c2, and c3)
	PI_HLS = PaletteInterp(C.GPI_HLS)
)

type ColorTable struct {
	cval C.GDALColorTableH
}

type ColorEntry struct {
	cval *C.GDALColorEntry
}

func (paletteInterp PaletteInterp) Name() string {
	return C.GoString(C.GDALGetPaletteInterpretationName(C.GDALPaletteInterp(paletteInterp)))
}

// Construct a new color table
func CreateColorTable(interp PaletteInterp) *ColorTable {
	ct := C.GDALCreateColorTable(C.GDALPaletteInterp(interp))
	if ct == nil {
		return nil
	}
	return &ColorTable{ct}
}

// Destroy the color table
func (ct *ColorTable) Destroy() {
	C.GDALDestroyColorTable(ct.cval)
}

// Make a copy of the color table
func (ct *ColorTable) Clone() *ColorTable {
	newCT := C.GDALCloneColorTable(ct.cval)
	if newCT == nil {
		return nil
	}
	return &ColorTable{newCT}
}

// Fetch palette interpretation
func (ct *ColorTable) PaletteInterpretation() PaletteInterp {
	pi := C.GDALGetPaletteInterpretation(ct.cval)
	return PaletteInterp(pi)
}

// Get number of color entries in table
func (ct *ColorTable) EntryCount() int {
	count := C.GDALGetColorEntryCount(ct.cval)
	return int(count)
}

// Fetch a color entry from table
func (ct *ColorTable) Entry(index int) *ColorEntry {
	entry := C.GDALGetColorEntry(ct.cval, C.int(index))
	if entry == nil {
		return nil
	}
	return &ColorEntry{entry}
}

// Unimplemented: EntryAsRGB

// Set entry in color table
func (ct *ColorTable) SetEntry(index int, entry *ColorEntry) {
	C.GDALSetColorEntry(ct.cval, C.int(index), entry.cval)
}

// Create color ramp
func (ct *ColorTable) CreateColorRamp(start, end int, startColor, endColor *ColorEntry) {
	C.GDALCreateColorRamp(ct.cval, C.int(start), startColor.cval, C.int(end), endColor.cval)
}
