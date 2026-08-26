# 5. How Reductions Are Used in Practice: Cook–Levin as the Root, Then All by Reduction

## 5.1 The Cook–Levin theorem

[Section 4.4](04-purpose-of-reductions.md#44-completeness) left an open problem: the definition of NP-completeness requires a reduction from *every* problem in NP, a requirement no proof can satisfy directly problem by problem — and transitivity only makes it practicable if a known NP-complete problem **already exists** to start from. The Cook–Levin theorem (1971; independently, Levin around the same time) is exactly that starting point: it proves that **SAT** (satisfiability of Boolean formulas, not necessarily given in conjunctive normal form) is NP-complete, and it is the only result in this whole story proved *from the raw definition*, without relying on any earlier reduction.

$\text{SAT} \in NP$ is the easy part ([Section 1.3](01-decision-problems-p-np.md#13-the-class-np-definition-via-verifier-and-certificate): certificate = assignment, verification = evaluate the formula, linear time). The substantial part is showing that **every** $L \in NP$ satisfies $L \le_p \text{SAT}$. The idea, at a level that makes its structure visible without developing the construction all the way (it is a project of its own, deliberately out of scope here — see the scope note in 5.3):

- By definition ([Section 1.3](01-decision-problems-p-np.md#13-the-class-np-definition-via-verifier-and-certificate)), $L \in NP$ means there is a polynomial verifier $V$ and a polynomial $p$ such that $x \in L \iff \exists\, y,\ \|y\| \le p(\|x\|),\ V(x,y)$ accepts. The verifier $V$ is, at bottom, a program — formalizable as a Turing machine that runs in a known polynomial time $q(\|x\|)$.
- A Boolean formula $\varphi_x$ is built that describes the entire execution of $V$ on input $x$ and on a certificate $y$ **not yet fixed** (the Boolean variables of $\varphi_x$ encode, cell by cell and time step by time step, what is written on $V$'s tape during the computation — the so-called *computation tableau*, a time × space grid of polynomial size). The clauses of $\varphi_x$ enforce three things: the initial configuration correctly encodes $x$ (the cells for $y$ are left free — they are the formula's degrees of freedom), every transition respects $V$'s rules (**local** constraints, since a cell's value at step $t+1$ depends only on a restricted neighborhood at step $t$), and an accepting state is reached.
- The key point: $\varphi_x$ is satisfiable if and only if there is a way to fill in the free cells (i.e., a certificate $y$) that makes $V$ accept — that is, if and only if $x \in L$. And building $\varphi_x$ from $x$ takes polynomial time, because the tableau has polynomial dimensions ($q(\|x\|)$ rows) and each individual clause is local, hence of constant size.

The result — $f(x) = \varphi_x$ — is, in every respect, a polynomial many-one reduction in the sense of [Section 3](03-polynomial-many-one-reductions.md), with $A = L$ and $B = \text{SAT}$. The difference from every other reduction in this project is **what** is being reduced: not one concrete combinatorial problem to another, but the *entire notion of verifiable computation* (any $L \in NP$, via its generic verifier $V$) to a single syntactic structure (a propositional formula). It is the only bridge between the definitional world (verifiers, Turing machines) and the combinatorial world (clauses, graphs, sets) in which every subsequent reduction lives — including the one used as this project's foundation.




## 5.2 From SAT to 3-SAT

Cook–Levin, in its classical form, produces general SAT — formulas with clauses of arbitrary width. The root problem chosen for this project is **3-SAT** (clauses of width exactly 3), which requires one additional step — but of a completely different nature from the previous one: $\text{SAT} \le_p \text{3-SAT}$ is proved with a purely **syntactic** transformation, touching no computational semantics at all. Every wide clause is split into a chain of 3-literal clauses by introducing auxiliary variables ($(l_1 \lor l_2 \lor l_3 \lor l_4)$ becomes $(l_1 \lor l_2 \lor z) \land (\lnot z \lor l_3 \lor l_4)$, and so on for longer clauses), and clauses that are too short are padded by duplicating a literal. The transformation preserves satisfiability clause by clause and is evidently polynomial.

It is worth noting the contrast explicitly: $\text{SAT} \le_p \text{3-SAT}$ is mechanical, almost administrative — it normalizes one syntactic form into another. Cook–Levin ($L \le_p \text{SAT}$ for every $L \in NP$) is the only reduction in the entire chain that *creates* difficulty out of nothing, in the sense that it compiles an entire model of computation into combinatorial structure. Every subsequent reduction in this project ([Section 6](06-pipeline-architecture.md)) resembles the second category more than the first in spirit — it moves difficulty that already exists from one combinatorial structure to another, never having to "invent" difficulty from scratch — but this still requires domain-specific structural insight, unlike the purely syntactic transformation just described.




## 5.3 The practical recipe — and what the project assumes instead of proving

Taken together, Cook–Levin and $\text{SAT} \le_p \text{3-SAT}$ fix a single anchoring fact: **3-SAT is NP-complete**. From this point on, every new NP-completeness proof always follows the same two-step recipe (it mirrors exactly the definition in [4.4](04-purpose-of-reductions.md#44-completeness)):

1. show that the new problem $B$ is in NP (exhibit a verifier/certificate — almost always the easy step, [Section 1.3](01-decision-problems-p-np.md#13-the-class-np-definition-via-verifier-and-certificate));
2. exhibit **a single** reduction $A_0 \le_p B$ from a problem $A_0$ **already** known to be NP-complete (never from the raw definition) — by transitivity ([Section 4.3](04-purpose-of-reductions.md#43-the-preorder-on-difficulty)), this suffices.

In practice, nobody after Cook and Levin has ever again had to repeat step 2 in its general form ("reduce from every $L \in NP$"): one picks as $A_0$ the already-known problem whose structure most resembles that of the new target — which is why recognizable "families" of reductions exist (graph problems reducing to other graph problems, numeric problems encoding the bits of a formula, and so on), rather than an arbitrary, disconnected collection.

**Explicit scope note**: this project **cites** the Cook–Levin theorem as an established fact (Section 5.1 describes its structure, without implementing it) and **assumes** that 3-SAT is NP-complete as a starting point — it does not prove this constructively in code. Building the Turing-machine-to-CNF compiler described in 5.1 is a separate project (a possible future extension, not a gap in this one). What this project *does* concretely is everything that comes after that anchoring: the chain $\text{3-SAT} \le_p \text{Independent Set}$, followed by two separate one-step reductions **out of** Independent Set — $\text{Independent Set} \le_p \text{Vertex Cover}$ (complementing the *vertex set* within the same graph, $S \leftrightarrow V \setminus S$) and $\text{Independent Set} \le_p \text{Clique}$ (complementing the *graph* itself and keeping the same set) — plus the branch $\text{3-SAT} \le_p \text{Subset Sum}$, each verified end-to-end against a real SAT solver used as the ground-truth oracle for 3-SAT instances. Vertex Cover and Clique are siblings one step from Independent Set, not a further link in a single line — a distinction worth keeping straight because their certificate maps $g$ are different (set-complement vs. identity).




## 5.4 Every reduction as a constructive mini-proof

With the anchor to 3-SAT established, every reduction function $f$ written in this project's code is, literally, the concrete witness of one step of the recipe in 5.3: not an isolated exercise, but a link that inherits — by transitivity — the NP-completeness of the entire chain built up to that point.

The distinction already anticipated in [1.6](01-decision-problems-p-np.md#16-an-asymmetry-worth-keeping-in-mind) must be kept in mind, between what the mathematics guarantees and what the code verifies empirically:

- the constructive direction ($x \in A \implies f(x) \in B$, with certificate) is the one the project makes fully executable: the function $g$ produces a concrete certificate for $B$, and that certificate is verified directly ([Section 1.3](01-decision-problems-p-np.md#13-the-class-np-definition-via-verifier-and-certificate));
- the direction "no solution is ever lost" ($x \notin A \implies f(x) \notin B$) is a universal statement over *all* instances, and remains a mathematical proof written separately for each reduction — the property test against the SAT oracle ([Section 6](06-pipeline-architecture.md)) **corroborates it empirically on sampled instances**, it does not replace it.

This distinction — a written proof for the universal statement, an executable check for the sample — is the thread that connects this report's theory to the code's architecture, the subject of the final section.

<br>

[§4](04-purpose-of-reductions.md) — §5 — [§6](06-pipeline-architecture.md)
