package db

import (
	"simplex/node"
	"github.com/intdxdt/geom"
	"simplex/pln"
	"simplex/rng"
	"simplex/seg"
	"github.com/intdxdt/mbr"
	"fmt"
	"github.com/intdxdt/random"
)

const NodeTblColumns = "fid, gob, geom"

//Node
type Node struct {
	Id          string
	FID         int
	NID         int
	Part        int
	Coordinates []*geom.Point
	Range       *rng.Range
	HullType    int
	WTK         string
	geom        geom.Geometry
	polyline    *pln.Polyline
}

func (n *Node) Geometry() geom.Geometry {
	if n.geom != nil {
		return n.geom
	}
	n.geom = geom.NewGeometry(n.WTK)
	return n.geom
}

func (n *Node) Polyline() *pln.Polyline {
	if n.polyline != nil {
		return n.polyline
	}
	n.polyline = pln.New(n.Coordinates)
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

//first point in coordinates
func (n *Node) First() *geom.Point {
	return n.Coordinates[0]
}

//last point in coordinates
func (n *Node) Last() *geom.Point {
	return n.Coordinates[len(n.Coordinates)-1]
}

//subnode ids
func (self *Node) SubNodeIds() (string, string) {
	return fmt.Sprintf("%v/a", self.Id), fmt.Sprintf("%v/b", self.Id)
}

//as segment
func (n *Node) Segment() *seg.Seg {
	var a, b = n.SegmentPoints()
	return seg.NewSeg(a, b, n.Range.I, n.Range.J)
}

//hull segment as polyline
func (n *Node) SegmentAsPolyline() *pln.Polyline {
	var a, b = n.SegmentPoints()
	return pln.New([]*geom.Point{a, b})
}

//segment points
func (n *Node) SegmentPoints() (*geom.Point, *geom.Point) {
	return n.First(), n.Last()
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

	var ai, aj = n.SegmentPoints()
	var bi, bj = other.SegmentPoints()

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

func New(coordinates []*geom.Point, r *rng.Range, fid, part int, gfn geom.GeometryFn, ids ...string) *Node {
	var id string
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = random.String(8)
	}
	var n = node.New(coordinates, r, gfn, id)
	var dn = NewDBNode(n)
	dn.FID, dn.Part = fid, part
	return dn
}

func NewDBNode(node *node.Node) *Node {
	return &Node{
		Id:          node.Id(),
		Coordinates: node.Polyline.Coordinates,
		Range:       node.Range,
		WTK:         node.Geom.WKT(),
		HullType:    node.Geom.Type().Value(),
	}
}

func NewNodeFromDB(ndb *Node) *node.Node {
	var n = &node.Node{
		Polyline: pln.New(ndb.Coordinates),
		Range:    ndb.Range,
		Geom:     geom.NewGeometry(ndb.WTK),
	}
	n.SetId(ndb.Id)
	return n
}
