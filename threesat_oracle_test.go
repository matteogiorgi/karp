package karp

import "testing"

// TestThreeSATDIMACS checks the exact text DIMACS produces against a
// hand-written expectation, since every downstream consumer — gophersat
// today, potentially any other DIMACS-reading solver tomorrow, per the
// swappable-oracle argument in docs/06-pipeline-architecture.md, §6.4 —
// depends on this format being exactly right, not just "close enough for
// gophersat's parser to guess correctly".
func TestThreeSATDIMACS(t *testing.T) {
	phi := ThreeSAT{
		NumVars: 3,
		Clauses: []Clause3{
			{1, 2, 3},
			{-1, -2, 3},
		},
	}

	want := "p cnf 3 2\n1 2 3 0\n-1 -2 3 0\n"
	if got := phi.DIMACS(); got != want {
		t.Errorf("DIMACS() = %q, want %q", got, want)
	}
}

// TestSolveThreeSATSatisfiable checks the "yes" side of the real oracle on
// a formula already known satisfiable from earlier tests in this package,
// and — critically — that the model gophersat returns actually passes
// this package's own Verify, not just that gophersat claims Sat. A correct
// SAT status paired with a wrongly re-indexed model would be exactly the
// kind of bug this check exists to catch.
func TestSolveThreeSATSatisfiable(t *testing.T) {
	// c0: (x1 ∨ x2 ∨ x2)   c1: (¬x1 ∨ x2 ∨ x2)  -- satisfiable, e.g. x2=true
	phi := ThreeSAT{
		NumVars: 2,
		Clauses: []Clause3{
			{1, 2, 2},
			{-1, 2, 2},
		},
	}

	ok, assign := SolveThreeSAT(phi)
	if !ok {
		t.Fatal("expected phi to be satisfiable")
	}
	if !phi.Verify(assign) {
		t.Fatalf("model %v returned by gophersat does not satisfy phi", assign)
	}
}

// TestSolveThreeSATUnusedVariable checks a case with real consequences for
// property-based testing: NumVars declares more variables than actually
// occur in any clause (x2 and x3 are declared but never referenced). A
// generator that picks NumVars and clause literals independently — genThreeSAT
// in property_test.go does exactly this — is very likely to produce
// instances shaped like this. Verified empirically before writing this
// test, ahead of genThreeSAT existing: gophersat's Model
// covers every variable declared in the DIMACS header, not just the ones
// a clause happens to mention, so re-indexing it into Assignment is safe
// here — but that was worth confirming rather than assuming, since
// SolveThreeSAT would panic on an index out of range if it were not true.
func TestSolveThreeSATUnusedVariable(t *testing.T) {
	phi := ThreeSAT{
		NumVars: 3,
		Clauses: []Clause3{{1, 1, 1}},
	}

	ok, assign := SolveThreeSAT(phi)
	if !ok {
		t.Fatal("expected phi to be satisfiable")
	}
	if len(assign) != phi.NumVars+1 {
		t.Fatalf("len(assign) = %d, want %d", len(assign), phi.NumVars+1)
	}
	if !phi.Verify(assign) {
		t.Fatalf("model %v does not satisfy phi", assign)
	}
}

// TestSolveThreeSATUnsatisfiable checks the "no" side on a formula that is
// unsatisfiable by construction: one clause forces x1 true, padded to
// width 3 by duplicating the literal (docs/05-cook-levin-and-practice.md,
// §5.2), and the other forces x1 false the same way. No assignment can
// satisfy both.
func TestSolveThreeSATUnsatisfiable(t *testing.T) {
	phi := ThreeSAT{
		NumVars: 1,
		Clauses: []Clause3{
			{1, 1, 1},
			{-1, -1, -1},
		},
	}

	if ok, assign := SolveThreeSAT(phi); ok {
		t.Fatalf("expected phi to be unsatisfiable, got model %v", assign)
	}
}

// TestSolveThreeSATAgreesWithIndependentSetOracle is the first end-to-end
// check spanning both halves of the pipeline described in docs/06-
// pipeline-architecture.md, §6.3: the real SAT oracle deciding the
// original ThreeSAT instance, and the brute-force oracle deciding its
// image under ThreeSATToIndependentSet, must agree on satisfiability — and
// the certificate map g must turn the brute-force oracle's witness into
// one the real oracle's own formula accepts. This is exactly the "yes"
// side of the four-check pipeline, run once by hand here on a fixed
// instance; TestThreeSATToIndependentSetProperty (property_test.go)
// repeats it on hundreds of random ones.
func TestSolveThreeSATAgreesWithIndependentSetOracle(t *testing.T) {
	phi := ThreeSAT{
		NumVars: 2,
		Clauses: []Clause3{
			{1, 2, 2},
			{-1, 2, 2},
		},
	}

	satOK, _ := SolveThreeSAT(phi)
	isOK, cert := SolveIndependentSet(ThreeSATToIndependentSet(phi))

	if satOK != isOK {
		t.Fatalf("SolveThreeSAT = %v, SolveIndependentSet(f(phi)) = %v: invariant broken", satOK, isOK)
	}
	if isOK {
		assign := CertificateToAssignment(phi, cert)
		if !phi.Verify(assign) {
			t.Fatalf("assignment %v decoded from %v does not satisfy phi", assign, cert)
		}
	}
}
