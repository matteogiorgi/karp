# Karp All-Reductions Pipeline

An explicit, typed pipeline of reductions between NP-complete problems, verified end-to-end against a real SAT solver — "reduction as compilation". The name is Richard Karp's, the author of the polynomial many-one reductions this project implements and of the 21 NP-complete problems that they connect.

## What it is

Every reduction between two NP-complete problems (e.g. 3-SAT → Independent Set) is implemented as a pure function `f: InstanceA → InstanceB`, paired with an inverse map on certificates `g: CertificateB → CertificateA`. Correctness is not just claimed: it is put to the test with property-based testing, comparing on randomly generated instances the answer of a real SAT solver (for 3-SAT) against a reference oracle for the reduced problem, and checking that `g` always reconstructs a valid certificate.

The reduction chain (deliberately small, to stay a self-contained project — see [Project status](#project-status) below for how much of it is actually built so far):

```
3-SAT → Independent Set → Vertex Cover      (set complement: same graph, S ↔ V∖S)
                        └→ Clique            (graph complement: same S, in Ḡ)
3-SAT → Subset Sum
```

Vertex Cover and Clique are both one step from Independent Set, not from each other: Vertex Cover complements the *vertex set* (`k' = |V| − k`); Clique complements the *graph* and keeps the same set and `k`. The two `g` certificate maps are correspondingly different (set-complement vs. identity).

3-SAT is taken as the root of the chain — its NP-completeness is cited (the Cook–Levin theorem) and assumed, not re-proved in code.

**Language: Go.** Chosen over R after a direct comparison on three axes: static types make each reduction's structural invariants (e.g. a clause has exactly 3 literals) checked by the compiler instead of at runtime; `rapid`'s property-testing generators are plain imperative Go code with automatic shrinking, versus the more combinator-heavy style the same thing needs in R's `hedgehog`; and the only in-process SAT binding available for R, `rpicosat`, has been archived from CRAN since 2022 — installable only from source, a maintenance risk this project prefers not to carry rather than a hard blocker on its own.

## Project status

- [x] Theoretical report (sections 1–6)
- [x] Types and verifiers for 3-SAT and Independent Set
- [x] Reduction 3-SAT ≤p Independent Set (`f` and its certificate map `g`)
- [x] Brute-force reference oracle
- [x] Real SAT oracle (`gophersat`, via DIMACS)
- [x] Property-based test for 3-SAT → Independent Set
- [ ] Reduction Independent Set ≤p Vertex Cover
- [ ] Reduction Independent Set ≤p Clique
- [ ] Reduction 3-SAT ≤p Subset Sum
- [ ] Property tests for the three reductions above

## Project layout

| File | What it is |
|---|---|
| [`go.mod`](go.mod), [`go.sum`](go.sum) | Go module definition (`github.com/matteogiorgi/karp`) and its two dependencies, `gophersat` and `rapid`. |
| [`doc.go`](doc.go) | Package-level doc comment: what the package is, and how to read it alongside `docs/`. |
| [`threesat.go`](threesat.go) | Types and verifier for 3-SAT — `Literal`, `Clause3`, `ThreeSAT`, `Assignment` — the root problem of the reduction chain. |
| [`threesat_test.go`](threesat_test.go) | Table-driven tests for `ThreeSAT.Verify`, including the duplicate-literal case a padded clause produces. |
| [`independent_set.go`](independent_set.go) | Types and verifier for Independent Set — `Graph`, `IndependentSet`, `VertexSet` — the first problem reduced from 3-SAT. |
| [`independent_set_test.go`](independent_set_test.go) | Table-driven tests for `IndependentSet.Verify`. |
| [`threesat_to_independent_set.go`](threesat_to_independent_set.go) | The reduction `f = ThreeSATToIndependentSet` (3-SAT ≤p Independent Set) and its certificate map `g = CertificateToAssignment`. |
| [`threesat_to_independent_set_test.go`](threesat_to_independent_set_test.go) | Structural test of the graph `f` builds, plus a hand-worked round trip through `f`, `g`, and both `Verify` methods. |
| [`oracle.go`](oracle.go) | The generic brute-force reference oracle `BruteForceOracle` (enumerate every candidate certificate, call `Verify`), and `SolveIndependentSet`, its instantiation for `IndependentSet`. |
| [`oracle_test.go`](oracle_test.go) | Tests for both the "yes" and "no" sides of the oracle, the universe-size safety cap, and its agreement with `ThreeSATToIndependentSet` end to end. |
| [`threesat_oracle.go`](threesat_oracle.go) | `DIMACS`, the literal boundary to the solver, and `SolveThreeSAT`, the real SAT oracle — `gophersat` called in-process via its DIMACS-reading entry point. |
| [`threesat_oracle_test.go`](threesat_oracle_test.go) | The exact DIMACS text for a hand-picked formula, the "yes" and "no" sides of the real oracle, and end-to-end agreement with the reduction and the brute-force oracle. |
| [`property_test.go`](property_test.go) | `genThreeSAT`, the one `rapid` generator this project needs, and the property test running the four-check pipeline of docs/06-pipeline-architecture.md, §6.3 on hundreds of random instances. |
| [`docs/`](docs/) | The theoretical report — see [Theoretical report](#theoretical-report) below for the section-by-section index. |
| [`DEVLOG.md`](DEVLOG.md) | A running record of concrete bugs and design gaps found while building this project, and what was done about them — both in the GitHub Pages rendering pipeline and in the Go implementation. |

The files below aren't part of the project itself — they only exist to make this repository browsable as a GitHub Pages site.

| File | What it is |
|---|---|
| `_layouts/default.html` | The site's page layout (the default GitHub Pages theme's own template), which is where `_includes/head-custom.html` gets pulled in. Not linked — Jekyll never publishes `_layouts/` as a browsable file, on GitHub Pages or in a local build. |
| `_includes/head-custom.html` | MathJax and Mermaid rendering, plus the footer-nav layout CSS, injected into every GitHub Pages build. Not linked, for the same reason as `_layouts/` above. |
| [`_config.yml`](_config.yml) | Forces Jekyll to publish `Gemfile`, `Gemfile.lock`, `.gitignore`, and itself (all excluded by default), so the links to them in this table don't 404 on the live site. |
| [`Gemfile`](Gemfile), [`Gemfile.lock`](Gemfile.lock) | Pin the `github-pages` gem version so a local Jekyll build matches GitHub Pages exactly. |
| [`.gitignore`](.gitignore) | Excludes build and tooling artifacts from version control — Jekyll's (`_site/`, `.bundle/`) and `rapid`'s (`testdata/rapid/`, its failure-reproduction files). |

## Theoretical report

1. [Decision Problems, P and NP](docs/01-decision-problems-p-np.md) — languages, the class P, NP via verifier+certificate, the equivalence with nondeterministic Turing machines, and the NP/co-NP asymmetry.
2. [What Complexity Classes Are For](docs/02-purpose-of-complexity-classes.md) — they classify the problem not the algorithm, they are robust with respect to the model of computation, they give a common vocabulary across domains, they transfer negative results.
3. [Polynomial Many-One Reductions](docs/03-polynomial-many-one-reductions.md) — Karp's definition, and why the project uses only that (and not the more general Turing/Cook reductions).
4. [What Reductions Are For](docs/04-purpose-of-reductions.md) — transferring difficulty in both directions, the order induced by `≤p`, and the notion of completeness.
5. [Cook–Levin and Practice](docs/05-cook-levin-and-practice.md) — the Cook–Levin theorem as the root, the SAT → 3-SAT normalization, and the two-step recipe used by every subsequent reduction.
6. [Pipeline Architecture](docs/06-pipeline-architecture.md) — the code's components (typed instances, reductions, oracles, the DIMACS boundary) and their one-to-one correspondence with the sections above.
