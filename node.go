package db

import (
	"fmt"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/random"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/pln"
	"github.com/TopoSimplify/rng"
	"github.com/TopoSimplify/seg"
	"github.com/intdxdt/iter"
)

//Node
type Node struct {
	Id          int
	Coordinates geom.Coords
	polyline    pln.Polyline
	Range       rng.Rng
	geom        geom.Geometry

	FID         int
	NID         int
	HullType    geom.GeoType
	WTK         string
}

func NewDBNode(id *iter.Igen, coordinates geom.Coords, r rng.Rng, fid int, gfn func(geom.Coords) geom.Geometry) Node {
	var n = NewDBNodeFromDPNode(node.CreateNode(id, coordinates, r, gfn))
	n.FID = fid
	return n
}

func NewDBNodeFromDPNode(node node.Node) Node {
	return Node{
		Id:          node.Id,
		Coordinates: node.Polyline.Coordinates,
		Range:       node.Range,
		WTK:         node.Geom.WKT(),
		HullType:    node.Geom.Type(),
		geom:        node.Geom,
	}
}

func (n *Node) UpdateSQL(nodeTable string, status int) string {
	return fmt.Sprintf(
		`UPDATE %v SET status=%v WHERE id=%v;`,
		nodeTable, status, n.NID,
	)
}

func (n *Node) DeleteSQL(nodeTable string, ) string {
	return fmt.Sprintf(
		`DELETE FROM %v WHERE id=%v;`, nodeTable, n.NID,
	)
}

func (n *Node) Geometry() geom.Geometry {
	if n.geom != nil {
		return n.geom
	}
	n.geom = geom.ReadGeometry(n.WTK)
	return n.geom
}

func (n *Node) Polyline() pln.Polyline {
	if n.polyline.LineString != nil {
		return n.polyline
	}
	n.polyline = pln.CreatePolyline(n.Coordinates)
	return n.polyline
}

//Implements bbox interface
func (n *Node) BBox() *mbr.MBR {
	return n.Geometry().BBox()
}

//stringer interface
func (n *Node) String() string {
	return n.Geometry().WKT()
}

//number of coordinates
func (n *Node) Size() int {
	return n.Coordinates.Len()
}

//first point in coordinates
func (n *Node) First() *geom.Point {
	return n.Coordinates.Pt(0)
}

//last point in coordinates
func (n *Node) Last() *geom.Point {
	return n.Coordinates.Pt(n.Coordinates.Len()-1)
}

//subnode ids
func (n *Node) SubNodeIds() (string, string) {
	return fmt.Sprintf("%v/a", n.Id), fmt.Sprintf("%v/b", n.Id)
}

//as segment
func (n *Node) Segment() *geom.Segment {
	var i, j = 0, n.Coordinates.Len()-1
	return geom.NewSegment(n.Coordinates, i, j)
}

//hull segment as polyline
func (n *Node) SegmentAsPolyline() pln.Polyline {
	return pln.CreatePolyline(geom.Coordinates([]geom.Point{*n.First(), *n.Last()}))
}

//Is node collapsible with respect to other
//self and other should be contiguous
func (n *Node) Collapsible(other *Node) bool {
	//segments are already collapsed
	if n.Range.Size() == 1 {
		return true
	}
	//or hull can be a linear for
	//colinear boundaries where self.range.size > 1
	if _, ok := n.Geometry().(*geom.LineString); ok {
		return true
	}

	var ai, aj = n.First(), n.Last()
	var bi, bj = other.First(), other.Last()

	var c *geom.Point
	if ai.Equals2D(bi) || aj.Equals2D(bi) {
		c = bi
	} else if ai.Equals2D(bj) || aj.Equals2D(bj) {
		c = bj
	} else {
		return true
	}

	var t = bj
	if c.Equals2D(t) {
		t = bi
	}
	if ply, ok := n.Geometry().(*geom.Polygon); ok {
		return !ply.Shell.PointCompletelyInRing(t)
	}
	panic("unimplemented : hull type is handled")
}
