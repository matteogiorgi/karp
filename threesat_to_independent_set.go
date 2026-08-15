package karp

// ThreeSATToIndependentSet is the witness function f for the reduction
// 3-SAT ≤p Independent Set, in the sense of docs/03-polynomial-many-one-
// reductions.md, §3.1: given a 3-SAT instance, it builds — in time
// polynomial in the size of phi — an Independent Set instance that has a
// solution if and only if phi is satisfiable.
//
// The construction is the standard textbook gadget. For every clause,
// three vertices are created, one per literal occurrence in that clause
// (vertexIndex(c, p) names the p-th literal of the c-th clause). Two
// kinds of edges are added:
//
//   - a triangle inside every clause, connecting all three of its
//     vertices to one another. Since the three vertices of a triangle are
//     pairwise adjacent, an independent set can contain at most one of
//     them — at most one literal "wins" per clause.
//   - a conflict edge between any two vertices, in different clauses,
//     that represent complementary literals (x and ¬x). This rules out an
//     independent set from "winning" a variable both true and false at
//     once.
//
// K is set to the number of clauses. With exactly one vertex allowed per
// triangle, and exactly that many triangles, an independent set of size K
// is only reachable by picking exactly one vertex from every clause — one
// satisfied literal per clause, consistently — which is what makes the
// reduction correct: phi is satisfiable exactly when such a set exists.
// The "yes ⟹ yes" half of that argument is made constructive by
// CertificateToAssignment below; the "no ⟹ no" half is a universal claim
// about every instance, argued here in the doc comment and corroborated,
// not replaced, by TestThreeSATToIndependentSetProperty (property_test.go)
// against a real SAT solver (see the distinction drawn in docs/05-cook-
// levin-and-practice.md, §5.4).
func ThreeSATToIndependentSet(phi ThreeSAT) IndependentSet {
	n := len(phi.Clauses)
	var edges [][2]int

	for c := 0; c < n; c++ {
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				edges = append(edges, [2]int{vertexIndex(c, i), vertexIndex(c, j)})
			}
		}
	}

	for c1 := 0; c1 < n; c1++ {
		for p1 := 0; p1 < 3; p1++ {
			for c2 := c1 + 1; c2 < n; c2++ {
				for p2 := 0; p2 < 3; p2++ {
					if phi.Clauses[c1][p1] == -phi.Clauses[c2][p2] {
						edges = append(edges, [2]int{vertexIndex(c1, p1), vertexIndex(c2, p2)})
					}
				}
			}
		}
	}

	return IndependentSet{
		G: Graph{NumVertices: 3 * n, Edges: edges},
		K: n,
	}
}

// vertexIndex maps a (clause, position-within-clause) pair to the vertex
// index ThreeSATToIndependentSet assigns it: clause c contributes vertices
// 3*c, 3*c+1, 3*c+2, one per literal position. clauseAndPosition below is
// its inverse.
func vertexIndex(clause, position int) int {
	return clause*3 + position
}

// clauseAndPosition is the inverse of vertexIndex: it recovers which
// clause and which literal position a vertex of the constructed graph
// stands for. CertificateToAssignment uses it to decode a chosen vertex
// back into the literal it names.
func clauseAndPosition(vertex int) (clause, position int) {
	return vertex / 3, vertex % 3
}

// CertificateToAssignment is the certificate map g for the reduction
// ThreeSATToIndependentSet: given the original 3-SAT instance and an
// independent set of size len(phi.Clauses) in f(phi), it reconstructs a
// satisfying Assignment for phi. This is the constructive half of the
// reduction's correctness discussed in docs/05-cook-levin-and-practice.md,
// §5.4 — the "x ∈ A ⟹ f(x) ∈ B, with certificate" direction made
// executable, matching the certificate-is-the-object-you-wanted shape of
// NP certificates described in docs/01-decision-problems-p-np.md, §1.3.
//
// The decoding relies on exactly the structure ThreeSATToIndependentSet
// builds: a valid independent set of that size must contain exactly one
// vertex per clause-triangle (no triangle can contribute two adjacent
// vertices), so every chosen vertex names a distinct clause and a literal
// within it; setting that literal to true is safe because the conflict
// edges rule out two chosen vertices ever naming complementary literals
// for the same variable. Variables that no chosen vertex happens to name
// are left false — they are unconstrained by phi's satisfaction and
// cannot break it.
//
// CertificateToAssignment does not itself check that cert is a valid
// certificate for f(phi); callers are expected to verify its output
// against phi.Verify, exactly as the pipeline's fourth check does (see
// docs/06-pipeline-architecture.md, §6.3). Handed a cert that is not
// actually an independent set of the right size, it may silently produce
// an assignment that does not satisfy phi.
func CertificateToAssignment(phi ThreeSAT, cert VertexSet) Assignment {
	assign := make(Assignment, phi.NumVars+1)
	for _, v := range cert {
		clause, position := clauseAndPosition(v)
		lit := phi.Clauses[clause][position]
		if lit > 0 {
			assign[lit] = true
		} else {
			assign[-lit] = false
		}
	}
	return assign
}
