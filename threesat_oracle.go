package karp

import (
	"fmt"
	"strings"

	"github.com/crillab/gophersat/solver"
)

// DIMACS serializes phi to the DIMACS CNF text format: a header line
// "p cnf <NumVars> <NumClauses>" followed by one line per clause, each a
// space-separated list of signed literals terminated by 0. This is the
// only representation of a ThreeSAT instance that ever leaves this
// package's own types — see docs/06-pipeline-architecture.md, §6.4: DIMACS
// is the literal boundary between the project's typed world and the
// solver, not just a narrative one, and this method is where that
// boundary is actually crossed.
func (phi ThreeSAT) DIMACS() string {
	var b strings.Builder
	fmt.Fprintf(&b, "p cnf %d %d\n", phi.NumVars, len(phi.Clauses))
	for _, c := range phi.Clauses {
		fmt.Fprintf(&b, "%d %d %d 0\n", c[0], c[1], c[2])
	}
	return b.String()
}

// SolveThreeSAT is the SAT oracle described in docs/06-pipeline-
// architecture.md, §6.1: the one point in this pipeline where a real,
// efficient solver is needed, because ThreeSAT is the only instance type
// generated directly rather than born as the image of a reduction (§6.1
// again). It hands gophersat the DIMACS text produced by phi.DIMACS(),
// in-process and in memory — no subprocess, no temporary file — which is
// exactly what §6.5 means by "in-process": the same binary, talking to
// gophersat through its DIMACS-reading entry point rather than its native
// Problem struct, so that the DIMACS boundary of §6.4 stays literal.
//
// If phi is satisfiable, SolveThreeSAT returns true together with a
// satisfying Assignment recovered from gophersat's model — gophersat's own
// model is a 0-indexed []bool naming variable i+1 at index i, which
// SolveThreeSAT re-indexes into this package's 1-indexed Assignment
// convention (docs/01-decision-problems-p-np.md, §1.3, and the doc comment
// on Assignment). If phi is unsatisfiable, it returns false and a nil
// Assignment.
//
// A malformed DIMACS string from phi.DIMACS(), or a solver status other
// than Sat or Unsat for what is a plain CNF problem with no external stop
// signal, are both programmer errors rather than expected outcomes, and
// SolveThreeSAT panics on either rather than silently returning a
// misleading answer.
func SolveThreeSAT(phi ThreeSAT) (bool, Assignment) {
	pb, err := solver.ParseCNF(strings.NewReader(phi.DIMACS()))
	if err != nil {
		panic(fmt.Sprintf("karp: SolveThreeSAT: gophersat rejected the generated DIMACS: %v", err))
	}

	s := solver.New(pb)
	switch status := s.Solve(); status {
	case solver.Sat:
		model := s.Model()
		assign := make(Assignment, phi.NumVars+1)
		for i, v := range model {
			assign[i+1] = v
		}
		return true, assign
	case solver.Unsat:
		return false, nil
	default:
		panic(fmt.Sprintf("karp: SolveThreeSAT: gophersat returned unexpected status %v for a plain CNF problem", status))
	}
}
