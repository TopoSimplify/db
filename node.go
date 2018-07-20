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
)

//Node
type Node struct {
	Id          string
	FID         int
	NID         int
	Coordinates []geom.Point
	Range       rng.Rng
	HullType    geom.GeoType
	WTK         string
	geom        geom.Geometry
	polyline    *pln.Polyline
}

func NewDBNode(coordinates []geom.Point, r rng.Rng, fid int, gfn geom.GeometryFn, ids ...string) *Node {
	var id string
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = random.String(8)
	}
	var n = NewDBNodeFromDPNode(node.New(coordinates, r, gfn, id))
	n.FID = fid
	return n
}

func NewDBNodeFromDPNode(node *node.Node) *Node {
	return &Node{
		Id:          node.Id(),
		Coordinates: node.Polyline.Coordinates,
		Range:       node.Range,
		WTK:         node.Geometry.WKT(),
		HullType:    node.Geometry.Type(),
		geom:        node.Geometry,
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

//number of coordinates
func (n *Node) Size() int {
	return len(n.Coordinates)
}

//first point in coordinates
func (n *Node) First() geom.Point {
	return n.Coordinates[0]
}

//last point in coordinates
func (n *Node) Last() geom.Point {
	return n.Coordinates[len(n.Coordinates)-1]
}

//subnode ids
func (n *Node) SubNodeIds() (string, string) {
	return fmt.Sprintf("%v/a", n.Id), fmt.Sprintf("%v/b", n.Id)
}

//as segment
func (n *Node) Segment() *seg.Seg {
	var i, j = 0, len(n.Coordinates)-1
	return seg.NewSeg(&n.Coordinates[i], &n.Coordinates[j], n.Range.I, n.Range.J)
}

//hull segment as polyline
func (n *Node) SegmentAsPolyline() *pln.Polyline {
	var i, j = 0, len(n.Coordinates)-1
	return pln.New([]geom.Point{n.Coordinates[i], n.Coordinates[j]})
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

	var ai, aj = &n.Coordinates[0], &n.Coordinates[n.Size()-1]
	var bi, bj = &other.Coordinates[0], &other.Coordinates[other.Size()-1]

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
