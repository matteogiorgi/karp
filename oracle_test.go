package karp

import "testing"

// TestBruteForceOracleFindsExistingSet checks the "yes" side of the
// oracle on the same 4-cycle used elsewhere in this package: {0,2} is a
// maximum independent set, so a K=2 instance must be found satisfiable,
// and the certificate the oracle returns must itself pass Verify — the
// oracle is not allowed to just say "yes", it has to produce the witness.
func TestBruteForceOracleFindsExistingSet(t *testing.T) {
	g := Graph{
		NumVertices: 4,
		Edges:       [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}},
	}
	inst := IndependentSet{g, 2}

	ok, cert := SolveIndependentSet(inst)
	if !ok {
		t.Fatal("expected a K=2 independent set to be found in a 4-cycle")
	}
	if !inst.Verify(cert) {
		t.Fatalf("oracle returned %v, which does not itself pass Verify", cert)
	}
}

// TestBruteForceOracleRejectsMissingSet checks the "no" side, which none
// of this package's earlier tests could establish by hand-picking a
// certificate: a complete graph on 4 vertices has no two non-adjacent
// vertices at all, so no independent set of size 2 exists, and exhausting
// every one of the 2^4 candidate subsets is the only way to be sure.
func TestBruteForceOracleRejectsMissingSet(t *testing.T) {
	g := Graph{
		NumVertices: 4,
		Edges: [][2]int{
			{0, 1}, {0, 2}, {0, 3},
			{1, 2}, {1, 3},
			{2, 3},
		},
	}
	inst := IndependentSet{g, 2}

	if ok, cert := SolveIndependentSet(inst); ok {
		t.Fatalf("expected no independent set of size 2 in K4, oracle returned %v", cert)
	}
}

// TestBruteForceOracleUniverseTooLarge checks that the safety cap actually
// fires instead of silently iterating for an impractically long time.
func TestBruteForceOracleUniverseTooLarge(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected BruteForceOracle to panic on a universe above maxBruteForceUniverse")
		}
	}()
	inst := IndependentSet{Graph{NumVertices: maxBruteForceUniverse + 1}, 1}
	SolveIndependentSet(inst)
}

// TestBruteForceOracleAgreesWithReduction ties the oracle to
// ThreeSATToIndependentSet: for a small satisfiable formula, the oracle
// must find an independent set in f(phi), and decoding that independent
// set with g = CertificateToAssignment must produce an assignment that
// satisfies phi — the same "yes" side of the four-check pipeline in
// docs/06-pipeline-architecture.md, §6.3, but now with the certificate
// discovered by the oracle instead of chosen by hand.
func TestBruteForceOracleAgreesWithReduction(t *testing.T) {
	// c0: (x1 ∨ x2 ∨ x2)   c1: (¬x1 ∨ x2 ∨ x2)  -- satisfiable, e.g. x2=true
	phi := ThreeSAT{
		NumVars: 2,
		Clauses: []Clause3{
			{1, 2, 2},
			{-1, 2, 2},
		},
	}
	inst := ThreeSATToIndependentSet(phi)

	ok, cert := SolveIndependentSet(inst)
	if !ok {
		t.Fatal("expected the reduced instance of a satisfiable formula to have an independent set")
	}

	assign := CertificateToAssignment(phi, cert)
	if !phi.Verify(assign) {
		t.Fatalf("assignment %v decoded from oracle certificate %v does not satisfy phi", assign, cert)
	}
}
