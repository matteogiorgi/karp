package karp

// Graph is a simple undirected graph on vertices 0..NumVertices-1. Edges
// are unordered pairs: {u, v} and {v, u} name the same edge, and nothing
// in this package relies on which of the two orderings a given edge
// happens to be stored in.
type Graph struct {
	NumVertices int
	Edges       [][2]int
}

// IndependentSet is an instance of the Independent Set decision problem:
// does G contain a set of at least K pairwise non-adjacent vertices? It is
// the first link of this project's reduction chain after 3-SAT itself
// (docs/05-cook-levin-and-practice.md, §5.3): ThreeSATToIndependentSet, in
// threesat_to_independent_set.go, is the concrete witness of 3-SAT ≤p
// IndependentSet in the sense of the many-one reduction defined in
// docs/03-polynomial-many-one-reductions.md, §3.1.
type IndependentSet struct {
	G Graph
	K int
}

// VertexSet is a certificate for an IndependentSet instance: the chosen
// vertices themselves. This mirrors the certificate examples worked
// through in docs/01-decision-problems-p-np.md, §1.3 — for a problem that
// asks "does a set with property P exist?", the certificate is, more or
// less, the set you were looking for, not some encoded proxy for it. The
// slice may contain duplicates or appear in any order; Verify normalizes
// it before checking anything else.
type VertexSet []int

// Verify reports whether cert names an independent set of size >= inst.K
// in inst.G: every chosen vertex must be a real vertex of inst.G (0 <= v <
// inst.G.NumVertices — a certificate naming something outside V is not a
// subset of V at all, and so cannot be an independent set of it, no matter
// how many such names it lists), no two chosen vertices may be joined by
// an edge, and at least K distinct vertices must have been chosen
// (duplicates in cert do not count twice). All of this is checked in time
// polynomial in the size of cert and inst.G, which is exactly the verifier
// side of the NP certificate definition in docs/01-decision-problems-p-
// np.md, §1.3 — Verify does not search for an independent set, it only
// checks one that is handed to it.
func (inst IndependentSet) Verify(cert VertexSet) bool {
	chosen := make(map[int]bool, len(cert))
	for _, v := range cert {
		if v < 0 || v >= inst.G.NumVertices {
			return false
		}
		chosen[v] = true
	}
	if len(chosen) < inst.K {
		return false
	}
	for _, e := range inst.G.Edges {
		if chosen[e[0]] && chosen[e[1]] {
			return false
		}
	}
	return true
}
