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
import (
	"fmt"
	"reflect"
	"unsafe"
)

// Fetch the pixel data type for this band
func (band *RasterBand) RasterDataType() DataType {
	return DataType(C.GDALGetRasterDataType(band.cval))
}

// Fetch the "natural" block size of this band
func (band *RasterBand) BlockSize() (int, int) {
	var xSize, ySize int
	C.GDALGetBlockSize(band.cval, (*C.int)(unsafe.Pointer(&xSize)), (*C.int)(unsafe.Pointer(&ySize)))
	return xSize, ySize
}

// Advise driver of upcoming read requests
func (band *RasterBand) AdviseRead(
	xOff, yOff, xSize, ySize, bufXSize, bufYSize int,
	dataType DataType,
	options []string,
) error {
	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALRasterAdviseRead(
		band.cval,
		C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize), C.int(bufXSize), C.int(bufYSize),
		C.GDALDataType(dataType),
		(**C.char)(unsafe.Pointer(&cOptions[0])),
	).Err()
}

// Read / Write a region of image data for this band
func (band *RasterBand) IO(
	rwFlag RWFlag,
	xOff, yOff, xSize, ySize int,
	buffer interface{},
	bufXSize, bufYSize int,
	pixelSpace, lineSpace int,
) error {
	var dataType DataType
	var dataPtr unsafe.Pointer
	switch data := buffer.(type) {
	case []int8:
		dataType = Byte
		dataPtr = unsafe.Pointer(&data[0])
	case []uint8:
		dataType = Byte
		dataPtr = unsafe.Pointer(&data[0])
	case []int16:
		dataType = Int16
		dataPtr = unsafe.Pointer(&data[0])
	case []uint16:
		dataType = UInt16
		dataPtr = unsafe.Pointer(&data[0])
	case []int32:
		dataType = Int32
		dataPtr = unsafe.Pointer(&data[0])
	case []uint32:
		dataType = UInt32
		dataPtr = unsafe.Pointer(&data[0])
	case []float32:
		dataType = Float32
		dataPtr = unsafe.Pointer(&data[0])
	case []float64:
		dataType = Float64
		dataPtr = unsafe.Pointer(&data[0])
	default:
		return fmt.Errorf("Error: buffer is not a valid data type (must be a valid numeric slice)")
	}

	return C.GDALRasterIO(
		band.cval,
		C.GDALRWFlag(rwFlag),
		C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize),
		dataPtr,
		C.int(bufXSize), C.int(bufYSize),
		C.GDALDataType(dataType),
		C.int(pixelSpace), C.int(lineSpace),
	).Err()
}

// Read a block of image data efficiently
func (band *RasterBand) ReadBlock(xOff, yOff int, dataPtr unsafe.Pointer) error {
	return C.GDALReadBlock(band.cval, C.int(xOff), C.int(yOff), dataPtr).Err()
}

// Write a block of image data efficiently
func (band *RasterBand) WriteBlock(xOff, yOff int, dataPtr unsafe.Pointer) error {
	return C.GDALWriteBlock(band.cval, C.int(xOff), C.int(yOff), dataPtr).Err()
}

// Fetch X size of raster
func (band *RasterBand) XSize() int {
	xSize := C.GDALGetRasterBandXSize(band.cval)
	return int(xSize)
}

// Fetch Y size of raster
func (band *RasterBand) YSize() int {
	ySize := C.GDALGetRasterBandYSize(band.cval)
	return int(ySize)
}

// Find out if we have update permission for this band
func (band *RasterBand) GetAccess() Access {
	access := C.GDALGetRasterAccess(band.cval)
	return Access(access)
}

// Fetch the band number of this raster band
func (band *RasterBand) Band() int {
	bandNumber := C.GDALGetBandNumber(band.cval)
	return int(bandNumber)
}

// Fetch the owning dataset handle
func (band *RasterBand) GetDataset() *Dataset {
	dataset := C.GDALGetBandDataset(band.cval)
	return &Dataset{dataset}
}

// How should this band be interpreted as color?
func (band *RasterBand) ColorInterp() ColorInterp {
	colorInterp := C.GDALGetRasterColorInterpretation(band.cval)
	return ColorInterp(colorInterp)
}

// Set color interpretation of the raster band
func (band *RasterBand) SetColorInterp(colorInterp ColorInterp) error {
	return C.GDALSetRasterColorInterpretation(band.cval, C.GDALColorInterp(colorInterp)).Err()
}

// Fetch the color table associated with this raster band
func (band *RasterBand) ColorTable() *ColorTable {
	ct := C.GDALGetRasterColorTable(band.cval)
	if ct == nil {
		return nil
	}
	return &ColorTable{ct}
}

// Set the raster color table for this raster band
func (band *RasterBand) SetColorTable(colorTable ColorTable) error {
	return C.GDALSetRasterColorTable(band.cval, colorTable.cval).Err()
}

// Check for arbitrary overviews
func (band *RasterBand) HasArbitraryOverviews() int {
	yes := C.GDALHasArbitraryOverviews(band.cval)
	return int(yes)
}

// Return the number of overview layers available
func (band *RasterBand) OverviewCount() int {
	count := C.GDALGetOverviewCount(band.cval)
	return int(count)
}

// Fetch overview raster band object
func (band *RasterBand) Overview(level int) *RasterBand {
	overview := C.GDALGetOverview(band.cval, C.int(level))
	if overview == nil {
		return nil
	}
	return &RasterBand{overview}
}

// Fetch the no data value for this band
func (band *RasterBand) NoDataValue() (val float64, valid bool) {
	var success int
	noDataVal := C.GDALGetRasterNoDataValue(band.cval, (*C.int)(unsafe.Pointer(&success)))
	return float64(noDataVal), success != 0
}

// Set the no data value for this band
func (band *RasterBand) SetNoDataValue(val float64) error {
	return C.GDALSetRasterNoDataValue(band.cval, C.double(val)).Err()
}

// Fetch the list of category names for this raster
func (band *RasterBand) CategoryNames() []string {
	p := C.GDALGetRasterCategoryNames(band.cval)
	var strings []string
	q := uintptr(unsafe.Pointer(p))
	for {
		p = (**C.char)(unsafe.Pointer(q))
		if *p == nil {
			break
		}
		strings = append(strings, C.GoString(*p))
		q += unsafe.Sizeof(q)
	}

	return strings
}

// Set the category names for this band
func (band *RasterBand) SetRasterCategoryNames(names []string) error {
	length := len(names)
	cStrings := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cStrings[i] = C.CString(names[i])
		defer C.free(unsafe.Pointer(cStrings[i]))
	}
	cStrings[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALSetRasterCategoryNames(band.cval, (**C.char)(unsafe.Pointer(&cStrings[0]))).Err()
}

// Fetch the minimum value for this band
func (band *RasterBand) GetMinimum() (val float64, valid bool) {
	var success int
	min := C.GDALGetRasterMinimum(band.cval, (*C.int)(unsafe.Pointer(&success)))
	return float64(min), success != 0
}

// Fetch the maximum value for this band
func (band *RasterBand) GetMaximum() (val float64, valid bool) {
	var success int
	max := C.GDALGetRasterMaximum(band.cval, (*C.int)(unsafe.Pointer(&success)))
	return float64(max), success != 0
}

// Fetch image statistics
func (band *RasterBand) GetStatistics(approxOK, force int) (min, max, mean, stdDev float64) {
	C.GDALGetRasterStatistics(
		band.cval,
		C.int(approxOK),
		C.int(force),
		(*C.double)(unsafe.Pointer(&min)),
		(*C.double)(unsafe.Pointer(&max)),
		(*C.double)(unsafe.Pointer(&mean)),
		(*C.double)(unsafe.Pointer(&stdDev)),
	)
	return min, max, mean, stdDev
}

// Compute image statistics
func (band *RasterBand) ComputeStatistics(
	approxOK int,
	progress ProgressFunc,
	data interface{},
) (min, max, mean, stdDev float64) {
	arg := &goGDALProgressFuncProxyArgs{progress, data}

	C.GDALComputeRasterStatistics(
		band.cval,
		C.int(approxOK),
		(*C.double)(unsafe.Pointer(&min)),
		(*C.double)(unsafe.Pointer(&max)),
		(*C.double)(unsafe.Pointer(&mean)),
		(*C.double)(unsafe.Pointer(&stdDev)),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return min, max, mean, stdDev
}

// Set statistics on raster band
func (band *RasterBand) SetStatistics(min, max, mean, stdDev float64) error {
	return C.GDALSetRasterStatistics(
		band.cval,
		C.double(min),
		C.double(max),
		C.double(mean),
		C.double(stdDev),
	).Err()
}

// Return raster unit type
func (band *RasterBand) GetUnitType() string {
	cString := C.GDALGetRasterUnitType(band.cval)
	return C.GoString(cString)
}

// Set unit type
func (band *RasterBand) SetUnitType(unit string) error {
	cString := C.CString(unit)
	defer C.free(unsafe.Pointer(cString))

	return C.GDALSetRasterUnitType(band.cval, cString).Err()
}

// Fetch the raster value offset
func (band *RasterBand) GetOffset() (float64, bool) {
	var success int
	val := C.GDALGetRasterOffset(band.cval, (*C.int)(unsafe.Pointer(&success)))
	return float64(val), success != 0
}

// Set scaling offset
func (band *RasterBand) SetOffset(offset float64) error {
	return C.GDALSetRasterOffset(band.cval, C.double(offset)).Err()
}

// Fetch the raster value scale
func (band *RasterBand) GetScale() (float64, bool) {
	var success int
	val := C.GDALGetRasterScale(band.cval, (*C.int)(unsafe.Pointer(&success)))
	return float64(val), success != 0
}

// Set scaling ratio
func (band *RasterBand) SetScale(scale float64) error {
	return C.GDALSetRasterScale(band.cval, C.double(scale)).Err()
}

// Compute the min / max values for a band
func (band *RasterBand) ComputeMinMax(approxOK int) (min, max float64) {
	var minmax [2]float64
	C.GDALComputeRasterMinMax(
		band.cval,
		C.int(approxOK),
		(*C.double)(unsafe.Pointer(&minmax[0])))
	return minmax[0], minmax[1]
}

// Flush raster data cache
func (band *RasterBand) FlushCache() {
	C.GDALFlushRasterCache(band.cval)
}

// Compute raster histogram
func (rb RasterBand) Histogram(
	min, max float64,
	buckets int,
	includeOutOfRange, approxOK int,
	progress ProgressFunc,
	data interface{},
) ([]uint64, error) {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	histogram := make([]C.GUIntBig, buckets)
	var err error
	if err = C.GDALGetRasterHistogramEx(
		rb.cval,
		C.double(min),
		C.double(max),
		C.int(buckets),
		(*C.GUIntBig)(unsafe.Pointer(&histogram[0])),
		C.int(includeOutOfRange),
		C.int(approxOK),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err(); err != nil {
		return nil, err
	} else {
		return CIntSliceToInt(histogram), nil
	}
	return nil, err
}

// Fetch default raster histogram
func (rb RasterBand) DefaultHistogram(
	force int,
	progress ProgressFunc,
	data interface{},
) (min, max float64, buckets int, histogram []uint64, err error) {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	var cHistogram *C.GUIntBig

	err = C.GDALGetDefaultHistogramEx(
		rb.cval,
		(*C.double)(&min),
		(*C.double)(&max),
		(*C.int)(unsafe.Pointer(&buckets)),
		&cHistogram,
		C.int(force),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()

	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&histogram))
	sliceHeader.Cap = buckets
	sliceHeader.Len = buckets
	sliceHeader.Data = uintptr(unsafe.Pointer(cHistogram))

	return min, max, buckets, histogram, err
}

// Set default raster histogram
// Unimplemented: SetDefaultHistogram

// Unimplemented: GetRandomRasterSample

// Fetch best sampling overviews
// Unimplemented: GetRasterSampleOverview

// Fill this band with a constant value
func (band *RasterBand) Fill(real, imaginary float64) error {
	return C.GDALFillRaster(band.cval, C.double(real), C.double(imaginary)).Err()
}

// Unimplemented: ComputeBandStats

// Unimplemented: OverviewMagnitudeCorrection

// Fetch default Raster Attribute Table
func (band *RasterBand) GetDefaultRAT() RasterAttributeTable {
	rat := C.GDALGetDefaultRAT(band.cval)
	return RasterAttributeTable{rat}
}

// Set default Raster Attribute Table
func (band *RasterBand) SetDefaultRAT(rat RasterAttributeTable) error {
	return C.GDALSetDefaultRAT(band.cval, rat.cval).Err()
}

// Unimplemented: AddDerivedBandPixelFunc

// Return the mask band associated with the band
func (band *RasterBand) GetMaskBand() *RasterBand {
	mask := C.GDALGetMaskBand(band.cval)
	if mask == nil {
		return nil
	}
	return &RasterBand{mask}
}

// Return the status flags of the mask band associated with the band
func (band *RasterBand) GetMaskFlags() int {
	flags := C.GDALGetMaskFlags(band.cval)
	return int(flags)
}

// Adds a mask band to the current band
func (band *RasterBand) CreateMaskBand(flags int) error {
	return C.GDALCreateMaskBand(band.cval, C.int(flags)).Err()
}

// Copy all raster band raster data
func (sourceRaster *RasterBand) RasterBandCopyWholeRaster(
	destRaster *RasterBand,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{progress, data}

	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALRasterBandCopyWholeRaster(
		sourceRaster.cval,
		destRaster.cval,
		(**C.char)(unsafe.Pointer(&cOptions[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}
