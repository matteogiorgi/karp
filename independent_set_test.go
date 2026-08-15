package karp

import "testing"

// TestIndependentSetVerify checks Verify on a 4-cycle, a graph small
// enough that every independent set in it can be enumerated by hand: the
// two maximum independent sets are the opposite-corner pairs {0,2} and
// {1,3}, and every adjacent pair is disqualified by definition. It also
// covers a certificate that names a vertex outside the graph — such a
// certificate is not a subset of V at all, so it must be rejected even if
// padding it with out-of-range "vertices" would otherwise reach K.
func TestIndependentSetVerify(t *testing.T) {
	// A 4-cycle: 0-1-2-3-0. Independent sets: {0,2}, {1,3} (size 2, max).
	g := Graph{
		NumVertices: 4,
		Edges:       [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}},
	}

	cases := []struct {
		name string
		inst IndependentSet
		cert VertexSet
		want bool
	}{
		{"opposite corners size 2", IndependentSet{g, 2}, VertexSet{0, 2}, true},
		{"adjacent vertices are not independent", IndependentSet{g, 2}, VertexSet{0, 1}, false},
		{"too small for required K", IndependentSet{g, 3}, VertexSet{0, 2}, false},
		{"single vertex always independent", IndependentSet{g, 1}, VertexSet{1}, true},
		{"out-of-range vertex must not count toward K", IndependentSet{g, 2}, VertexSet{0, 99}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.inst.Verify(c.cert); got != c.want {
				t.Errorf("Verify(%v) = %v, want %v", c.cert, got, c.want)
			}
		})
	}
}
