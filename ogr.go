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
	"unsafe"
)

func init() {
	C.OGRRegisterAll()
}

// OGR errors
var (
	ErrOGRNotEnoughData           = errors.New("not enough data")
	ErrOGRNotEnoughMemory         = errors.New("not enough memory")
	ErrOGRUnsupportedGeometryType = errors.New("unsupported geometry type")
	ErrOGRUnsupportedOperation    = errors.New("unsupported operation")
	ErrOGRCorruptData             = errors.New("corrupt data")
	ErrOGRFailure                 = errors.New("failure")
	ErrOGRUnsupportedSRS          = errors.New("unsupported SRS")
	ErrOGRInvalidHandle           = errors.New("invalid handle")
	ErrOGRNonExistingFeature      = errors.New("non-existing feature")
	EOL                           = errors.New("end of features")
)

type Envelope struct {
	cval C.OGREnvelope
}

func (env Envelope) MinX() float64 {
	return float64(env.cval.MinX)
}

func (env Envelope) MaxX() float64 {
	return float64(env.cval.MaxX)
}

func (env Envelope) MinY() float64 {
	return float64(env.cval.MinY)
}

func (env Envelope) MaxY() float64 {
	return float64(env.cval.MaxY)
}

func (env *Envelope) SetMinX(val float64) {
	env.cval.MinX = C.double(val)
}

func (env *Envelope) SetMaxX(val float64) {
	env.cval.MaxX = C.double(val)
}

func (env *Envelope) SetMinY(val float64) {
	env.cval.MinY = C.double(val)
}

func (env *Envelope) SetMaxY(val float64) {
	env.cval.MaxY = C.double(val)
}

func (env Envelope) IsInit() bool {
	return env.cval.MinX != 0 || env.cval.MinY != 0 || env.cval.MaxX != 0 || env.cval.MaxY != 0
}

func min(a, b C.double) C.double {
	if a < b {
		return a
	}
	return b
}

func max(a, b C.double) C.double {
	if a > b {
		return a
	}
	return b
}

// Return the union of this envelope with another one
func (env Envelope) Union(other Envelope) Envelope {
	if env.IsInit() {
		env.cval.MinX = min(env.cval.MinX, other.cval.MinX)
		env.cval.MinY = min(env.cval.MinY, other.cval.MinY)
		env.cval.MaxX = max(env.cval.MaxX, other.cval.MaxX)
		env.cval.MaxY = max(env.cval.MaxY, other.cval.MaxY)
	} else {
		env.cval.MinX = other.cval.MinX
		env.cval.MinY = other.cval.MinY
		env.cval.MaxX = other.cval.MaxX
		env.cval.MaxY = other.cval.MaxY
	}
	return env
}

// Return the intersection of this envelope with another
func (env Envelope) Intersect(other Envelope) {
	if env.Intersects(other) {
		if env.IsInit() {
			env.cval.MinX = max(env.cval.MinX, other.cval.MinX)
			env.cval.MinY = max(env.cval.MinY, other.cval.MinY)
			env.cval.MaxX = min(env.cval.MaxX, other.cval.MaxX)
			env.cval.MaxY = min(env.cval.MaxY, other.cval.MaxY)
		} else {
			env.cval.MinX = other.cval.MinX
			env.cval.MinY = other.cval.MinY
			env.cval.MaxX = other.cval.MaxX
			env.cval.MaxY = other.cval.MaxY
		}
	} else {
		env.cval.MinX = 0
		env.cval.MinY = 0
		env.cval.MaxX = 0
		env.cval.MaxY = 0
	}
}

// Test if one envelope intersects another
func (env Envelope) Intersects(other Envelope) bool {
	return env.cval.MinX <= other.cval.MaxX &&
		env.cval.MaxX >= other.cval.MinX &&
		env.cval.MinY <= other.cval.MaxY &&
		env.cval.MaxY >= other.cval.MinY
}

// Test if one envelope completely contains another
func (env Envelope) Contains(other Envelope) bool {
	return env.cval.MinX <= other.cval.MinX &&
		env.cval.MaxX >= other.cval.MaxX &&
		env.cval.MinY <= other.cval.MinY &&
		env.cval.MaxY >= other.cval.MaxY
}

/* -------------------------------------------------------------------- */
/*      Misc functions                                                  */
/* -------------------------------------------------------------------- */

// Clean up all OGR related resources
func CleanupOGR() {
	C.OGRCleanupAll()
}

// Convert a go bool to a C int
func BoolToCInt(in bool) (out C.int) {
	if in {
		out = 1
	} else {
		out = 0
	}
	return
}

/* -------------------------------------------------------------------- */
/*      Field definition functions                                      */
/* -------------------------------------------------------------------- */

// List of well known binary geometry types
type FieldType int

const (
	FT_Integer     = FieldType(C.OFTInteger)
	FT_IntegerList = FieldType(C.OFTIntegerList)
	FT_Real        = FieldType(C.OFTReal)
	FT_RealList    = FieldType(C.OFTRealList)
	FT_String      = FieldType(C.OFTString)
	FT_StringList  = FieldType(C.OFTStringList)
	FT_Binary      = FieldType(C.OFTBinary)
	FT_Date        = FieldType(C.OFTDate)
	FT_Time        = FieldType(C.OFTTime)
	FT_DateTime    = FieldType(C.OFTDateTime)
)

type Justification int

const (
	J_Undefined = Justification(C.OJUndefined)
	J_Left      = Justification(C.OJLeft)
	J_Right     = Justification(C.OJRight)
)

type FieldDefinition struct {
	cval C.OGRFieldDefnH
}

type Field struct {
	cval *C.OGRField
}

// Create a new field definition
func CreateFieldDefinition(name string, fieldType FieldType) *FieldDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fieldDef := C.OGR_Fld_Create(cName, C.OGRFieldType(fieldType))
	if fieldDef == nil {
		return nil
	}
	return &FieldDefinition{fieldDef}
}

// Destroy the field definition
func (fd *FieldDefinition) Destroy() {
	C.OGR_Fld_Destroy(fd.cval)
}

// Fetch the name of the field
func (fd *FieldDefinition) Name() string {
	name := C.OGR_Fld_GetNameRef(fd.cval)
	return C.GoString(name)
}

// Set the name of the field
func (fd *FieldDefinition) SetName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.OGR_Fld_SetName(fd.cval, cName)
}

// Fetch the type of this field
func (fd *FieldDefinition) Type() FieldType {
	fType := C.OGR_Fld_GetType(fd.cval)
	return FieldType(fType)
}

// Set the type of this field
func (fd *FieldDefinition) SetType(fType FieldType) {
	C.OGR_Fld_SetType(fd.cval, C.OGRFieldType(fType))
}

// Fetch the justification for this field
func (fd *FieldDefinition) Justification() Justification {
	justify := C.OGR_Fld_GetJustify(fd.cval)
	return Justification(justify)
}

// Set the justification for this field
func (fd *FieldDefinition) SetJustification(justify Justification) {
	C.OGR_Fld_SetJustify(fd.cval, C.OGRJustification(justify))
}

// Fetch the formatting width for this field
func (fd *FieldDefinition) Width() int {
	width := C.OGR_Fld_GetWidth(fd.cval)
	return int(width)
}

// Set the formatting width for this field
func (fd *FieldDefinition) SetWidth(width int) {
	C.OGR_Fld_SetWidth(fd.cval, C.int(width))
}

// Fetch the precision for this field
func (fd *FieldDefinition) Precision() int {
	precision := C.OGR_Fld_GetPrecision(fd.cval)
	return int(precision)
}

// Set the precision for this field
func (fd *FieldDefinition) SetPrecision(precision int) {
	C.OGR_Fld_SetPrecision(fd.cval, C.int(precision))
}

// Set defining parameters of field in a single call
func (fd *FieldDefinition) Set(
	name string,
	fType FieldType,
	width, precision int,
	justify Justification,
) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.OGR_Fld_Set(
		fd.cval,
		cName,
		C.OGRFieldType(fType),
		C.int(width),
		C.int(precision),
		C.OGRJustification(justify),
	)
}

// Fetch whether this field should be ignored when fetching features
func (fd *FieldDefinition) IsIgnored() bool {
	ignore := C.OGR_Fld_IsIgnored(fd.cval)
	return ignore != 0
}

// Set whether this field should be ignored when fetching features
func (fd *FieldDefinition) SetIgnored(ignore bool) {
	C.OGR_Fld_SetIgnored(fd.cval, BoolToCInt(ignore))
}

// Fetch human readable name for the field type
func (ft FieldType) Name() string {
	name := C.OGR_GetFieldTypeName(C.OGRFieldType(ft))
	return C.GoString(name)
}

/* -------------------------------------------------------------------- */
/*      Feature definition functions                                    */
/* -------------------------------------------------------------------- */

type FeatureDefinition struct {
	cval C.OGRFeatureDefnH
}

// Create a new feature definition object
func CreateFeatureDefinition(name string) *FeatureDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fd := C.OGR_FD_Create(cName)
	if fd == nil {
		return nil
	}
	return &FeatureDefinition{fd}
}

// Destroy a feature definition object
func (fd *FeatureDefinition) Destroy() {
	C.OGR_FD_Destroy(fd.cval)
}

// Drop a reference, and delete object if no references remain
func (fd *FeatureDefinition) Release() {
	C.OGR_FD_Release(fd.cval)
}

// Fetch the name of this feature definition
func (fd *FeatureDefinition) Name() string {
	name := C.OGR_FD_GetName(fd.cval)
	return C.GoString(name)
}

// Fetch the number of fields in the feature definition
func (fd *FeatureDefinition) FieldCount() int {
	count := C.OGR_FD_GetFieldCount(fd.cval)
	return int(count)
}

// Fetch the definition of the indicated field
func (fd *FeatureDefinition) FieldDefinition(index int) *FieldDefinition {
	fieldDefn := C.OGR_FD_GetFieldDefn(fd.cval, C.int(index))
	if fieldDefn == nil {
		return nil
	}
	return &FieldDefinition{fieldDefn}
}

// Fetch the index of the named field
func (fd *FeatureDefinition) FieldIndex(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	index := C.OGR_FD_GetFieldIndex(fd.cval, cName)
	return int(index)
}

// Add a new field definition to this feature definition
func (fd *FeatureDefinition) AddFieldDefinition(fieldDefn *FieldDefinition) {
	C.OGR_FD_AddFieldDefn(fd.cval, fieldDefn.cval)
}

// Delete a field definition from this feature definition
func (fd *FeatureDefinition) DeleteFieldDefinition(index int) error {
	return C.OGR_FD_DeleteFieldDefn(fd.cval, C.int(index)).Err()
}

// Fetch the geometry base type of this feature definition
func (fd *FeatureDefinition) GeometryType() GeometryType {
	gt := C.OGR_FD_GetGeomType(fd.cval)
	return GeometryType(gt)
}

// Set the geometry base type for this feature definition
func (fd *FeatureDefinition) SetGeometryType(geomType GeometryType) {
	C.OGR_FD_SetGeomType(fd.cval, C.OGRwkbGeometryType(geomType))
}

// Fetch if the geometry can be ignored when fetching features
func (fd *FeatureDefinition) IsGeometryIgnored() bool {
	isIgnored := C.OGR_FD_IsGeometryIgnored(fd.cval)
	return isIgnored != 0
}

// Set whether the geometry can be ignored when fetching features
func (fd *FeatureDefinition) SetGeometryIgnored(val bool) {
	C.OGR_FD_SetGeometryIgnored(fd.cval, BoolToCInt(val))
}

// Fetch if the style can be ignored when fetching features
func (fd *FeatureDefinition) IsStyleIgnored() bool {
	isIgnored := C.OGR_FD_IsStyleIgnored(fd.cval)
	return isIgnored != 0
}

// Set whether the style can be ignored when fetching features
func (fd *FeatureDefinition) SetStyleIgnored(val bool) {
	C.OGR_FD_SetStyleIgnored(fd.cval, BoolToCInt(val))
}

// Increment the reference count by one
func (fd *FeatureDefinition) Reference() int {
	count := C.OGR_FD_Reference(fd.cval)
	return int(count)
}

// Decrement the reference count by one
func (fd *FeatureDefinition) Dereference() int {
	count := C.OGR_FD_Dereference(fd.cval)
	return int(count)
}

// Fetch the current reference count
func (fd *FeatureDefinition) ReferenceCount() int {
	count := C.OGR_FD_GetReferenceCount(fd.cval)
	return int(count)
}

/* -------------------------------------------------------------------- */
/*      Data source functions                                           */
/* -------------------------------------------------------------------- */

type DataSource struct {
	cval C.OGRDataSourceH
}

// Open a file / data source with one of the registered drivers
func OpenDataSource(name string, update int) *DataSource {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ds := C.OGROpen(cName, C.int(update), nil)
	if ds == nil {
		return nil
	}
	return &DataSource{ds}
}

// Open a shared file / data source with one of the registered drivers
func OpenSharedDataSource(name string, update int) *DataSource {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ds := C.OGROpenShared(cName, C.int(update), nil)
	if ds == nil {
		return nil
	}
	return &DataSource{ds}
}

// Drop a reference to this datasource and destroy if reference is zero
func (ds *DataSource) Release() error {
	return C.OGRReleaseDataSource(ds.cval).Err()
}

// Return the number of opened data sources
func OpenDataSourceCount() int {
	count := C.OGRGetOpenDSCount()
	return int(count)
}

// Return the i'th datasource opened
func OpenDataSourceByIndex(index int) *DataSource {
	ds := C.OGRGetOpenDS(C.int(index))
	if ds == nil {
		return nil
	}
	return &DataSource{ds}
}

// Closes datasource and releases resources
func (ds *DataSource) Destroy() {
	C.OGR_DS_Destroy(ds.cval)
}

// Fetch the name of the data source
func (ds *DataSource) Name() string {
	name := C.OGR_DS_GetName(ds.cval)
	return C.GoString(name)
}

// Fetch the number of layers in this data source
func (ds *DataSource) LayerCount() int {
	count := C.OGR_DS_GetLayerCount(ds.cval)
	return int(count)
}

// Fetch a layer of this data source by index
func (ds *DataSource) LayerByIndex(index int) *Layer {
	layer := C.OGR_DS_GetLayer(ds.cval, C.int(index))
	if layer == nil {
		return nil
	}
	return &Layer{layer}
}

// Fetch a layer of this data source by name
func (ds *DataSource) LayerByName(name string) *Layer {
	cString := C.CString(name)
	defer C.free(unsafe.Pointer(cString))
	layer := C.OGR_DS_GetLayerByName(ds.cval, cString)
	if layer == nil {
		return nil
	}
	return &Layer{layer}
}

// Delete the layer from the data source
func (ds *DataSource) Delete(index int) error {
	return C.OGR_DS_DeleteLayer(ds.cval, C.int(index)).Err()
}

// Fetch the driver that the data source was opened with
func (ds *DataSource) Driver() *OGRDriver {
	driver := C.OGR_DS_GetDriver(ds.cval)
	if driver == nil {
		return nil
	}
	return &OGRDriver{driver}
}

// Create a new layer on the data source
func (ds *DataSource) CreateLayer(
	name string,
	sr SpatialReference,
	geomType GeometryType,
	options []string,
) *Layer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	layer := C.OGR_DS_CreateLayer(
		ds.cval,
		cName,
		sr.cval,
		C.OGRwkbGeometryType(geomType),
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	if layer == nil {
		return nil
	}
	return &Layer{layer}
}

// Duplicate an existing layer
func (ds *DataSource) CopyLayer(
	source Layer,
	name string,
	options []string,
) *Layer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	layer := C.OGR_DS_CopyLayer(
		ds.cval,
		source.cval,
		cName,
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	if layer == nil {
		return nil
	}
	return &Layer{layer}
}

// Test if the data source has the indicated capability
func (ds *DataSource) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_DS_TestCapability(ds.cval, cString)
	return val != 0
}

// Execute an SQL statement against the data source
func (ds *DataSource) ExecuteSQL(sql string, filter Geometry, dialect string) *Layer {
	cSQL := C.CString(sql)
	defer C.free(unsafe.Pointer(cSQL))
	cDialect := C.CString(dialect)
	defer C.free(unsafe.Pointer(cDialect))

	layer := C.OGR_DS_ExecuteSQL(ds.cval, cSQL, filter.cval, cDialect)
	if layer == nil {
		return nil
	}
	return &Layer{layer}
}

// Release the results of ExecuteSQL
func (ds *DataSource) ReleaseResultSet(layer *Layer) {
	C.OGR_DS_ReleaseResultSet(ds.cval, layer.cval)
}

// Flush pending changes to the data source
func (ds *DataSource) Sync() error {
	return C.OGR_DS_SyncToDisk(ds.cval).Err()
}

/* -------------------------------------------------------------------- */
/*      Driver functions                                                */
/* -------------------------------------------------------------------- */

type OGRDriver struct {
	cval C.OGRSFDriverH
}

// Fetch name of driver (file format)
func (driver *OGRDriver) Name() string {
	name := C.OGR_Dr_GetName(driver.cval)
	return C.GoString(name)
}

// Attempt to open file with this driver
func (driver *OGRDriver) Open(filename string, update int) (newDS *DataSource, ok bool) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	ds := C.OGR_Dr_Open(driver.cval, cFilename, C.int(update))
	if ds == nil {
		return nil, false
	}
	return &DataSource{ds}, true
}

// Test if this driver supports the named capability
func (driver *OGRDriver) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_Dr_TestCapability(driver.cval, cString)
	return val != 0
}

// Create a new data source based on this driver
func (driver *OGRDriver) Create(name string, options []string) (newDS *DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))
	ds := C.OGR_Dr_CreateDataSource(driver.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	if ds == nil {
		return nil, false
	}
	return &DataSource{ds}, true
}

// Create a new datasource with this driver by copying all layers of the existing datasource
func (driver *OGRDriver) Copy(source DataSource, name string, options []string) (newDS *DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CopyDataSource(driver.cval, source.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	if ds == nil {
		return nil, false
	}
	return &DataSource{ds}, true
}

// Delete a data source
func (driver *OGRDriver) Delete(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return C.OGR_Dr_DeleteDataSource(driver.cval, cFilename).Err()
}

// Add a driver to the list of registered drivers
func (driver *OGRDriver) Register() {
	C.OGRRegisterDriver(driver.cval)
}

// Remove a driver from the list of registered drivers
func (driver *OGRDriver) Deregister() {
	C.OGRDeregisterDriver(driver.cval)
}

// Fetch the number of registered drivers
func OGRDriverCount() int {
	count := C.OGRGetDriverCount()
	return int(count)
}

// Fetch the indicated driver by index
func OGRDriverByIndex(index int) *OGRDriver {
	driver := C.OGRGetDriver(C.int(index))
	if driver == nil {
		return nil
	}
	return &OGRDriver{driver}
}

// Fetch the indicated driver by name
func OGRDriverByName(name string) *OGRDriver {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	driver := C.OGRGetDriverByName(cName)
	if driver == nil {
		return nil
	}
	return &OGRDriver{driver}
}
