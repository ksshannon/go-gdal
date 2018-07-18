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

type Layer struct {
	cval C.OGRLayerH
}

// Return the layer name
func (layer *Layer) Name() string {
	name := C.OGR_L_GetName(layer.cval)
	return C.GoString(name)
}

// Return the layer geometry type
func (layer *Layer) GeomType() GeometryType {
	gt := C.OGR_L_GetGeomType(layer.cval)
	return GeometryType(gt)
}

// Return the current spatial filter for this layer
func (layer *Layer) SpatialFilter() *Geometry {
	geom := C.OGR_L_GetSpatialFilter(layer.cval)
	if geom == nil {
		return nil
	}
	return &Geometry{geom}
}

// Set a new spatial filter for this layer
func (layer *Layer) SetSpatialFilter(filter *Geometry) {
	if filter == nil {
		C.OGR_L_SetSpatialFilter(layer.cval, nil)
	} else {
		C.OGR_L_SetSpatialFilter(layer.cval, filter.cval)
	}
}

// Set a new rectangular spatial filter for this layer
func (layer *Layer) SetSpatialFilterRect(minX, minY, maxX, maxY float64) {
	C.OGR_L_SetSpatialFilterRect(
		layer.cval,
		C.double(minX), C.double(minY), C.double(maxX), C.double(maxY),
	)
}

// Set a new attribute query filter
func (layer *Layer) SetAttributeFilter(filter string) error {
	cFilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cFilter))
	return C.OGR_L_SetAttributeFilter(layer.cval, cFilter).Err()
}

// Reset reading to start on the first feature
func (layer *Layer) ResetReading() {
	C.OGR_L_ResetReading(layer.cval)
}

// Fetch the next available feature from this layer
func (layer *Layer) NextFeature() *Feature {
	feature := C.OGR_L_GetNextFeature(layer.cval)
	if feature == nil {
		return nil
	}
	return &Feature{feature}
}

// Move read cursor to the provided index
func (layer *Layer) SetNextByIndex(index int) error {
	return C.OGR_L_SetNextByIndex(layer.cval, C.GIntBig(index)).Err()
}

// Fetch a feature by its index
func (layer *Layer) Feature(index int) *Feature {
	feature := C.OGR_L_GetFeature(layer.cval, C.GIntBig(index))
	if feature == nil {
		return nil
	}
	return &Feature{feature}
}

// Rewrite the provided feature
func (layer *Layer) SetFeature(feature *Feature) error {
	return C.OGR_L_SetFeature(layer.cval, feature.cval).Err()
}

// Create and write a new feature within a layer
func (layer *Layer) CreateFeature(feature *Feature) error {
	return C.OGR_L_CreateFeature(layer.cval, feature.cval).Err()
}

// Delete indicated feature from layer
func (layer *Layer) DeleteFeature(index int) error {
	return C.OGR_L_DeleteFeature(layer.cval, C.GIntBig(index)).Err()
}

// Fetch the schema information for this layer
func (layer *Layer) Definition() *FeatureDefinition {
	defn := C.OGR_L_GetLayerDefn(layer.cval)
	if defn == nil {
		return nil
	}
	return &FeatureDefinition{defn}
}

// Fetch the spatial reference system for this layer
func (layer *Layer) SpatialRef() *SpatialReference {
	sr := C.OGR_L_GetSpatialRef(layer.cval)
	if sr == nil {
		return nil
	}
	return &SpatialReference{sr}
}

// Fetch the feature count for this layer
func (layer *Layer) FeatureCount(force bool) (count int, ok bool) {
	count = int(C.OGR_L_GetFeatureCount(layer.cval, BoolToCInt(force)))
	return count, count != -1
}

// Fetch the extent of this layer
func (layer *Layer) Extent(force bool) (env Envelope, err error) {
	err = C.OGR_L_GetExtent(layer.cval, &env.cval, BoolToCInt(force)).Err()
	return
}

// Test if this layer supports the named capability
func (layer *Layer) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_L_TestCapability(layer.cval, cString)
	return val != 0
}

// Create a new field on a layer
func (layer *Layer) CreateField(fd *FieldDefinition, approxOK bool) error {
	return C.OGR_L_CreateField(layer.cval, fd.cval, BoolToCInt(approxOK)).Err()
}

// Delete a field from the layer
func (layer *Layer) DeleteField(index int) error {
	return C.OGR_L_DeleteField(layer.cval, C.int(index)).Err()
}

// Reorder all the fields of a layer
func (layer *Layer) ReorderFields(layerMap []int) error {
	return C.OGR_L_ReorderFields(layer.cval, (*C.int)(unsafe.Pointer(&layerMap[0]))).Err()
}

// Reorder an existing field of a layer
func (layer *Layer) ReorderField(oldIndex, newIndex int) error {
	return C.OGR_L_ReorderField(layer.cval, C.int(oldIndex), C.int(newIndex)).Err()
}

// Alter the definition of an existing field of a layer
func (layer *Layer) AlterFieldDefn(index int, newDefn FieldDefinition, flags int) error {
	return C.OGR_L_AlterFieldDefn(layer.cval, C.int(index), newDefn.cval, C.int(flags)).Err()
}

// Begin a transation on data sources which support it
func (layer *Layer) StartTransaction() error {
	return C.OGR_L_StartTransaction(layer.cval).Err()
}

// Commit a transaction on data sources which support it
func (layer *Layer) CommitTransaction() error {
	return C.OGR_L_CommitTransaction(layer.cval).Err()
}

// Roll back the current transaction on data sources which support it
func (layer *Layer) RollbackTransaction() error {
	return C.OGR_L_RollbackTransaction(layer.cval).Err()
}

// Flush pending changes to the layer
func (layer *Layer) Sync() error {
	return C.OGR_L_SyncToDisk(layer.cval).Err()
}

// Fetch the name of the FID column
func (layer *Layer) FIDColumn() string {
	name := C.OGR_L_GetFIDColumn(layer.cval)
	return C.GoString(name)
}

// Fetch the name of the geometry column
func (layer *Layer) GeometryColumn() string {
	name := C.OGR_L_GetGeometryColumn(layer.cval)
	return C.GoString(name)
}

// Set which fields can be ignored when retrieving features from the layer
func (layer *Layer) SetIgnoredFields(names []string) error {
	length := len(names)
	cNames := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cNames[i] = C.CString(names[i])
		defer C.free(unsafe.Pointer(cNames[i]))
	}
	cNames[length] = (*C.char)(unsafe.Pointer(nil))

	return C.OGR_L_SetIgnoredFields(layer.cval, (**C.char)(unsafe.Pointer(&cNames[0]))).Err()
}

// Return the intersection of two layers
// Unimplemented: Intersection
// Will be new in 2.0

// Return the union of two layers
// Unimplemented: Union
// Will be new in 2.0

// Return the symmetric difference of two layers
// Unimplemented: SymDifference
// Will be new in 2.0

// Identify features in this layer with ones from the provided layer
// Unimplemented: Identity
// Will be new in 2.0

// Update this layer with features from the provided layer
// Unimplemented: Update
// Will be new in 2.0

// Clip off areas that are not covered by the provided layer
// Unimplemented: Clip
// Will be new in 2.0

// Remove areas that are covered by the provided layer
// Unimplemented: Erase
// Will be new in 2.0
