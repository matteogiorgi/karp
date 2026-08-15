package karp

// Literal is a single occurrence of a boolean variable inside a clause: a
// positive value v denotes the variable v itself, a negative value -v
// denotes its negation ¬v. Zero is never a valid literal, since variables
// are indexed starting at 1 — this leaves the sign of the int doing double
// duty as the literal's polarity, which is why every place that reads a
// Literal has to recover both the variable index and the polarity from it
// (see clauseSatisfied below).
type Literal int

// Clause3 is a 3-SAT clause: a disjunction of exactly three literals,
// (l1 ∨ l2 ∨ l3). Fixing the array's length to 3, instead of using a slice,
// makes "a clause has exactly three literals" part of the type itself — a
// clause of the wrong width simply cannot be constructed, so the invariant
// is enforced by the compiler rather than by a runtime check written and
// maintained by hand. This is the concrete payoff of the "explicit over
// magic" design discussed for this project's type choices in
// docs/01-decision-problems-p-np.md, §1.1.
//
// Duplicate literals inside a clause are allowed and meaningful, not a
// malformed input: they arise naturally when a wider SAT clause is
// normalized down to width 3 by padding a literal (docs/05-cook-levin-and-
// practice.md, §5.2), and Verify treats a repeated literal exactly like any
// other one — no special-casing is required, or present, below.
type Clause3 [3]Literal

// ThreeSAT is an instance of the 3-SAT decision problem: NumVars boolean
// variables, indexed 1..NumVars, and a conjunction of Clauses over them. A
// ThreeSAT value plays the role of an instance x of the language L = 3-SAT
// in the sense of docs/01-decision-problems-p-np.md, §1.1 — "solving" it
// means answering whether some assignment of the variables satisfies every
// clause, which is exactly what Verify checks for one candidate assignment
// at a time.
//
// 3-SAT is the root of this project's reduction chain: its NP-completeness
// is cited (the Cook–Levin theorem, docs/05-cook-levin-and-practice.md,
// §5.1) and assumed here, not re-proved in code — see the explicit scope
// note in §5.3 of that same section.
type ThreeSAT struct {
	NumVars int
	Clauses []Clause3
}

// Assignment is a certificate for a ThreeSAT instance: a truth value per
// variable, indexed 1..NumVars. Index 0 is always false and never read
// meaningfully; it is kept only so a literal's variable index can be used
// directly as a slice index, without an off-by-one adjustment scattered
// across every call site that touches an Assignment.
type Assignment []bool

// Verify reports whether assign satisfies every clause of phi — in other
// words, whether assign is a valid certificate for phi under the
// verifier-and-certificate definition of NP in docs/01-decision-problems-
// p-np.md, §1.3. Verify only ever checks a candidate assignment that is
// already in hand; it never searches for one, and it is not supposed to —
// the definition of NP only requires the checking side to be efficient,
// and this method is exactly that check, one clause at a time.
func (phi ThreeSAT) Verify(assign Assignment) bool {
	for _, clause := range phi.Clauses {
		if !clauseSatisfied(clause, assign) {
			return false
		}
	}
	return true
}

// clauseSatisfied reports whether at least one literal in clause evaluates
// to true under assign. A clause is a disjunction, so a single satisfied
// literal is enough to satisfy the whole clause; the loop below stops at
// the first one it finds rather than examining the remaining two.
func clauseSatisfied(clause Clause3, assign Assignment) bool {
	for _, lit := range clause {
		v := int(lit)
		if v < 0 {
			v = -v
		}
		val := assign[v]
		if lit < 0 {
			val = !val
		}
		if val {
			return true
		}
	}
	return false
}
