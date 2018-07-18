package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"

#cgo linux  pkg-config: gdal
#cgo darwin pkg-config: gdal
#cgo windows LDFLAGS: -Lc:/gdal/release-1600-x64/lib -lgdal_i
#cgo windows CFLAGS: -IC:/gdal/release-1600-x64/include
*/
import "C"
import "unsafe"

type RATFieldType int

const (
	GFT_Integer = RATFieldType(C.GFT_Integer)
	GFT_Real    = RATFieldType(C.GFT_Real)
	GFT_String  = RATFieldType(C.GFT_String)
)

type RATFieldUsage int

const (
	GFU_Generic    = RATFieldUsage(C.GFU_Generic)
	GFU_PixelCount = RATFieldUsage(C.GFU_PixelCount)
	GFU_Name       = RATFieldUsage(C.GFU_Name)
	GFU_Min        = RATFieldUsage(C.GFU_Min)
	GFU_Max        = RATFieldUsage(C.GFU_Max)
	GFU_MinMax     = RATFieldUsage(C.GFU_MinMax)
	GFU_Red        = RATFieldUsage(C.GFU_Red)
	GFU_Green      = RATFieldUsage(C.GFU_Green)
	GFU_Blue       = RATFieldUsage(C.GFU_Blue)
	GFU_Alpha      = RATFieldUsage(C.GFU_Alpha)
	GFU_RedMin     = RATFieldUsage(C.GFU_RedMin)
	GFU_GreenMin   = RATFieldUsage(C.GFU_GreenMin)
	GFU_BlueMin    = RATFieldUsage(C.GFU_BlueMin)
	GFU_AlphaMin   = RATFieldUsage(C.GFU_AlphaMin)
	GFU_RedMax     = RATFieldUsage(C.GFU_RedMax)
	GFU_GreenMax   = RATFieldUsage(C.GFU_GreenMax)
	GFU_BlueMax    = RATFieldUsage(C.GFU_BlueMax)
	GFU_AlphaMax   = RATFieldUsage(C.GFU_AlphaMax)
	GFU_MaxCount   = RATFieldUsage(C.GFU_MaxCount)
)

// Construct empty raster attribute table
func CreateRasterAttributeTable() *RasterAttributeTable {
	rat := C.GDALCreateRasterAttributeTable()
	if rat == nil {
		return nil
	}
	return &RasterAttributeTable{rat}
}

// Destroy a RAT
func (rat *RasterAttributeTable) Destroy() {
	C.GDALDestroyRasterAttributeTable(rat.cval)
}

// Fetch table column count
func (rat *RasterAttributeTable) ColumnCount() int {
	count := C.GDALRATGetColumnCount(rat.cval)
	return int(count)
}

// Fetch the name of indicated column
func (rat *RasterAttributeTable) NameOfCol(index int) string {
	name := C.GDALRATGetNameOfCol(rat.cval, C.int(index))
	return C.GoString(name)
}

// Fetch the usage of indicated column
func (rat *RasterAttributeTable) UsageOfCol(index int) RATFieldUsage {
	rfu := C.GDALRATGetUsageOfCol(rat.cval, C.int(index))
	return RATFieldUsage(rfu)
}

// Fetch the type of indicated column
func (rat *RasterAttributeTable) TypeOfCol(index int) RATFieldType {
	rft := C.GDALRATGetTypeOfCol(rat.cval, C.int(index))
	return RATFieldType(rft)
}

// Fetch column index for indicated usage
func (rat *RasterAttributeTable) ColOfUsage(rfu RATFieldUsage) int {
	index := C.GDALRATGetColOfUsage(rat.cval, C.GDALRATFieldUsage(rfu))
	return int(index)
}

// Fetch row count
func (rat *RasterAttributeTable) RowCount() int {
	count := C.GDALRATGetRowCount(rat.cval)
	return int(count)
}

// Fetch field value as string
func (rat *RasterAttributeTable) ValueAsString(row, field int) string {
	cString := C.GDALRATGetValueAsString(rat.cval, C.int(row), C.int(field))
	return C.GoString(cString)
}

// Fetch field value as integer
func (rat *RasterAttributeTable) ValueAsInt(row, field int) int {
	val := C.GDALRATGetValueAsInt(rat.cval, C.int(row), C.int(field))
	return int(val)
}

// Fetch field value as float64
func (rat *RasterAttributeTable) ValueAsFloat64(row, field int) float64 {
	val := C.GDALRATGetValueAsDouble(rat.cval, C.int(row), C.int(field))
	return float64(val)
}

// Set field value from string
func (rat *RasterAttributeTable) SetValueAsString(row, field int, val string) {
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	C.GDALRATSetValueAsString(rat.cval, C.int(row), C.int(field), cVal)
}

// Set field value from integer
func (rat *RasterAttributeTable) SetValueAsInt(row, field, val int) {
	C.GDALRATSetValueAsInt(rat.cval, C.int(row), C.int(field), C.int(val))
}

// Set field value from float64
func (rat *RasterAttributeTable) SetValueAsFloat64(row, field int, val float64) {
	C.GDALRATSetValueAsDouble(rat.cval, C.int(row), C.int(field), C.double(val))
}

// Set row count
func (rat *RasterAttributeTable) SetRowCount(count int) {
	C.GDALRATSetRowCount(rat.cval, C.int(count))
}

// Create new column
func (rat *RasterAttributeTable) CreateColumn(name string, rft RATFieldType, rfu RATFieldUsage) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.GDALRATCreateColumn(rat.cval, cName, C.GDALRATFieldType(rft), C.GDALRATFieldUsage(rfu)).Err()
}

// Set linear binning information
func (rat *RasterAttributeTable) SetLinearBinning(row0min, binsize float64) error {
	return C.GDALRATSetLinearBinning(rat.cval, C.double(row0min), C.double(binsize)).Err()
}

// Fetch linear binning information
func (rat *RasterAttributeTable) LinearBinning() (row0min, binsize float64, exists bool) {
	success := C.GDALRATGetLinearBinning(rat.cval, (*C.double)(&row0min), (*C.double)(&binsize))
	return row0min, binsize, success != 0
}

// Initialize RAT from color table
func (rat *RasterAttributeTable) FromColorTable(ct ColorTable) error {
	return C.GDALRATInitializeFromColorTable(rat.cval, ct.cval).Err()
}

// Translate RAT to a color table
func (rat *RasterAttributeTable) ToColorTable(count int) ColorTable {
	ct := C.GDALRATTranslateToColorTable(rat.cval, C.int(count))
	return ColorTable{ct}
}

// Dump RAT in readable form to a file
// Unimplemented: DumpReadable

// Get row for pixel value
func (rat *RasterAttributeTable) RowOfValue(val float64) (int, bool) {
	row := C.GDALRATGetRowOfValue(rat.cval, C.double(val))
	return int(row), row != -1
}
