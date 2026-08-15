package karp

import "testing"

// TestThreeSATToIndependentSetStructure checks the shape of the graph
// ThreeSATToIndependentSet builds against a hand-worked two-clause
// formula: the vertex count, K, the triangle edges within each clause,
// and the one conflict edge that the single complementary pair (x1, ¬x1)
// across the two clauses should produce.
func TestThreeSATToIndependentSetStructure(t *testing.T) {
	// c0: (x1 ∨ x2 ∨ x2)   c1: (¬x1 ∨ x2 ∨ x2)
	phi := ThreeSAT{
		NumVars: 2,
		Clauses: []Clause3{
			{1, 2, 2},
			{-1, 2, 2},
		},
	}

	inst := ThreeSATToIndependentSet(phi)

	if inst.G.NumVertices != 6 {
		t.Fatalf("NumVertices = %d, want 6 (3 per clause)", inst.G.NumVertices)
	}
	if inst.K != 2 {
		t.Fatalf("K = %d, want 2 (one clause per vertex needed)", inst.K)
	}

	edgeSet := make(map[[2]int]bool)
	for _, e := range inst.G.Edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		edgeSet[e] = true
	}

	// The two triangles: {0,1,2} for clause 0, {3,4,5} for clause 1.
	triangleEdges := [][2]int{{0, 1}, {0, 2}, {1, 2}, {3, 4}, {3, 5}, {4, 5}}
	for _, e := range triangleEdges {
		if !edgeSet[e] {
			t.Errorf("missing expected triangle edge %v", e)
		}
	}

	// Vertex 0 is clause 0's literal x1; vertex 3 is clause 1's literal
	// ¬x1 — the only complementary pair across the two clauses.
	if !edgeSet[[2]int{0, 3}] {
		t.Error("missing expected conflict edge {0,3} for the complementary pair (x1, ¬x1)")
	}

	// No other cross-clause edge should exist: every other literal pair
	// across the two clauses is either non-complementary (x2 vs x2) or
	// unrelated (x1 vs x2), so the edge count should be exactly the two
	// triangles (3+3) plus the one conflict edge.
	if got, want := len(inst.G.Edges), 7; got != want {
		t.Errorf("len(Edges) = %d, want %d", got, want)
	}
}

// TestThreeSATToIndependentSetRoundTrip exercises f and g together on the
// same formula: it picks an independent set by hand (one vertex per
// clause, no two naming complementary literals), confirms it really is a
// valid certificate via IndependentSet.Verify, decodes it with
// CertificateToAssignment, and confirms the decoded assignment really
// satisfies the original formula via ThreeSAT.Verify. This is the "yes"
// side of the four-check pipeline described in docs/06-pipeline-
// architecture.md, §6.3, run here on one hand-picked instance instead of
// the hundreds of random ones TestThreeSATToIndependentSetProperty
// (property_test.go) runs against the real SAT solver.
func TestThreeSATToIndependentSetRoundTrip(t *testing.T) {
	// c0: (x1 ∨ x2 ∨ x2)   c1: (¬x1 ∨ x2 ∨ x2)
	phi := ThreeSAT{
		NumVars: 2,
		Clauses: []Clause3{
			{1, 2, 2},
			{-1, 2, 2},
		},
	}
	inst := ThreeSATToIndependentSet(phi)

	// Vertex 0 = clause 0's x1, vertex 4 = clause 1's second x2: two
	// different clauses, two different variables, no conflict edge
	// between them.
	cert := VertexSet{0, 4}

	if !inst.Verify(cert) {
		t.Fatalf("%v is not accepted as a valid independent set of %+v", cert, inst)
	}

	assign := CertificateToAssignment(phi, cert)
	if !phi.Verify(assign) {
		t.Fatalf("assignment %v decoded from %v does not satisfy phi", assign, cert)
	}
}
