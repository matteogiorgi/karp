// Package karp implements an explicit, typed pipeline of Karp reductions
// between NP-complete decision problems — "reduction as compilation": every
// reduction is a pure function between two problem-specific instance types,
// paired with an inverse map that carries a certificate of the reduced
// problem back to a certificate of the original one.
//
// The package does not stand on its own: it is the executable half of a
// theoretical report that lives in docs/ at the root of the repository.
// Every exported type and function below corresponds to a specific notion
// introduced there, and its doc comment says which — reading the two side
// by side is the intended way to use this codebase, not just the code
// alone. The correspondence is summarized in docs/06-pipeline-architecture.md,
// §6.1.
package karp
