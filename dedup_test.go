package topojson

import (
	"reflect"
	"testing"

	"github.com/cheekybits/is"
	"github.com/chrisdinn/go.geojson"
)

// See https://github.com/mbostock/topojson/blob/master/test/topology/dedup-test.js

// dedup exact duplicate lines ABC & ABC share an arc
func TestDedupDuplicates(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("abc2", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc2")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 0)
	is.Equal(o2.Arc.End, 2)
	is.Nil(o2.Arc.Next)
}

// dedup reversed duplicate lines ABC & CBA share an arc
func TestDedupReversedDuplicates(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "cba")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 2)
	is.Equal(o2.Arc.End, 0)
	is.Nil(o2.Arc.Next)
}

// dedup exact duplicate rings ABCA & ABCA share an arc
func TestDedupDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})),
		NewTestFeature("abca2", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "abca2")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)
}

// dedup reversed duplicate rings ACBA & ABCA share an arc
func TestDedupReversedDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})),
		NewTestFeature("acba", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {2, 0}, {1, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "acba")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 3)
	is.Equal(o2.Arcs[0].End, 0)
	is.Nil(o2.Arcs[0].Next)
}

// dedup rotated duplicate rings BCAB & ABCA share an arc
func TestDedupRotatedDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})),
		NewTestFeature("bcab", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {2, 0}, {0, 0}, {1, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "bcab")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)
}

// dedup ring ABCA & line ABCA have no cuts
func TestDedupRingLine(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcaLine", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {0, 0},
		})),
		NewTestFeature("abcaPolygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcaLine")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abcaPolygon")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)
}

// dedup ring BCAB & line ABCA have no cuts
func TestDedupRingLineReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcaLine", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {0, 0},
		})),
		NewTestFeature("bcabPolygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {2, 0}, {0, 0}, {1, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcaLine")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "bcabPolygon")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)

	is.True(reflect.DeepEqual(topo.coordinates[4:8], [][]float64{
		{0, 0}, {1, 0}, {2, 0}, {0, 0},
	}))
}

// dedup ring ABCA & line BCAB have no cuts
func TestDedupRingLineReversed2(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("bcabLine", geojson.NewLineStringGeometry([][]float64{
			{1, 0}, {2, 0}, {0, 0}, {1, 0},
		})),
		NewTestFeature("abcaPolygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "bcabLine")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abcaPolygon")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)
}

// dedup when an old arc ABC extends a new arc AB, ABC is cut into AB-BC
func TestDedupOldArcExtends(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "ab")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 0)
	is.Equal(o2.Arc.End, 1)
	is.Nil(o2.Arc.Next)
}

// dedup when a reversed old arc CBA extends a new arc AB, CBA is cut into CB-BA
func TestDedupReversedOldArcExtends(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})),
		NewTestFeature("ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "ab")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 2)
	is.Equal(o2.Arc.End, 1)
	is.Nil(o2.Arc.Next)
}

// dedup when a new arc ADE shares its start with an old arc ABC, there are no cuts
func TestDedupNewArcSharesStart(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("ade", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 1}, {2, 1},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "ade")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 5)
	is.Nil(o2.Arc.Next)
}

// dedup ring ABA has no cuts
func TestDedupRingNoCuts(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("aba", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "aba")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 2)
	is.Nil(o1.Arcs[0].Next)
}

// dedup ring AA has no cuts
func TestDedupRingAANoCuts(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("aa", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "aa")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 1)
	is.Nil(o1.Arcs[0].Next)
}

// dedup degenerate ring A has no cuts
func TestDedupRingANoCuts(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("a", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "a")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 0)
	is.Nil(o1.Arcs[0].Next)
}

// dedup when a new line DEC shares its end with an old line ABC, there are no cuts
func TestDedupNewLineSharesEnd(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("dec", geojson.NewLineStringGeometry([][]float64{
			{0, 1}, {1, 1}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "dec")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 5)
	is.Nil(o2.Arc.Next)
}

// dedup when a new line ABC extends an old line AB, ABC is cut into AB-BC
func TestDedupNewLineExtends(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0},
		})),
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "ab")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 0)
	is.Equal(o2.Arc.End, 1)
	is.Equal(o2.Arc.Next.Start, 3)
	is.Equal(o2.Arc.Next.End, 4)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line ABC extends a reversed old line BA, ABC is cut into AB-BC
func TestDedupNewLineExtendsReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("ba", geojson.NewLineStringGeometry([][]float64{
			{1, 0}, {0, 0},
		})),
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "ba")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Nil(o1.Arc.Next)

	o2 := GetFeature(topo, "abc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 1)
	is.Equal(o2.Arc.End, 0)
	is.Equal(o2.Arc.Next.Start, 3)
	is.Equal(o2.Arc.Next.End, 4)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line starts BC in the middle of an old line ABC, ABC is cut into AB-BC
func TestDedupNewStartsMiddle(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("bc", geojson.NewLineStringGeometry([][]float64{
			{1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "bc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 1)
	is.Equal(o2.Arc.End, 2)
	is.Nil(o2.Arc.Next)
}

// dedup when a new line BC starts in the middle of a reversed old line CBA, CBA is cut into CB-BA
func TestDedupNewStartsMiddleReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})),
		NewTestFeature("bc", geojson.NewLineStringGeometry([][]float64{
			{1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "bc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 1)
	is.Equal(o2.Arc.End, 0)
	is.Nil(o2.Arc.Next)
}

// dedup when a new line ABD deviates from an old line ABC, ABD is cut into AB-BD and ABC is cut into AB-BC
func TestDedupNewLineDeviates(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("abd", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {3, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "abd")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 0)
	is.Equal(o2.Arc.End, 1)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line ABD deviates from a reversed old line CBA, CBA is cut into CB-BA and ABD is cut into AB-BD
func TestDedupNewLineDeviatesReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})),
		NewTestFeature("abd", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {3, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "abd")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 2)
	is.Equal(o2.Arc.End, 1)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line DBC merges into an old line ABC, DBC is cut into DB-BC and ABC is cut into AB-BC
func TestDedupNewLineMerges(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("dbc", geojson.NewLineStringGeometry([][]float64{
			{3, 0}, {1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "dbc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 1)
	is.Equal(o2.Arc.Next.End, 2)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line DBC merges into a reversed old line CBA, DBC is cut into DB-BC and CBA is cut into CB-BA
func TestDedupNewLineMergesReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})),
		NewTestFeature("dbc", geojson.NewLineStringGeometry([][]float64{
			{3, 0}, {1, 0}, {2, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "cba")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "dbc")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 1)
	is.Equal(o2.Arc.Next.End, 0)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line DBE shares a single midpoint with an old line ABC, DBE is cut into DB-BE and ABC is cut into AB-BC
func TestDedupNewLineSharesMidpoint(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("dbe", geojson.NewLineStringGeometry([][]float64{
			{0, 1}, {1, 0}, {2, 1},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abc")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	is.Equal(o1.Arc.Next.Start, 1)
	is.Equal(o1.Arc.Next.End, 2)
	is.Nil(o1.Arc.Next.Next)

	o2 := GetFeature(topo, "dbe")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 3)
	is.Equal(o2.Arc.End, 4)
	is.Equal(o2.Arc.Next.Start, 4)
	is.Equal(o2.Arc.Next.End, 5)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a new line ABDE skips a point with an old line ABCDE, ABDE is cut into AB-BD-DE and ABCDE is cut into AB-BCD-DE
func TestDedupNewLineSkipsPoint(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcde", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
		})),
		NewTestFeature("adbe", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {3, 0}, {4, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcde")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	o1Next := o1.Arc.Next
	is.Equal(o1Next.Start, 1)
	is.Equal(o1Next.End, 3)
	is.Equal(o1Next.Next.Start, 3)
	is.Equal(o1Next.Next.End, 4)
	is.Nil(o1Next.Next.Next)

	o2 := GetFeature(topo, "adbe")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 0)
	is.Equal(o2.Arc.End, 1)
	o2Next := o2.Arc.Next
	is.Equal(o2Next.Start, 6)
	is.Equal(o2Next.End, 7)
	is.Equal(o2Next.Next.Start, 3)
	is.Equal(o2Next.Next.End, 4)
	is.Nil(o2Next.Next.Next)
}

// dedup when a new line ABDE skips a point with a reversed old line EDCBA, ABDE is cut into AB-BD-DE and EDCBA is cut into ED-DCB-BA
func TestDedupNewLineSkipsPointReversed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("edcba", geojson.NewLineStringGeometry([][]float64{
			{4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0},
		})),
		NewTestFeature("adbe", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {3, 0}, {4, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "edcba")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	o1Next := o1.Arc.Next
	is.Equal(o1Next.Start, 1)
	is.Equal(o1Next.End, 3)
	is.Equal(o1Next.Next.Start, 3)
	is.Equal(o1Next.Next.End, 4)
	is.Nil(o1Next.Next.Next)

	o2 := GetFeature(topo, "adbe")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 4)
	is.Equal(o2.Arc.End, 3)
	o2Next := o2.Arc.Next
	is.Equal(o2Next.Start, 6)
	is.Equal(o2Next.End, 7)
	is.Equal(o2Next.Next.Start, 1)
	is.Equal(o2Next.Next.End, 0)
	is.Nil(o2Next.Next.Next)
}

// dedup when a line ABCDBE self-intersects with its middle, it is not cut
func TestDedupSelfIntersectsMiddle(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbe", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {3, 0}, {1, 0}, {4, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcdbe")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 5)
	is.Nil(o1.Arc.Next)
}

// dedup when a line ABACD self-intersects with its start, it is cut into ABA-ACD
func TestDedupSelfIntersectsStart(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abacd", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {0, 0}, {3, 0}, {4, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abacd")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Equal(o1.Arc.Next.Start, 2)
	is.Equal(o1.Arc.Next.End, 4)
	is.Nil(o1.Arc.Next.Next)
}

// dedup when a line ABDCD self-intersects with its end, it is cut into ABD-DCD
func TestDedupSelfIntersectsEnd(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbd", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {4, 0}, {3, 0}, {4, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcdbd")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 2)
	is.Equal(o1.Arc.Next.Start, 2)
	is.Equal(o1.Arc.Next.End, 4)
	is.Nil(o1.Arc.Next.Next)
}

// dedup when an old line ABCDBE self-intersects and shares a point B, ABCDBE is cut into AB-BCDB-BE and FBG is cut into FB-BG
func TestDedupSelfIntersectsShares(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcdbe", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {3, 0}, {1, 0}, {4, 0},
		})),
		NewTestFeature("fbg", geojson.NewLineStringGeometry([][]float64{
			{0, 1}, {1, 0}, {2, 1},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcdbe")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 1)
	o1Next := o1.Arc.Next
	is.Equal(o1Next.Start, 1)
	is.Equal(o1Next.End, 4)
	is.Equal(o1Next.Next.Start, 4)
	is.Equal(o1Next.Next.End, 5)
	is.Nil(o1Next.Next.Next)

	o2 := GetFeature(topo, "fbg")
	is.Equal(o2.Type, geojson.GeometryLineString)
	is.Equal(o2.Arc.Start, 6)
	is.Equal(o2.Arc.End, 7)
	is.Equal(o2.Arc.Next.Start, 7)
	is.Equal(o2.Arc.Next.End, 8)
	is.Nil(o2.Arc.Next.Next)
}

// dedup when a line ABCA is closed, there are no cuts
func TestDedupLineClosed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {0, 1}, {0, 0},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryLineString)
	is.Equal(o1.Arc.Start, 0)
	is.Equal(o1.Arc.End, 3)
	is.Nil(o1.Arc.Next)
}

// dedup when a ring ABCA is closed, there are no cuts
func TestDedupRingClosed(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)
}

// dedup exact duplicate rings ABCA & ABCA have no cuts
func TestDedupDuplicateRingsShare(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
		NewTestFeature("abca2", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "abca2")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)
}

// dedup reversed duplicate rings ABCA & ACBA have no cuts
func TestDedupDuplicateRingsReversedShare(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
		NewTestFeature("acba", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {0, 1}, {1, 0}, {0, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "acba")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 3)
	is.Equal(o2.Arcs[0].End, 0)
	is.Nil(o2.Arcs[0].Next)
}

// dedup coincident rings ABCA & BCAB have no cuts
func TestDedupCoincidentRings(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
		NewTestFeature("bcab", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {0, 1}, {0, 0}, {1, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "bcab")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 0)
	is.Equal(o2.Arcs[0].End, 3)
	is.Nil(o2.Arcs[0].Next)
}

// dedup coincident reversed rings ABCA & BACB have no cuts
func TestDedupCoincidentRings2(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
		NewTestFeature("bacb", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {0, 0}, {0, 1}, {1, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "bacb")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 3)
	is.Equal(o2.Arcs[0].End, 0)
	is.Nil(o2.Arcs[0].Next)
}

// dedup coincident rings ABCDA, EFAE & GHCG are cut into ABC-CDA, EFAE and GHCG
func TestDedupCoincidentRings3(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcda", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0},
			},
		})),
		NewTestFeature("efae", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, -1}, {1, -1}, {0, 0}, {0, -1},
			},
		})),
		NewTestFeature("ghcg", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 2}, {1, 2}, {1, 1}, {0, 2},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcda")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 2)
	is.Equal(o1.Arcs[0].Next.Start, 2)
	is.Equal(o1.Arcs[0].Next.End, 4)
	is.Nil(o1.Arcs[0].Next.Next)

	o2 := GetFeature(topo, "efae")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 5)
	is.Equal(o2.Arcs[0].End, 8)
	is.Nil(o2.Arcs[0].Next)

	o3 := GetFeature(topo, "ghcg")
	is.Equal(o3.Type, geojson.GeometryPolygon)
	is.Equal(len(o3.Arcs), 1)
	is.Equal(o3.Arcs[0].Start, 9)
	is.Equal(o3.Arcs[0].End, 12)
	is.Nil(o3.Arcs[0].Next)
}

// dedup coincident rings ABCA & DBED have no cuts, but are rotated to share B
func TestDedupNoCutsButRotated(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 1}, {0, 0},
			},
		})),
		NewTestFeature("dbed", geojson.NewPolygonGeometry([][][]float64{
			{
				{2, 1}, {1, 0}, {2, 2}, {2, 1},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abca")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 3)
	is.Nil(o1.Arcs[0].Next)

	o2 := GetFeature(topo, "dbed")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 4)
	is.Equal(o2.Arcs[0].End, 7)
	is.Nil(o2.Arcs[0].Next)

	is.True(reflect.DeepEqual(topo.coordinates[0:4], [][]float64{
		{1, 0}, {0, 1}, {0, 0}, {1, 0},
	}))
	is.True(reflect.DeepEqual(topo.coordinates[4:8], [][]float64{
		{1, 0}, {2, 2}, {2, 1}, {1, 0},
	}))
}

// dedup overlapping rings ABCDA and BEFCB are cut into BC-CDAB and BEFC-CB
func TestDedupOverlapping(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("abcda", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}, // rotated to BCDAB, cut BC-CDAB
			},
		})),
		NewTestFeature("befcb", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {2, 0}, {2, 1}, {1, 1}, {1, 0},
			},
		})),
	}

	topo := &Topology{input: in}
	topo.extract()
	topo.cut()
	topo.dedup()

	o1 := GetFeature(topo, "abcda")
	is.Equal(o1.Type, geojson.GeometryPolygon)
	is.Equal(len(o1.Arcs), 1)
	is.Equal(o1.Arcs[0].Start, 0)
	is.Equal(o1.Arcs[0].End, 1)
	is.Equal(o1.Arcs[0].Next.Start, 1)
	is.Equal(o1.Arcs[0].Next.End, 4)
	is.Nil(o1.Arcs[0].Next.Next)

	o2 := GetFeature(topo, "befcb")
	is.Equal(o2.Type, geojson.GeometryPolygon)
	is.Equal(len(o2.Arcs), 1)
	is.Equal(o2.Arcs[0].Start, 5)
	is.Equal(o2.Arcs[0].End, 8)
	is.Equal(o2.Arcs[0].Next.Start, 1)
	is.Equal(o2.Arcs[0].Next.End, 0)
	is.Nil(o2.Arcs[0].Next.Next)
}
