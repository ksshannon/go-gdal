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
	"unsafe"
)

// List of well known binary geometry types
type GeometryType uint32

const (
	GT_Unknown               = GeometryType(C.wkbUnknown)
	GT_Point                 = GeometryType(C.wkbPoint)
	GT_LineString            = GeometryType(C.wkbLineString)
	GT_Polygon               = GeometryType(C.wkbPolygon)
	GT_MultiPoint            = GeometryType(C.wkbMultiPoint)
	GT_MultiLineString       = GeometryType(C.wkbMultiLineString)
	GT_MultiPolygon          = GeometryType(C.wkbMultiPolygon)
	GT_GeometryCollection    = GeometryType(C.wkbGeometryCollection)
	GT_None                  = GeometryType(C.wkbNone)
	GT_LinearRing            = GeometryType(C.wkbLinearRing)
	GT_Point25D              = GeometryType(C.wkbPoint25D)
	GT_LineString25D         = GeometryType(C.wkbLineString25D)
	GT_Polygon25D            = GeometryType(C.wkbPolygon25D)
	GT_MultiPoint25D         = GeometryType(C.wkbMultiPoint25D)
	GT_MultiLineString25D    = GeometryType(C.wkbMultiLineString25D)
	GT_MultiPolygon25D       = GeometryType(C.wkbMultiPolygon25D)
	GT_GeometryCollection25D = GeometryType(C.wkbGeometryCollection25D)
)

type Geometry struct {
	cval C.OGRGeometryH
}

//Create a geometry object from its well known binary representation
func CreateFromWKB(wkb []uint8, srs *SpatialReference, bytes int) (*Geometry, error) {
	buf := (*C.uchar)(unsafe.Pointer(&wkb[0]))
	var newGeom *Geometry
	return newGeom, C.OGR_G_CreateFromWkb(
		buf, srs.cval, &newGeom.cval, C.int(bytes),
	).Err()
}

//Create a geometry object from its well known text representation
func CreateFromWKT(wkt string, srs SpatialReference) (*Geometry, error) {
	cString := C.CString(wkt)
	defer C.free(unsafe.Pointer(cString))
	var g C.OGRGeometryH
	rc := C.OGR_G_CreateFromWkt(&cString, srs.cval, &g)
	if rc != C.OGRERR_NONE {
		return nil, rc.Err()
	}
	return &Geometry{g}, nil
}

//Create a geometry object from its GeoJSON representation
func CreateFromJson(_json string) *Geometry {
	cString := C.CString(_json)
	defer C.free(unsafe.Pointer(cString))
	var newGeom *Geometry
	newGeom.cval = C.OGR_G_CreateGeometryFromJson(cString)
	return newGeom
}

// Destroy geometry object
func (geometry *Geometry) Destroy() {
	C.OGR_G_DestroyGeometry(geometry.cval)
}

// Create an empty geometry of the desired type
func Create(geomType GeometryType) *Geometry {
	geom := C.OGR_G_CreateGeometry(C.OGRwkbGeometryType(geomType))
	if geom == nil {
		return nil
	}
	return &Geometry{geom}
}

// Stroke arc to linestring
func ApproximateArcAngles(
	x, y, z,
	primaryRadius,
	secondaryRadius,
	rotation,
	startAngle,
	endAngle,
	stepSizeDegrees float64,
) *Geometry {
	geom := C.OGR_G_ApproximateArcAngles(
		C.double(x),
		C.double(y),
		C.double(z),
		C.double(primaryRadius),
		C.double(secondaryRadius),
		C.double(rotation),
		C.double(startAngle),
		C.double(endAngle),
		C.double(stepSizeDegrees))
	return &Geometry{geom}
}

// Convert to polygon
func (geom *Geometry) ForceToPolygon() *Geometry {
	newGeom := C.OGR_G_ForceToPolygon(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Convert to multipolygon
func (geom *Geometry) ForceToMultiPolygon() *Geometry {
	newGeom := C.OGR_G_ForceToMultiPolygon(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Convert to multipoint
func (geom *Geometry) ForceToMultiPoint() *Geometry {
	newGeom := C.OGR_G_ForceToMultiPoint(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Convert to multilinestring
func (geom *Geometry) ForceToMultiLineString() *Geometry {
	newGeom := C.OGR_G_ForceToMultiLineString(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Get the dimension of this geometry
func (geom *Geometry) Dimension() int {
	dim := C.OGR_G_GetDimension(geom.cval)
	return int(dim)
}

// Get the dimension of the coordinates in this geometry
func (geom *Geometry) CoordinateDimension() int {
	dim := C.OGR_G_GetCoordinateDimension(geom.cval)
	return int(dim)
}

// Set the dimension of the coordinates in this geometry
func (geom *Geometry) SetCoordinateDimension(dim int) {
	C.OGR_G_SetCoordinateDimension(geom.cval, C.int(dim))
}

// Create a copy of this geometry
func (geom *Geometry) Clone() *Geometry {
	newGeom := C.OGR_G_Clone(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute and return the bounding envelope for this geometry
func (geom *Geometry) Envelope() Envelope {
	var env Envelope
	C.OGR_G_GetEnvelope(geom.cval, &env.cval)
	return env
}

// Unimplemented: GetEnvelope3D

// Assign a geometry from well known binary data
func (geom *Geometry) FromWKB(wkb []uint8, bytes int) error {
	buf := (*C.uchar)(unsafe.Pointer(&wkb[0]))
	return C.OGR_G_ImportFromWkb(geom.cval, buf, C.int(bytes)).Err()
}

// Convert a geometry to well known binary data
func (geom *Geometry) ToWKB() ([]uint8, error) {
	b := make([]uint8, geom.WKBSize())
	cString := (*C.uchar)(unsafe.Pointer(&b[0]))
	err := C.OGR_G_ExportToWkb(geom.cval, C.OGRwkbByteOrder(C.wkbNDR), cString).Err()
	return b, err
}

// Returns size of related binary representation
func (geom *Geometry) WKBSize() int {
	size := C.OGR_G_WkbSize(geom.cval)
	return int(size)
}

// Assign geometry object from its well known text representation
func (geom *Geometry) FromWKT(wkt string) error {
	cString := C.CString(wkt)
	defer C.free(unsafe.Pointer(cString))
	return C.OGR_G_ImportFromWkt(geom.cval, &cString).Err()
}

// Fetch geometry as WKT
func (geom *Geometry) ToWKT() (string, error) {
	var p *C.char
	err := C.OGR_G_ExportToWkt(geom.cval, &p).Err()
	wkt := C.GoString(p)
	return wkt, err
}

// Fetch geometry type
func (geom *Geometry) Type() GeometryType {
	gt := C.OGR_G_GetGeometryType(geom.cval)
	return GeometryType(gt)
}

// Fetch geometry name
func (geom *Geometry) Name() string {
	name := C.OGR_G_GetGeometryName(geom.cval)
	return C.GoString(name)
}

// Unimplemented: DumpReadable

// Convert geometry to strictly 2D
func (geom *Geometry) FlattenTo2D() {
	C.OGR_G_FlattenTo2D(geom.cval)
}

// Force rings to be closed
func (geom *Geometry) CloseRings() {
	C.OGR_G_CloseRings(geom.cval)
}

// Create a geometry from its GML representation
func CreateFromGML(gml string) *Geometry {
	cString := C.CString(gml)
	defer C.free(unsafe.Pointer(cString))
	geom := C.OGR_G_CreateFromGML(cString)
	if geom == nil {
		return nil
	}
	return &Geometry{geom}
}

// Convert a geometry to GML format
func (geom *Geometry) ToGML() string {
	val := C.OGR_G_ExportToGML(geom.cval)
	return C.GoString(val)
}

// Convert a geometry to GML format with options
func (geom *Geometry) ToGML_Ex(options []string) string {
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	val := C.OGR_G_ExportToGMLEx(geom.cval, (**C.char)(unsafe.Pointer(&opts[0])))
	return C.GoString(val)
}

// Convert a geometry to KML format
func (geom *Geometry) ToKML() string {
	val := C.OGR_G_ExportToKML(geom.cval, nil)
	return C.GoString(val)
}

// Convert a geometry to JSON format
func (geom *Geometry) ToJSON() string {
	val := C.OGR_G_ExportToJson(geom.cval)
	return C.GoString(val)
}

// Convert a geometry to JSON format with options
func (geom *Geometry) ToJSON_ex(options []string) string {
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	val := C.OGR_G_ExportToJsonEx(geom.cval, (**C.char)(unsafe.Pointer(&opts[0])))
	return C.GoString(val)
}

// Fetch the spatial reference associated with this geometry
func (geom *Geometry) SpatialReference() *SpatialReference {
	spatialRef := C.OGR_G_GetSpatialReference(geom.cval)
	if spatialRef == nil {
		return nil
	}
	return &SpatialReference{spatialRef}
}

// Assign a spatial reference to this geometry
func (geom *Geometry) SetSpatialReference(spatialRef *SpatialReference) {
	C.OGR_G_AssignSpatialReference(geom.cval, spatialRef.cval)
}

// Apply coordinate transformation to geometry
func (geom *Geometry) Transform(ct *CoordinateTransform) error {
	return C.OGR_G_Transform(geom.cval, ct.cval).Err()
}

// Transform geometry to new spatial reference system
func (geom *Geometry) TransformTo(sr *SpatialReference) error {
	return C.OGR_G_TransformTo(geom.cval, sr.cval).Err()
}

// Simplify the geometry
func (geom *Geometry) Simplify(tolerance float64) *Geometry {
	newGeom := C.OGR_G_Simplify(geom.cval, C.double(tolerance))
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Simplify the geometry while preserving topology
func (geom *Geometry) SimplifyPreservingTopology(tolerance float64) *Geometry {
	newGeom := C.OGR_G_SimplifyPreserveTopology(geom.cval, C.double(tolerance))
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Modify the geometry such that it has no line segment longer than the given distance
func (geom *Geometry) Segmentize(distance float64) {
	C.OGR_G_Segmentize(geom.cval, C.double(distance))
}

// Return true if these features intersect
func (geom *Geometry) Intersects(other *Geometry) bool {
	val := C.OGR_G_Intersects(geom.cval, other.cval)
	return val != 0
}

// Return true if these features are equal
func (geom *Geometry) Equals(other *Geometry) bool {
	val := C.OGR_G_Equals(geom.cval, other.cval)
	return val != 0
}

// Return true if the features are disjoint
func (geom *Geometry) Disjoint(other *Geometry) bool {
	val := C.OGR_G_Disjoint(geom.cval, other.cval)
	return val != 0
}

// Return true if this feature touches the other
func (geom *Geometry) Touches(other *Geometry) bool {
	val := C.OGR_G_Touches(geom.cval, other.cval)
	return val != 0
}

// Return true if this feature crosses the other
func (geom *Geometry) Crosses(other *Geometry) bool {
	val := C.OGR_G_Crosses(geom.cval, other.cval)
	return val != 0
}

// Return true if this geometry is within the other
func (geom *Geometry) Within(other *Geometry) bool {
	val := C.OGR_G_Within(geom.cval, other.cval)
	return val != 0
}

// Return true if this geometry contains the other
func (geom *Geometry) Contains(other *Geometry) bool {
	val := C.OGR_G_Contains(geom.cval, other.cval)
	return val != 0
}

// Return true if this geometry overlaps the other
func (geom *Geometry) Overlaps(other *Geometry) bool {
	val := C.OGR_G_Overlaps(geom.cval, other.cval)
	return val != 0
}

// Compute boundary for the geometry
func (geom *Geometry) Boundary() *Geometry {
	newGeom := C.OGR_G_Boundary(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute convex hull for the geometry
func (geom *Geometry) ConvexHull() *Geometry {
	newGeom := C.OGR_G_ConvexHull(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute buffer of the geometry
func (geom *Geometry) Buffer(distance float64, segments int) *Geometry {
	newGeom := C.OGR_G_Buffer(geom.cval, C.double(distance), C.int(segments))
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute intersection of this geometry with the other
func (geom *Geometry) Intersection(other *Geometry) *Geometry {
	newGeom := C.OGR_G_Intersection(geom.cval, other.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute union of this geometry with the other
func (geom *Geometry) Union(other *Geometry) *Geometry {
	newGeom := C.OGR_G_Union(geom.cval, other.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Unimplemented: UnionCascaded

// Unimplemented: PointOn Surface (until 2.0)
// Return a point guaranteed to lie on the surface
// func (geom Geometry) PointOnSurface() Geometry {
//	newGeom := C.OGR_G_PointOnSurface(geom.cval)
//	return Geometry{newGeom}
// }

// Compute difference between this geometry and the other
func (geom *Geometry) Difference(other *Geometry) *Geometry {
	newGeom := C.OGR_G_Difference(geom.cval, other.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute symmetric difference between this geometry and the other
func (geom *Geometry) SymmetricDifference(other *Geometry) *Geometry {
	newGeom := C.OGR_G_SymDifference(geom.cval, other.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Compute distance between thie geometry and the other
func (geom *Geometry) Distance(other *Geometry) float64 {
	dist := C.OGR_G_Distance(geom.cval, other.cval)
	return float64(dist)
}

// Compute length of geometry
func (geom *Geometry) Length() float64 {
	length := C.OGR_G_Length(geom.cval)
	return float64(length)
}

// Compute area of geometry
func (geom *Geometry) Area() float64 {
	area := C.OGR_G_Area(geom.cval)
	return float64(area)
}

// Compute centroid of geometry
func (geom *Geometry) Centroid() *Geometry {
	var centroid *Geometry
	C.OGR_G_Centroid(geom.cval, centroid.cval)
	return centroid
}

// Clear the geometry to its uninitialized state
func (geom *Geometry) Empty() {
	C.OGR_G_Empty(geom.cval)
}

// Test if the geometry is empty
func (geom *Geometry) IsEmpty() bool {
	val := C.OGR_G_IsEmpty(geom.cval)
	return val != 0
}

// Test if the geometry is valid
func (geom *Geometry) IsValid() bool {
	val := C.OGR_G_IsValid(geom.cval)
	return val != 0
}

// Test if the geometry is simple
func (geom *Geometry) IsSimple() bool {
	val := C.OGR_G_IsSimple(geom.cval)
	return val != 0
}

// Test if the geometry is a ring
func (geom *Geometry) IsRing() bool {
	val := C.OGR_G_IsRing(geom.cval)
	return val != 0
}

// Polygonize a set of sparse edges
func (geom *Geometry) Polygonize() *Geometry {
	newGeom := C.OGR_G_Polygonize(geom.cval)
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Fetch number of points in the geometry
func (geom *Geometry) PointCount() int {
	count := C.OGR_G_GetPointCount(geom.cval)
	return int(count)
}

// Fetch the X coordinate of a point in the geometry
func (geom *Geometry) X(index int) float64 {
	x := C.OGR_G_GetX(geom.cval, C.int(index))
	return float64(x)
}

// Fetch the Y coordinate of a point in the geometry
func (geom *Geometry) Y(index int) float64 {
	y := C.OGR_G_GetY(geom.cval, C.int(index))
	return float64(y)
}

// Fetch the Z coordinate of a point in the geometry
func (geom *Geometry) Z(index int) float64 {
	z := C.OGR_G_GetZ(geom.cval, C.int(index))
	return float64(z)
}

// Fetch the coordinates of a point in the geometry
func (geom *Geometry) Point(index int) (x, y, z float64) {
	C.OGR_G_GetPoint(
		geom.cval,
		C.int(index),
		(*C.double)(&x),
		(*C.double)(&y),
		(*C.double)(&z))
	return
}

// Set the coordinates of a point in the geometry
func (geom *Geometry) SetPoint(index int, x, y, z float64) {
	C.OGR_G_SetPoint(
		geom.cval,
		C.int(index),
		C.double(x),
		C.double(y),
		C.double(z))
}

// Set the coordinates of a point in the geometry, ignoring the 3rd dimension
func (geom *Geometry) SetPoint2D(index int, x, y float64) {
	C.OGR_G_SetPoint_2D(geom.cval, C.int(index), C.double(x), C.double(y))
}

// Add a new point to the geometry (line string or polygon only)
func (geom *Geometry) AddPoint(x, y, z float64) {
	C.OGR_G_AddPoint(geom.cval, C.double(x), C.double(y), C.double(z))
}

// Add a new point to the geometry (line string or polygon only), ignoring the 3rd dimension
func (geom *Geometry) AddPoint2D(x, y float64) {
	C.OGR_G_AddPoint_2D(geom.cval, C.double(x), C.double(y))
}

// Fetch the number of elements in the geometry, or number of geometries in the container
func (geom *Geometry) GeometryCount() int {
	count := C.OGR_G_GetGeometryCount(geom.cval)
	return int(count)
}

// Fetch geometry from a geometry container
func (geom *Geometry) Geometry(index int) *Geometry {
	newGeom := C.OGR_G_GetGeometryRef(geom.cval, C.int(index))
	if newGeom == nil {
		return nil
	}
	return &Geometry{newGeom}
}

// Add a geometry to a geometry container
func (geom *Geometry) AddGeometry(other *Geometry) error {
	return C.OGR_G_AddGeometry(geom.cval, other.cval).Err()
}

// Add a geometry to a geometry container and assign ownership to that container
func (geom *Geometry) AddGeometryDirectly(other *Geometry) error {
	return C.OGR_G_AddGeometryDirectly(geom.cval, other.cval).Err()
}

// Remove a geometry from the geometry container
func (geom Geometry) RemoveGeometry(index int, delete bool) error {
	return C.OGR_G_RemoveGeometry(geom.cval, C.int(index), BoolToCInt(delete)).Err()
}

// Build a polygon / ring from a set of lines
func (geom Geometry) BuildPolygonFromEdges(autoClose bool, tolerance float64) (Geometry, error) {
	var cErr C.OGRErr
	newGeom := C.OGRBuildPolygonFromEdges(
		geom.cval,
		0,
		BoolToCInt(autoClose),
		C.double(tolerance),
		&cErr,
	)
	return Geometry{newGeom}, cErr.Err()
}
