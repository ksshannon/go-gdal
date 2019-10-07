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
	"errors"
	"fmt"
	"unsafe"
)

func init() {
	C.GDALAllRegister()
}

const (
	VERSION_MAJOR = int(C.GDAL_VERSION_MAJOR)
	VERSION_MINOR = int(C.GDAL_VERSION_MINOR)
	VERSION_REV   = int(C.GDAL_VERSION_REV)
	VERSION_BUILD = int(C.GDAL_VERSION_BUILD)
	VERSION_NUM   = int(C.GDAL_VERSION_NUM)
	RELEASE_DATE  = int(C.GDAL_RELEASE_DATE)
	RELEASE_NAME  = string(C.GDAL_RELEASE_NAME)
)

var (
	ErrDebug   = errors.New("Debug Error")
	ErrWarning = errors.New("Warning Error")
	ErrFailure = errors.New("Failure Error")
	ErrFatal   = errors.New("Fatal Error")
	ErrIllegal = errors.New("Illegal Error")
)

// Error handling.  The following is bare-bones, and needs to be replaced with something more useful.
func (err C.CPLErr) Err() error {
	switch err {
	case 0:
		return nil
	case 1:
		return ErrDebug
	case 2:
		return ErrWarning
	case 3:
		return ErrFailure
	case 4:
		return ErrFailure
	}
	return ErrIllegal
}

func (err C.OGRErr) Err() error {
	switch err {
	case 0:
		return nil
	case 1:
		return ErrDebug
	case 2:
		return ErrWarning
	case 3:
		return ErrFailure
	case 4:
		return ErrFailure
	}
	return ErrIllegal
}

// Pixel data types
type DataType int

const (
	Unknown  = DataType(C.GDT_Unknown)
	Byte     = DataType(C.GDT_Byte)
	UInt16   = DataType(C.GDT_UInt16)
	Int16    = DataType(C.GDT_Int16)
	UInt32   = DataType(C.GDT_UInt32)
	Int32    = DataType(C.GDT_Int32)
	Float32  = DataType(C.GDT_Float32)
	Float64  = DataType(C.GDT_Float64)
	CInt16   = DataType(C.GDT_CInt16)
	CInt32   = DataType(C.GDT_CInt32)
	CFloat32 = DataType(C.GDT_CFloat32)
	CFloat64 = DataType(C.GDT_CFloat64)
)

// Get data type size in bits.
func (dataType DataType) Size() int {
	return int(C.GDALGetDataTypeSize(C.GDALDataType(dataType)))
}

func (dataType DataType) IsComplex() int {
	return int(C.GDALDataTypeIsComplex(C.GDALDataType(dataType)))
}

func (dataType DataType) Name() string {
	return C.GoString(C.GDALGetDataTypeName(C.GDALDataType(dataType)))
}

func (dataType DataType) Union(dataTypeB DataType) DataType {
	return DataType(
		C.GDALDataTypeUnion(C.GDALDataType(dataType), C.GDALDataType(dataTypeB)),
	)
}

//Safe array conversion
func IntSliceToCInt(data []int) []C.int {
	sliceSz := len(data)
	result := make([]C.int, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = C.int(data[i])
	}
	return result
}

//Safe array conversion
func CIntSliceToInt(data []C.GUIntBig) []uint64 {
	sliceSz := len(data)
	result := make([]uint64, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = uint64(data[i])
	}
	return result
}

// status of the asynchronous stream
type AsyncStatusType int

const (
	AR_Pending  = AsyncStatusType(C.GARIO_PENDING)
	AR_Update   = AsyncStatusType(C.GARIO_UPDATE)
	AR_Error    = AsyncStatusType(C.GARIO_ERROR)
	AR_Complete = AsyncStatusType(C.GARIO_COMPLETE)
)

func (statusType AsyncStatusType) Name() string {
	return C.GoString(C.GDALGetAsyncStatusTypeName(C.GDALAsyncStatusType(statusType)))
}

func GetAsyncStatusTypeByName(statusTypeName string) AsyncStatusType {
	name := C.CString(statusTypeName)
	defer C.free(unsafe.Pointer(name))
	return AsyncStatusType(C.GDALGetAsyncStatusTypeByName(name))
}

// Flag indicating read/write, or read-only access to data.
type Access uint

const (
	// Read only (no update) access
	ReadOnly = Access(C.GA_ReadOnly)
	// Read/write access.
	Update = Access(C.GA_Update)

	// GDAL_OF flags for opening datasets with OpenEx()
	// Note: we define GDAL_OF_READONLY and GDAL_OF_UPDATE to be on purpose
	// equals to GA_ReadOnly and GA_Update defined above.

	// Allow raster and vector drivers to be used.
	AllDrivers = Access(C.GDAL_OF_ALL)

	// Allow raster drivers to be used.
	RasterDrivers = Access(C.GDAL_OF_RASTER)

	// Allow vector drivers to be used.
	VectorDrivers = Access(C.GDAL_OF_VECTOR)

	// Allow gnm drivers to be used.
	GNMDrivers = Access(C.GDAL_OF_GNM)

	// Unsure
	KindMask = Access(C.GDAL_OF_KIND_MASK)

	// Open in shared mode.
	Shared = Access(C.GDAL_OF_SHARED)

	// Emit error message in case of failed open.
	Verbose = Access(C.GDAL_OF_VERBOSE_ERROR)

	// Open as internal dataset. Such dataset isn't registered in the global list
	// of opened dataset. Cannot be used with GDAL_OF_SHARED.
	Internal = Access(C.GDAL_OF_INTERNAL)

// 2.1+ flags not supported yet
)

// Read/Write flag for RasterIO() method
type RWFlag int

const (
	// Read data
	Read = RWFlag(C.GF_Read)
	// Write data
	Write = RWFlag(C.GF_Write)
)

// Types of color interpretation for raster bands.
type ColorInterp int

const (
	CI_Undefined      = ColorInterp(C.GCI_Undefined)
	CI_GrayIndex      = ColorInterp(C.GCI_GrayIndex)
	CI_PaletteIndex   = ColorInterp(C.GCI_PaletteIndex)
	CI_RedBand        = ColorInterp(C.GCI_RedBand)
	CI_GreenBand      = ColorInterp(C.GCI_GreenBand)
	CI_BlueBand       = ColorInterp(C.GCI_BlueBand)
	CI_AlphaBand      = ColorInterp(C.GCI_AlphaBand)
	CI_HueBand        = ColorInterp(C.GCI_HueBand)
	CI_SaturationBand = ColorInterp(C.GCI_SaturationBand)
	CI_LightnessBand  = ColorInterp(C.GCI_LightnessBand)
	CI_CyanBand       = ColorInterp(C.GCI_CyanBand)
	CI_MagentaBand    = ColorInterp(C.GCI_MagentaBand)
	CI_YellowBand     = ColorInterp(C.GCI_YellowBand)
	CI_BlackBand      = ColorInterp(C.GCI_BlackBand)
	CI_YCbCr_YBand    = ColorInterp(C.GCI_YCbCr_YBand)
	CI_YCbCr_CbBand   = ColorInterp(C.GCI_YCbCr_CbBand)
	CI_YCbCr_CrBand   = ColorInterp(C.GCI_YCbCr_CrBand)
	CI_Max            = ColorInterp(C.GCI_Max)
)

func (colorInterp ColorInterp) Name() string {
	return C.GoString(C.GDALGetColorInterpretationName(C.GDALColorInterp(colorInterp)))
}

func GetColorInterpretationByName(name string) ColorInterp {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ColorInterp(C.GDALGetColorInterpretationByName(cName))
}

// "well known" metadata items.
const (
	MD_AREA_OR_POINT = string(C.GDALMD_AREA_OR_POINT)
	MD_AOP_AREA      = string(C.GDALMD_AOP_AREA)
	MD_AOP_POINT     = string(C.GDALMD_AOP_POINT)
)

/* -------------------------------------------------------------------- */
/*      Define handle types related to various internal classes.        */
/* -------------------------------------------------------------------- */

type MajorObject struct {
	cval C.GDALMajorObjectH
}

type Dataset struct {
	cval C.GDALDatasetH
}

type RasterBand struct {
	cval C.GDALRasterBandH
}

type Driver struct {
	cval C.GDALDriverH
}

type RasterAttributeTable struct {
	cval C.GDALRasterAttributeTableH
}

type AsyncReader struct {
	cval C.GDALAsyncReaderH
}

/* -------------------------------------------------------------------- */
/*      Callback "progress" function.                                   */
/* -------------------------------------------------------------------- */

type ProgressFunc func(complete float64, message string, progressArg interface{}) int

func DummyProgress(complete float64, message string, data interface{}) int {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))

	retval := C.GDALDummyProgress(C.double(complete), msg, unsafe.Pointer(nil))
	return int(retval)
}

func TermProgress(complete float64, message string, data interface{}) int {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))

	retval := C.GDALTermProgress(C.double(complete), msg, unsafe.Pointer(nil))
	return int(retval)
}

func ScaledProgress(complete float64, message string, data interface{}) int {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))

	retval := C.GDALScaledProgress(C.double(complete), msg, unsafe.Pointer(nil))
	return int(retval)
}

func CreateScaledProgress(min, max float64, progress ProgressFunc, data unsafe.Pointer) unsafe.Pointer {
	panic("not implemented!")
	return nil
}

func DestroyScaledProgress(data unsafe.Pointer) {
	C.GDALDestroyScaledProgress(data)
}

type goGDALProgressFuncProxyArgs struct {
	progresssFunc ProgressFunc
	data          interface{}
}

//export goGDALProgressFuncProxyA
func goGDALProgressFuncProxyA(complete C.double, message *C.char, data unsafe.Pointer) int {
	arg := (*goGDALProgressFuncProxyArgs)(data)
	return arg.progresssFunc(
		float64(complete), C.GoString(message), arg.data,
	)
}

/* ==================================================================== */
/*      Registration/driver related.                                    */
/* ==================================================================== */

const (
	DMD_LONGNAME           = string(C.GDAL_DMD_LONGNAME)
	DMD_HELPTOPIC          = string(C.GDAL_DMD_HELPTOPIC)
	DMD_MIMETYPE           = string(C.GDAL_DMD_MIMETYPE)
	DMD_EXTENSION          = string(C.GDAL_DMD_EXTENSION)
	DMD_CREATIONOPTIONLIST = string(C.GDAL_DMD_CREATIONOPTIONLIST)
	DMD_CREATIONDATATYPES  = string(C.GDAL_DMD_CREATIONDATATYPES)
	DCAP_CREATE            = string(C.GDAL_DCAP_CREATE)
	DCAP_CREATECOPY        = string(C.GDAL_DCAP_CREATECOPY)
	DCAP_VIRTUALIO         = string(C.GDAL_DCAP_VIRTUALIO)
)

// Open an existing dataset
func Open(filename string, access Access) (*Dataset, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	dataset := C.GDALOpen(cFilename, C.GDALAccess(access))
	if dataset == nil {
		return nil, fmt.Errorf("Error: dataset '%s' open error", filename)
	}
	return &Dataset{dataset}, nil
}

// Open a shared existing dataset
func OpenShared(filename string, access Access) *Dataset {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	dataset := C.GDALOpenShared(cFilename, C.GDALAccess(access))
	if dataset == nil {
		return nil
	}
	return &Dataset{dataset}
}

// TODO(kyle): deprecate Open(), rename OpenEx->Open
func OpenEx(filename string, flags Access, allowedDrivers []string, options []string, siblingFiles []string) (*Dataset, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	n := len(allowedDrivers)
	cDrivers := make([]*C.char, n+1)
	for i := 0; i < n; i++ {
		cDrivers[i] = C.CString(allowedDrivers[i])
		defer C.free(unsafe.Pointer(cDrivers[i]))
	}
	cDrivers[n] = (*C.char)(unsafe.Pointer(nil))

	n = len(options)
	cOptions := make([]*C.char, n+1)
	for i := 0; i < n; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[n] = (*C.char)(unsafe.Pointer(nil))

	n = len(siblingFiles)
	cSiblings := make([]*C.char, n+1)
	for i := 0; i < n; i++ {
		cSiblings[i] = C.CString(siblingFiles[i])
		defer C.free(unsafe.Pointer(cSiblings[i]))
	}
	cSiblings[n] = (*C.char)(unsafe.Pointer(nil))

	var dataset C.GDALDatasetH
	// The allowed drivers argument has to be handled specially.
	// nil -> all drivers
	// null terminated list (as the default code would produce) means no drivers
	if allowedDrivers == nil {
		dataset = C.GDALOpenEx(
			cFilename,
			C.uint(flags),
			nil,
			(**C.char)(unsafe.Pointer(&cOptions[0])),
			(**C.char)(unsafe.Pointer(&cSiblings[0])))
	} else {
		dataset = C.GDALOpenEx(
			cFilename,
			C.uint(flags),
			(**C.char)(unsafe.Pointer(&cDrivers[0])),
			(**C.char)(unsafe.Pointer(&cOptions[0])),
			(**C.char)(unsafe.Pointer(&cSiblings[0])))
	}
	if dataset == nil {
		return nil, fmt.Errorf("Error: dataset '%s' open error", filename)
	}
	return &Dataset{dataset}, nil
}

// Unimplemented: DumpOpenDatasets
/* ==================================================================== */
/*      GDAL_GCP                                                        */
/* ==================================================================== */

// Unimplemented: InitGCPs
// Unimplemented: DeinitGCPs
// Unimplemented: DuplicateGCPs
// Unimplemented: GCPsToGeoTransform
// Unimplemented: ApplyGeoTransform

/* ==================================================================== */
/*      major objects (dataset, and, driver, drivermanager).            */
/* ==================================================================== */

// Fetch object description
func (object MajorObject) Description() string {
	cObject := object.cval
	desc := C.GoString(C.GDALGetDescription(cObject))
	return desc
}

// Set object description
func (object MajorObject) SetDescription(desc string) {
	cObject := object.cval
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	C.GDALSetDescription(cObject, cDesc)
}

// Fetch metadata
func (object MajorObject) Metadata(domain string) []string {
	panic("not implemented!")
	return nil
}

// Set metadata
func (object MajorObject) SetMetadata(metadata []string, domain string) {
	panic("not implemented!")
	return
}

// Fetch a single metadata item
func (object MajorObject) MetadataItem(name, domain string) string {
	panic("not implemented!")
	return ""
}

// Set a single metadata item
func (object MajorObject) SetMetadataItem(name, value, domain string) {
	panic("not implemented!")
	return
}

func (dataset *Dataset) Metadata(domain string) []string {
	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	p := C.GDALGetMetadata(C.GDALMajorObjectH(dataset.cval), c_domain)
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

func (band *RasterBand) Metadata(domain string) []string {
	cDomain := C.CString(domain)
	defer C.free(unsafe.Pointer(cDomain))

	p := C.GDALGetMetadata(C.GDALMajorObjectH(band.cval), cDomain)
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

// TODO: Make korrekt class hirerarchy via interfaces

func (object *RasterBand) SetMetadataItem(name, value, domain string) error {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))

	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	return C.GDALSetMetadataItem(
		C.GDALMajorObjectH(unsafe.Pointer(object.cval)),
		c_name, c_value, c_domain,
	).Err()
}

func (band *RasterBand) MetadataItem(name, domain string) string {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cDomain := C.CString(domain)
	defer C.free(unsafe.Pointer(cDomain))

	md := C.GDALGetMetadataItem(C.GDALMajorObjectH(band.cval), cName, cDomain)
	if md != nil {
		return ""
	}
	return C.GoString(md)
}

// TODO: Make korrekt class hirerarchy via interfaces

func (object *Dataset) SetMetadataItem(name, value, domain string) error {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))

	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	return C.GDALSetMetadataItem(
		C.GDALMajorObjectH(unsafe.Pointer(object.cval)),
		c_name, c_value, c_domain,
	).Err()
}

// Fetch single metadata item.
func (object *Driver) MetadataItem(name, domain string) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	return C.GoString(
		C.GDALGetMetadataItem(
			C.GDALMajorObjectH(unsafe.Pointer(object.cval)),
			c_name, c_domain,
		),
	)
}

/* ==================================================================== */
/*      GDALDataset class ... normally this represents one file.        */
/* ==================================================================== */

// Get the driver to which this dataset relates
func (dataset *Dataset) Driver() *Driver {
	return &Driver{C.GDALGetDatasetDriver(dataset.cval)}
}

// Fetch files forming the dataset.
func (dataset *Dataset) FileList() []string {
	p := C.GDALGetFileList(dataset.cval)
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

// Close the dataset
func (dataset *Dataset) Close() {
	C.GDALClose(dataset.cval)
	return
}

// Fetch X size of raster
func (dataset *Dataset) RasterXSize() int {
	return int(C.GDALGetRasterXSize(dataset.cval))
}

// Fetch Y size of raster
func (dataset *Dataset) RasterYSize() int {
	return int(C.GDALGetRasterYSize(dataset.cval))
}

// Fetch the number of raster bands in the dataset
func (dataset *Dataset) RasterCount() int {
	return int(C.GDALGetRasterCount(dataset.cval))
}

// ErrInvalidBand represents an invalid band number when requested.  It is set
// when the return from GDALGetRasterBand(n) is NULL.
var ErrIllegalBand = errors.New("illegal band #")

// RasterBand returns the RasterBand at index band, where:
//
// 0 > band >= Dataset.RasterCount()
//
// If the band is invalid, nil and ErrInvalidBand is returned
func (dataset *Dataset) RasterBand(band int) (*RasterBand, error) {
	p := C.GDALGetRasterBand(dataset.cval, C.int(band))
	if p == nil {
		return nil, ErrIllegalBand
	}
	return &RasterBand{p}, nil
}

// Add a band to a dataset
func (dataset *Dataset) AddBand(dataType DataType, options []string) error {
	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALAddBand(
		dataset.cval,
		C.GDALDataType(dataType),
		(**C.char)(unsafe.Pointer(&cOptions[0])),
	).Err()
}

type ResampleAlg int

const (
	GRA_NearestNeighbour = ResampleAlg(0)
	GRA_Bilinear         = ResampleAlg(1)
	GRA_Cubic            = ResampleAlg(2)
	GRA_CubicSpline      = ResampleAlg(3)
	GRA_Lanczos          = ResampleAlg(4)
)

func (dataset *Dataset) AutoCreateWarpedVRT(srcWKT, dstWKT string, resampleAlg ResampleAlg) (*Dataset, error) {
	c_srcWKT := C.CString(srcWKT)
	defer C.free(unsafe.Pointer(c_srcWKT))
	c_dstWKT := C.CString(dstWKT)
	defer C.free(unsafe.Pointer(c_dstWKT))
	h := C.GDALAutoCreateWarpedVRT(dataset.cval, c_srcWKT, c_dstWKT, C.GDALResampleAlg(resampleAlg), 0.0, nil)
	if h == nil {
		return nil, fmt.Errorf("AutoCreateWarpedVRT failed")
	}
	return &Dataset{h}, nil
}

// Unimplemented: GDALBeginAsyncReader
// Unimplemented: GDALEndAsyncReader

// Read / write a region of image data from multiple bands
func (dataset *Dataset) IO(
	rwFlag RWFlag,
	xOff, yOff, xSize, ySize int,
	buffer interface{},
	bufXSize, bufYSize int,
	bandCount int,
	bandMap []int,
	pixelSpace, lineSpace, bandSpace int,
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

	return C.GDALDatasetRasterIO(
		dataset.cval,
		C.GDALRWFlag(rwFlag),
		C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize),
		dataPtr,
		C.int(bufXSize), C.int(bufYSize),
		C.GDALDataType(dataType),
		C.int(bandCount),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(bandMap)[0])),
		C.int(pixelSpace), C.int(lineSpace), C.int(bandSpace),
	).Err()
}

// Advise driver of upcoming read requests
func (dataset *Dataset) AdviseRead(
	rwFlag RWFlag,
	xOff, yOff, xSize, ySize, bufXSize, bufYSize int,
	dataType DataType,
	bandCount int,
	bandMap []int,
	options []string,
) error {
	length := len(options)
	cOptions := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cOptions[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(cOptions[i]))
	}
	cOptions[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALDatasetAdviseRead(
		dataset.cval,
		C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize),
		C.int(bufXSize), C.int(bufYSize),
		C.GDALDataType(dataType),
		C.int(bandCount),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(bandMap)[0])),
		(**C.char)(unsafe.Pointer(&cOptions[0])),
	).Err()
}

// Fetch the projection definition string for this dataset
func (dataset *Dataset) ProjectionRef() string {
	return C.GoString(C.GDALGetProjectionRef(dataset.cval))
}

// Set the projection reference string
func (dataset *Dataset) SetProjection(proj string) error {
	cProj := C.CString(proj)
	defer C.free(unsafe.Pointer(cProj))

	return C.GDALSetProjection(dataset.cval, cProj).Err()
}

// Get the affine transformation coefficients
func (dataset *Dataset) GeoTransform() [6]float64 {
	var transform [6]float64
	C.GDALGetGeoTransform(dataset.cval, (*C.double)(unsafe.Pointer(&transform[0])))
	return transform
}

// Set the affine transformation coefficients
func (dataset *Dataset) SetGeoTransform(transform [6]float64) error {
	return C.GDALSetGeoTransform(
		dataset.cval,
		(*C.double)(unsafe.Pointer(&transform[0])),
	).Err()
}

// Return the inverted transform
func (dataset *Dataset) InvGeoTransform() [6]float64 {
	return InvGeoTransform(dataset.GeoTransform())
}

// Invert the supplied transform
func InvGeoTransform(transform [6]float64) [6]float64 {
	var result [6]float64
	C.GDALInvGeoTransform((*C.double)(unsafe.Pointer(&transform[0])), (*C.double)(unsafe.Pointer(&result[0])))
	return result
}

// Get number of GCPs
func (dataset *Dataset) GDALGetGCPCount() int {
	return int(C.GDALGetGCPCount(dataset.cval))
}

// Unimplemented: GDALGetGCPProjection
// Unimplemented: GDALGetGCPs
// Unimplemented: GDALSetGCPs

// Fetch a format specific internally meaningful handle
func (dataset *Dataset) GDALGetInternalHandle(request string) unsafe.Pointer {
	cRequest := C.CString(request)
	defer C.free(unsafe.Pointer(cRequest))

	ptr := C.GDALGetInternalHandle(dataset.cval, cRequest)
	return ptr
}

// Add one to dataset reference count
func (dataset *Dataset) GDALReferenceDataset() int {
	count := C.GDALReferenceDataset(dataset.cval)
	return int(count)
}

// Subtract one from dataset reference count
func (dataset *Dataset) GDALDereferenceDataset() int {
	count := C.GDALDereferenceDataset(dataset.cval)
	return int(count)
}

// Build raster overview(s)
func (dataset *Dataset) BuildOverviews(
	resampling string,
	nOverviews int,
	overviewList []int,
	nBands int,
	bandList []int,
	progress ProgressFunc,
	data interface{},
) error {
	cResampling := C.CString(resampling)
	defer C.free(unsafe.Pointer(cResampling))

	arg := &goGDALProgressFuncProxyArgs{progress, data}

	return C.GDALBuildOverviews(
		dataset.cval,
		cResampling,
		C.int(nOverviews),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(overviewList)[0])),
		C.int(nBands),
		(*C.int)(unsafe.Pointer(&IntSliceToCInt(bandList)[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Unimplemented: GDALGetOpenDatasets

// Return access flag
func (dataset *Dataset) Access() Access {
	accessVal := C.GDALGetAccess(dataset.cval)
	return Access(accessVal)
}

// Write all write cached data to disk
func (dataset *Dataset) FlushCache() {
	C.GDALFlushCache(dataset.cval)
	return
}

// Adds a mask band to the dataset
func (dataset *Dataset) CreateMaskBand(flags int) error {
	return C.GDALCreateDatasetMaskBand(dataset.cval, C.int(flags)).Err()
}

// Copy all dataset raster data
func (sourceDataset *Dataset) CopyWholeRaster(
	destDataset *Dataset,
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

	return C.GDALDatasetCopyWholeRaster(
		sourceDataset.cval,
		destDataset.cval,
		(**C.char)(unsafe.Pointer(&cOptions[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// LayerCount gets the number of layers in this dataset.
func (ds *Dataset) LayerCount() int {
	return int(C.GDALDatasetGetLayerCount(ds.cval))
}

// LayerByIndex fetches a layer by index.
//
// The returned layer remains owned by the GDALDataset and should not be deleted
// by the application.
//
// This function is the same as the C++ method GDALDataset::GetLayer()
func (ds *Dataset) Layer(layer int) (*Layer, error) {
	lyr := C.GDALDatasetGetLayer(ds.cval, C.int(layer))
	if lyr == nil {
		return nil, fmt.Errorf("failed to get layer")
	}
	return &Layer{lyr}, nil
}

// LayerByName fetches a layer by name.
//
// The returned layer remains owned by the GDALDataset and should not be deleted
// by the application.
//
// This function is the same as the C++ method GDALDataset::GetLayerByName()
func (ds *Dataset) LayerByName(name string) (*Layer, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	lyr := C.GDALDatasetGetLayerByName(ds.cval, cName)
	if lyr == nil {
		return nil, fmt.Errorf("failed to get layer")
	}
	return &Layer{lyr}, nil
}

// ExecuteSQL Executes an SQL statement against the data store.
//
// The result of an SQL query is either NULL for statements that are in error,
// or that have no results set, or an OGRLayer pointer representing a results
// set from the query. Note that this OGRLayer is in addition to the layers in
// the data store and must be destroyed with ReleaseResultSet() before the
// dataset is closed (destroyed).
//
// This method is the same as the C++ method GDALDataset::ExecuteSQL()
//
// For more information on the SQL dialect supported internally by OGR review
// the OGR SQL document. Some drivers (i.e. Oracle and PostGIS) pass the SQL
// directly through to the underlying RDBMS.
func (ds *Dataset) ExecuteSQL(sql string, spatialFilter Geometry, dialect string) (*Layer, error) {
	cSQL := C.CString(sql)
	defer C.free(unsafe.Pointer(cSQL))
	var cDialect *C.char
	if dialect == "" {
		cDialect = nil
	} else {
		cDialect = C.CString(dialect)
		defer C.free(unsafe.Pointer(cDialect))
	}
	lyr := C.GDALDatasetExecuteSQL(ds.cval, cSQL, spatialFilter.cval, cDialect)
	if lyr == nil {
		return nil, fmt.Errorf("failed to execute SQL")
	}
	return &Layer{lyr}, nil
}

// Generate downsampled overviews
// Unimplemented: RegenerateOverviews

/* ==================================================================== */
/*     GDALAsyncReader                                                  */
/* ==================================================================== */

// Unimplemented: GetNextUpdatedRegion
// Unimplemented: LockBuffer
// Unimplemented: UnlockBuffer
