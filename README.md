# Karp All-Reductions Pipeline

An explicit, typed pipeline of reductions between NP-complete problems, verified end-to-end against a real SAT solver — "reduction as compilation". The name is Richard Karp's, the author of the polynomial many-one reductions this project implements and of the 21 NP-complete problems that they connect.

## What it is

Every reduction between two NP-complete problems (e.g. 3-SAT → Independent Set) is implemented as a pure function `f: InstanceA → InstanceB`, paired with an inverse map on certificates `g: CertificateB → CertificateA`. Correctness is not just claimed: it is put to the test with property-based testing, comparing on randomly generated instances the answer of a real SAT solver (for 3-SAT) against a reference oracle for the reduced problem, and checking that `g` always reconstructs a valid certificate.

Chain implemented (deliberately small, to stay a self-contained project):

```
3-SAT → Independent Set → Vertex Cover (→ Clique by complementation)
3-SAT → Subset Sum
```

3-SAT is taken as the root of the chain — its NP-completeness is cited (the Cook–Levin theorem) and assumed, not re-proved in code.

**Language: Go.** Chosen over R after a direct comparison on three axes (clarity of the reduction code, ergonomics of property testing, maturity of the SAT solver available in each language) — the last one was decisive, since the only in-process SAT binding for R (`rpicosat`) has been archived from CRAN since 2022.

## Project status

- [x] Theoretical report (sections 1–6)
- [ ] Go implementation of the reductions
- [ ] Property tests against the SAT oracle

## Theoretical report

1. [Decision Problems, P and NP](docs/01-decision-problems-p-np.md) — languages, the class P, NP via verifier+certificate, the equivalence with nondeterministic Turing machines, and the NP/co-NP asymmetry.
2. [What Complexity Classes Are For](docs/02-purpose-of-complexity-classes.md) — they classify the problem not the algorithm, they are robust with respect to the model of computation, they give a common vocabulary across domains, they transfer negative results.
3. [Polynomial Many-One Reductions](docs/03-polynomial-many-one-reductions.md) — Karp's definition, and why the project uses only that (and not the more general Turing/Cook reductions).
4. [What Reductions Are For](docs/04-purpose-of-reductions.md) — transferring difficulty in both directions, the order induced by `≤p`, and the notion of completeness.
5. [Cook–Levin and Practice](docs/05-cook-levin-and-practice.md) — the Cook–Levin theorem as the root, the SAT → 3-SAT normalization, and the two-step recipe used by every subsequent reduction.
6. [Pipeline Architecture](docs/06-pipeline-architecture.md) — the code's components (typed instances, reductions, oracles, the DIMACS boundary) and their one-to-one correspondence with the sections above.
