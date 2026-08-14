# 4. What Reductions Are For

## 4.1 Transferring difficulty: the "easy" direction

If `A ≤p B` via `f`, and `B ∈ P`, then `A ∈ P`. The proof is the direct composition of the two algorithms: to decide `x ∈ A`, compute `y = f(x)` (polynomial time by definition of reduction), then run on `y` the polynomial algorithm for `B` (polynomial time in `|y|`).

The point that makes the argument non-trivial is the composition of sizes: `f` runs in polynomial time, so its *output* `y` has length polynomial in `|x|` (it cannot write more characters than the number of steps it takes). Hence the second algorithm, polynomial in `|y|`, is in turn polynomial in `|x|` — the composition of two polynomials is a polynomial. This is precisely why the definition of reduction ([Section 3.1](03-polynomial-many-one-reductions.md#31-definition)) requires `f` to be computable in polynomial time, and not merely computable: an `f` that is computable but exponential would still transfer the relation `x ∈ A ⟺ f(x) ∈ B`, but it would break the efficiency chain — and with it, the only reason the reduction is interesting.

## 4.2 Transferring difficulty: the "hard" direction (by contraposition)

The practical content lies in the contrapositive of the previous statement: if `A ≤p B` and `A ∉ P`, then `B ∉ P`. (If `B ∈ P` were true, then by 4.1 `A ∈ P` would also hold — a contradiction.)

This is the direction that is actually used. To show that a newly encountered problem `B` is (presumably) hard, one does **not** try to reduce `B` to something else — one exhibits a reduction **from** an already-known hard problem `A` **to** `B` (`A ≤p B`). This is the reason — often a source of confusion the first time the argument is encountered — why every reduction in this project runs "forward" starting from 3-SAT: 3-SAT is the known-hard reference problem ([Section 5](05-cook-levin-and-practice.md)), and each reduction `3-SAT ≤p X` is a way of saying "`X` inherits the difficulty of 3-SAT", not the other way around.

## 4.3 The (pre)order on difficulty

`≤p` is **reflexive** (the identity function is a trivial reduction of `A` to itself) and **transitive**: if `A ≤p B` via `f` and `B ≤p C` via `f'`, then `A ≤p C` via the composition `f' ∘ f` — which is still polynomial, by the same size-composition argument as in 4.1. These two properties make `≤p` a **preorder** on the set of decision problems: not a total order (not all problems are comparable), and not a strict partial order either (`A ≤p B` and `B ≤p A` do not imply `A = B`, only that the two problems have the same difficulty "up to polynomial reduction").

A note on where the interest of this relation actually lies: **all** non-trivial problems in `P` (i.e., different from the empty language and the universal one — two degenerate cases that must be excluded for a technical reason tied to how the constant reduction is built) reduce to one another in polynomial time — so `≤p`, restricted to `P`, collapses into a single indistinguishable class. `≤p` becomes a discriminating tool only when looking at problems suspected of lying *outside* `P`: that is where the order stops being trivial and starts stratifying problems by relative difficulty.

## 4.4 Completeness

A problem `B` is **NP-complete** if both of the following hold:

1. `B ∈ NP`;
2. for **every** `A ∈ NP`, `A ≤p B`.

An NP-complete problem is therefore, by construction, the hardest possible problem inside NP with respect to `≤p`: everything else in the class reduces to it. Combining this with 4.2: if even a single NP-complete problem turned out to be in `P`, then *every* problem in NP would be too (just compose its reduction to `B` with the polynomial algorithm for `B`) — i.e., `P = NP`. This is the real weight of the statement "`B` is NP-complete": it is not just "it's hard", it is "it's exactly as hard as the hardest problem in the entire class NP, whatever that turns out to be".

Taken literally, the definition is unusable as a direct proof tool: showing `A ≤p B` for *every* `A ∈ NP` — an infinite set of problems, in general not even explicitly enumerable — is not something that can be done by hand, problem by problem. The way out is the transitivity from 4.3: if a single problem `A₀` is already known to be NP-complete, and `A₀ ≤p B` is exhibited, then for every `A ∈ NP` we have `A ≤p A₀ ≤p B` by transitivity — i.e., the single reduction `A₀ ≤p B` suffices to inherit *all* the infinitely many reductions the definition requires.

This is exactly why this project can limit itself to building a **chain** of reductions (3-SAT → Independent Set → Vertex Cover → …) instead of starting over from the definition every time: each new link in the chain only requires a single reduction from the previous link, and NP-completeness propagates by transitivity along the whole chain. A single starting point remains to be justified outside this mechanism — the first link, the one that makes 3-SAT itself NP-complete without being able to rely on any previous problem. That is the content of the Cook–Levin theorem, the subject of the next section.

---

[← Previous: 3. Polynomial Many-One Reductions](03-polynomial-many-one-reductions.md) · [README](../README.md) · [Next: 5. Cook–Levin and Practice →](05-cook-levin-and-practice.md)
