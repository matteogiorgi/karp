package karp

import "testing"

// TestThreeSATVerify checks Verify against a small, hand-worked formula
// where the satisfying and unsatisfying assignments are easy to confirm by
// inspection, rather than relying on the very code under test to tell us
// what the "right" answer is.
func TestThreeSATVerify(t *testing.T) {
	// (x1 ∨ x2 ∨ x3) ∧ (¬x1 ∨ ¬x2 ∨ x3)
	phi := ThreeSAT{
		NumVars: 3,
		Clauses: []Clause3{
			{1, 2, 3},
			{-1, -2, 3},
		},
	}

	cases := []struct {
		name   string
		assign Assignment
		want   bool
	}{
		{"all true satisfies both clauses", Assignment{false, true, true, true}, true},
		{"x1=x2=true forces x3 for second clause", Assignment{false, true, true, false}, false},
		{"all false fails first clause", Assignment{false, false, false, false}, false},
		{"x1 false, x3 true satisfies both", Assignment{false, false, false, true}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := phi.Verify(c.assign); got != c.want {
				t.Errorf("Verify(%v) = %v, want %v", c.assign, got, c.want)
			}
		})
	}
}

// TestThreeSATVerifyDuplicateLiteral guards the specific case the report
// flags in docs/05-cook-levin-and-practice.md, §5.2: a clause padded by
// repeating a literal, as the SAT -> 3-SAT normalization would produce for
// a clause shorter than three literals. Clause3 does not forbid this
// shape, and Verify must not treat it as anything other than an ordinary
// (if slightly redundant) disjunction — a regression here would silently
// break every reduction downstream of a padded clause.
func TestThreeSATVerifyDuplicateLiteral(t *testing.T) {
	// (x1 ∨ x1 ∨ x2)
	phi := ThreeSAT{
		NumVars: 2,
		Clauses: []Clause3{{1, 1, 2}},
	}

	if !phi.Verify(Assignment{false, true, false}) {
		t.Error("expected duplicate-literal clause to be satisfied when x1 is true")
	}
	if phi.Verify(Assignment{false, false, false}) {
		t.Error("expected duplicate-literal clause to fail when both x1 and x2 are false")
	}
}
