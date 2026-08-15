package karp

// Verifier is the generic shape every problem's per-type verifier has:
// given a candidate certificate of type C, decide whether it is valid for
// this instance. Every concrete instance type in this package — ThreeSAT,
// IndependentSet, and whatever is added after them — satisfies
// Verifier[C] for its own certificate type simply by having a
// Verify(C) bool method; nothing needs to declare that it implements this
// interface, Go's structural typing does that for free.
type Verifier[C any] interface {
	Verify(C) bool
}

// IndexSet is the constraint on certificate types BruteForceOracle can
// enumerate: any named slice-of-int type, such as VertexSet. Every
// certificate in this project other than ThreeSAT's own Assignment is, at
// bottom, a subset of some finite universe (chosen vertices, chosen
// items, ...) named directly by their indices — see docs/01-decision-
// problems-p-np.md, §1.3, on why the certificate for a set-existence
// problem is "the set you were looking for", not some encoded proxy for
// it.
type IndexSet interface {
	~[]int
}

// maxBruteForceUniverse caps the size of the universe BruteForceOracle
// will enumerate the power set of. 2^30 candidate subsets is already far
// beyond what any well-behaved instance in this pipeline should reach —
// per docs/06-pipeline-architecture.md, §6.1, every instance other than
// ThreeSAT is born as the image of a reduction and stays as small as the
// ThreeSAT instance that produced it. Hitting this cap means a generator
// somewhere is handing the oracle an instance too large for it to remain
// the small-but-obviously-correct reference it is meant to be — not that
// BruteForceOracle itself has a bug worth chasing.
const maxBruteForceUniverse = 30

// BruteForceOracle decides an instance by enumerating every possible
// certificate over a universe of universeSize indices (0..universeSize-1)
// and calling inst.Verify on each one, returning the first one that
// succeeds. It is not an efficient algorithm — it is, deliberately, the
// least efficient one that could possibly work: the "try every y with
// |y| <= p(|x|)" enumeration already implicit in the certificate
// definition of NP (docs/01-decision-problems-p-np.md, §1.3), made
// literal instead of left as an existence argument. Its correctness
// therefore never depends on any problem-specific insight, only on
// inst.Verify being correct — which is exactly why it can serve as an
// independent, structurally trustworthy reference oracle for every
// problem other than ThreeSAT (docs/06-pipeline-architecture.md, §6.2):
// unlike ThreeSAT, whose instances are generated directly and need a real
// solver, every other instance in this pipeline is born as the image of a
// reduction and stays small enough for "correct but slow" to be an
// acceptable trade, not "correct but fast".
//
// BruteForceOracle panics if universeSize exceeds maxBruteForceUniverse,
// since silently iterating an astronomically large power set would look
// like a hang rather than a bug to whoever is waiting on it.
func BruteForceOracle[C IndexSet](inst Verifier[C], universeSize int) (satisfiable bool, cert C) {
	if universeSize > maxBruteForceUniverse {
		panic("karp: BruteForceOracle: universe too large to enumerate exhaustively")
	}
	for mask := 0; mask < (1 << universeSize); mask++ {
		var candidate []int
		for i := 0; i < universeSize; i++ {
			if mask&(1<<i) != 0 {
				candidate = append(candidate, i)
			}
		}
		c := C(candidate)
		if inst.Verify(c) {
			return true, c
		}
	}
	var zero C
	return false, zero
}

// SolveIndependentSet decides inst by brute force, enumerating candidate
// vertex sets over inst.G's vertices. It is the reference oracle for
// IndependentSet described in docs/06-pipeline-architecture.md, §6.1 —
// the fixed, independent point of comparison TestThreeSATToIndependentSetProperty
// (property_test.go) checks ThreeSATToIndependentSet against, alongside
// the real SAT oracle deciding the original ThreeSAT instance.
func SolveIndependentSet(inst IndependentSet) (bool, VertexSet) {
	return BruteForceOracle[VertexSet](inst, inst.G.NumVertices)
}
