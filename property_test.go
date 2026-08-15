package karp

import (
	"testing"

	"pgregory.net/rapid"
)

// genThreeSAT is a rapid generator for random ThreeSAT instances. It is
// the only generator this package needs, per docs/06-pipeline-
// architecture.md, §6.1: ThreeSAT is the only instance type ever produced
// directly rather than born as the image of a reduction, so it is the
// only type whose test cases have to be generated rather than derived.
//
// numVars and numClauses are drawn independently of each other and of the
// literals inside each clause, deliberately: this is exactly what makes
// it likely, not just possible, for a generated formula to declare a
// variable that no clause ends up mentioning — the case
// TestSolveThreeSATUnusedVariable in threesat_oracle_test.go verified
// SolveThreeSAT ahead of time, precisely because this generator was going
// to produce it sooner or later.
//
// The two ranges below are not arbitrary; they come from two measurements
// made before settling on them, not from a guess:
//
//   - Uniform random 3-SAT with few clauses relative to variables is
//     almost always satisfiable. An early version of this generator, drawn
//     from IntRange(1,5) for both numVars and numClauses independently,
//     produced satisfiable instances 97.5% of the time (measured over
//     2000 samples) — leaving the "no" side of the invariant in step 3 of
//     TestThreeSATToIndependentSetProperty below exercised only two or
//     three times per default 100-check run. Capping numVars at 3 and
//     raising the floor on numClauses to 3 brings that down to roughly
//     92% satisfiable / 8% unsatisfiable, a far healthier mix for the
//     same sample size.
//   - BruteForceOracle's cost is exponential in IndependentSet's vertex
//     count (3 per clause, docs/06-pipeline-architecture.md, §6.1), and
//     that exponential is real, not theoretical: a worst-case (genuinely
//     unsatisfiable, forcing exhaustive enumeration) instance measured
//     13ms at 5 clauses, 140ms at 6, 1.35s at 7, 13s at 8, and over two
//     minutes at 9 — an order of magnitude per extra clause, not the
//     factor of 8 raw subset counting alone would suggest, because each
//     candidate also costs O(edges) to check. maxBruteForceUniverse in
//     oracle.go only guards against an outright hang (up to 30 clauses);
//     it does not by itself guarantee a *fast* test suite. numClauses is
//     capped at 6 here specifically because that is where the worst case
//     was still comfortably under a fifth of a second, well before rapid's
//     shrinking — which can re-run the property function many times over
//     — would make a slow case actually painful to wait out.
func genThreeSAT(t *rapid.T) ThreeSAT {
	numVars := rapid.IntRange(1, 3).Draw(t, "numVars")
	numClauses := rapid.IntRange(3, 6).Draw(t, "numClauses")

	clauses := make([]Clause3, numClauses)
	for i := range clauses {
		for j := 0; j < 3; j++ {
			v := rapid.IntRange(1, numVars).Draw(t, "var")
			sign := rapid.SampledFrom([]int{1, -1}).Draw(t, "sign")
			clauses[i][j] = Literal(sign * v)
		}
	}

	return ThreeSAT{NumVars: numVars, Clauses: clauses}
}

// TestThreeSATToIndependentSetProperty is the property test the whole
// pipeline has been building toward: the four-check data flow described
// in docs/06-pipeline-architecture.md, §6.3, run by rapid on hundreds of
// randomly generated 3-SAT instances instead of the handful chosen by
// hand in threesat_to_independent_set_test.go and oracle_test.go.
//
//  1. Construction: y = f(x). No check here, just execution of the
//     reduction — ThreeSATToIndependentSet.
//  2. Double, independent decision: the real SAT oracle decides phi, the
//     brute-force oracle decides f(phi). Neither knows the other exists.
//  3. Invariant A(x) = B(f(x)): the two boolean answers must agree. This
//     is the correctness of f as a map that preserves the yes/no answer
//     (docs/03-polynomial-many-one-reductions.md, §3.1).
//  4. Certificate: if satisfiable, g = CertificateToAssignment must turn
//     the brute-force oracle's witness into one phi's own Verify accepts.
//     This is checked independently of step 3 — a bug in g that produces
//     a wrong certificate would not make the invariant in step 3 fail
//     (that one only concerns yes/no), but it would fail here.
//
// If a case fails either check, rapid automatically shrinks phi to a
// minimal counterexample before reporting it.
func TestThreeSATToIndependentSetProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		phi := genThreeSAT(t)

		inst := ThreeSATToIndependentSet(phi)

		satOK, _ := SolveThreeSAT(phi)
		isOK, cert := SolveIndependentSet(inst)

		if satOK != isOK {
			t.Fatalf("invariant broken: SolveThreeSAT(phi) = %v, SolveIndependentSet(f(phi)) = %v, phi = %+v",
				satOK, isOK, phi)
		}

		if isOK {
			assign := CertificateToAssignment(phi, cert)
			if !phi.Verify(assign) {
				t.Fatalf("g produced an invalid certificate: cert = %v, decoded assign = %v, phi = %+v",
					cert, assign, phi)
			}
		}
	})
}
