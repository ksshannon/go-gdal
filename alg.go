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
import (
	"fmt"
	"unsafe"
)

// Compute optimal PCT for RGB image
func ComputeMedianCutPCT(
	red, green, blue *RasterBand,
	colors int,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	err := C.GDALComputeMedianCutPCT(
		red.cval,
		green.cval,
		blue.cval,
		nil,
		C.int(colors),
		ct.cval,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return int(err)
}

// 24bit to 8bit conversion with dithering
func DitherRGB2PCT(
	red, green, blue, target *RasterBand,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	err := C.GDALDitherRGB2PCT(
		red.cval,
		green.cval,
		blue.cval,
		target.cval,
		ct.cval,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return int(err)
}

// Compute checksum for image region
func (rb *RasterBand) Checksum(xOff, yOff, xSize, ySize int) int {
	sum := C.GDALChecksumImage(rb.cval, C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize))
	return int(sum)
}

// Compute the proximity of all pixels in the image to a set of pixels in the source image
func (src *RasterBand) ComputeProximity(
	dest *RasterBand,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALComputeProximity(
		src.cval,
		dest.cval,
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Fill selected raster regions by interpolation from the edges
func (src *RasterBand) FillNoData(
	mask *RasterBand,
	distance float64,
	iterations int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALFillNodata(
		src.cval,
		mask.cval,
		C.double(distance),
		0,
		C.int(iterations),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Create polygon coverage from raster data using an integer buffer
func (src *RasterBand) Polygonize(
	mask *RasterBand,
	layer *Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALPolygonize(
		src.cval,
		mask.cval,
		layer.cval,
		C.int(fieldIndex),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Create polygon coverage from raster data using a floating point buffer
func (src *RasterBand) FPolygonize(
	mask *RasterBand,
	layer *Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALFPolygonize(
		src.cval,
		mask.cval,
		layer.cval,
		C.int(fieldIndex),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Removes small raster polygons
func (src *RasterBand) SieveFilter(
	mask, dest *RasterBand,
	threshold, connectedness int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALSieveFilter(
		src.cval,
		mask.cval,
		dest.cval,
		C.int(threshold),
		C.int(connectedness),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

/* --------------------------------------------- */
/* Warp functions                                */
/* --------------------------------------------- */

//Unimplemented: CreateGenImgProjTransformer
//Unimplemented: CreateGenImgProjTransformer2
//Unimplemented: CreateGenImgProjTransformer3
//Unimplemented: SetGenImgProjTransformerDstGeoTransform
//Unimplemented: DestroyGenImgProjTransformer
//Unimplemented: GenImgProjTransform

//Unimplemented: CreateReprojectionTransformer
//Unimplemented: DestroyReprojection
//Unimplemented: ReprojectionTransform
//Unimplemented: CreateGCPTransformer
//Unimplemented: CreateGCPRefineTransformer
//Unimplemented: DestroyGCPTransformer
//Unimplemented: GCPTransform

//Unimplemented: CreateTPSTransformer
//Unimplemented: DestroyTPSTransformer
//Unimplemented: TPSTransform

//Unimplemented: CreateRPCTransformer
//Unimplemented: DestroyRPCTransformer
//Unimplemented: RPCTransform

//Unimplemented: CreateGeoLocTransformer
//Unimplemented: DestroyGeoLocTransformer
//Unimplemented: GeoLocTransform

//Unimplemented: CreateApproxTransformer
//Unimplemented: DestroyApproxTransformer
//Unimplemented: ApproxTransform

//Unimplemented: SimpleImageWarp
//Unimplemented: SuggestedWarpOutput
//Unimplemented: SuggsetedWarpOutput2
//Unimplemented: SerializeTransformer
//Unimplemented: DeserializeTransformer

//Unimplemented: TransformGeolocations

/* --------------------------------------------- */
/* Contour line functions                        */
/* --------------------------------------------- */

//Unimplemented: CreateContourGenerator
//Unimplemented: FeedLine
//Unimplemented: Destroy
//Unimplemented: ContourWriter
//Unimplemented: ContourGenerate

/* --------------------------------------------- */
/* Rasterizer functions                          */
/* --------------------------------------------- */

// Burn geometries into raster
//Unimplmemented: RasterizeGeometries

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayers

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayersBuf

/* --------------------------------------------- */
/* Gridding functions                            */
/* --------------------------------------------- */

// TODO(kyle): this is not right
type GridAlgorithm uint8

const (
	InverseDistanceToAPower                = GridAlgorithm(C.GGA_InverseDistanceToAPower)
	MovingAverage                          = GridAlgorithm(C.GGA_MovingAverage)
	NearestNeighbor                        = GridAlgorithm(C.GGA_NearestNeighbor)
	MetricMinimum                          = GridAlgorithm(C.GGA_MetricMinimum)
	MetricMaximum                          = GridAlgorithm(C.GGA_MetricMaximum)
	MetricRange                            = GridAlgorithm(C.GGA_MetricRange)
	MetricCount                            = GridAlgorithm(C.GGA_MetricCount)
	MetricAverageDistance                  = GridAlgorithm(C.GGA_MetricAverageDistance)
	MetricAverageDistancePts               = GridAlgorithm(C.GGA_MetricAverageDistancePts)
	Linear                                 = GridAlgorithm(C.GGA_Linear)
	InverseDistanceToAPowerNearestNeighbor = GridAlgorithm(C.GGA_InverseDistanceToAPowerNearestNeighbor)
)

type GridOptions struct {
	Power           float64
	Smoothing       float64
	AnisotropyRatio float64
	AnisotropyAngle float64
	Radius1         float64
	Radius2         float64
	Angle           float64
	MaxPoints       uint32
	MinPoints       uint32
	NoDataValue     float64
}

// Fill selected raster regions by interpolation from the edges
func Grid(
	algorithm GridAlgorithm,
	options interface{}, // maybe define an interface.
	x, y, z []float64,
	xmin, xmax float64,
	ymin, ymax float64,
	xsize, ysize uint32,
	buffer interface{},
	progress ProgressFunc,
	data interface{},
) error {

	var dataPtr unsafe.Pointer
	var dataType DataType
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

	return C.GDALGridCreate(
		C.GDALGridAlgorithm(algorithm),
		unsafe.Pointer(&options),
		C.GUInt32(len(x)),
		(*C.double)(&x[0]),
		(*C.double)(&y[0]),
		(*C.double)(&z[0]),
		C.double(xmin), C.double(xmax),
		C.double(ymin), C.double(ymax),
		C.GUInt32(xsize),
		C.GUInt32(ysize),
		C.GDALDataType(dataType),
		dataPtr,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(dataPtr),
	).Err()
}

//Unimplemented: CreateGrid
//Unimplemented: ComputeMatchingPoints
