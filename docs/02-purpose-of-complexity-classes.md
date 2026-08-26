# 2. What Complexity Classes Are For

## 2.1 They classify the problem, not an algorithm

A statement like "Vertex Cover is NP-complete" is not a claim about a specific algorithm — it is a claim about the entire problem: no algorithm, among all possible ones, solves it in polynomial time in the worst case, unless $P = NP$. This is an important distinction: knowing that a problem sits in a class says something about its intrinsic structure, independent of how good you are at writing code.

This is also why a proof of membership in NP (via certificate, [Section 1.3](01-decision-problems-p-np.md#13-the-class-np-definition-via-verifier-and-certificate)) and a proof of NP-completeness (via reduction, [Section 3](03-polynomial-many-one-reductions.md)) are so different in form: the first exhibits a concrete object (the verifier algorithm); the second is an *indirect, proof-by-contradiction-style* argument — it shows that a hypothetical fast algorithm for the problem would imply one for every problem in NP, which is believed to be impossible.




## 2.2 They are robust with respect to the model of computation

The "polynomial time" threshold is not arbitrary: it is the coarsest threshold that stays stable across reasonable changes of computational model. Single-tape Turing machine, multi-tape, RAM machine with unit-cost arithmetic — these models simulate one another with **polynomial** overhead (typically quadratic or less). So "a polynomial algorithm exists" is a property of the *problem*, not of the model chosen to formalize it: if it is true in one model, it is true in all of them.

This is the same idea as the Church–Turing thesis, applied not to computability but to efficiency: the exact exponent (which depends on the model and implementation details) does not matter; the qualitative polynomial/non-polynomial threshold is what matters, because it is the only one that stays invariant. It is also why the definition of P does not specify the degree of the polynomial: doing so would make it model-dependent, and therefore theoretically less interesting — even though, as already noted, it is less useful for distinguishing "efficient" from "inefficient" in actual practice.




## 2.3 They give a common vocabulary across different domains

The most immediate practical value: complexity classes let you compare problems that, on the surface, have nothing in common. 3-SAT is a propositional logic problem; Independent Set is a graph theory problem; Subset Sum is an arithmetic problem. There is no obvious way to say "these three problems are equally hard" while staying inside the language of each domain — the notion of complexity class (and, downstream, of reduction) is precisely the bridge that turns that into a precise, provable statement. It is exactly the bridge this project crosses explicitly in code: the same "difficulty" manifests as a formula, as a graph, and as a set of numbers.




## 2.4 They transfer negative results, and they guide practice

Saying "I haven't found a polynomial algorithm for B" is a weak statement — it might just mean you haven't tried hard enough. Saying "B is NP-complete" is a strong statement: it means that *nobody*, across fifty years of research on thousands of interconnected NP-complete problems, has found a polynomial algorithm for *any* of them — and finding one for B would automatically solve all the others. This is the sense in which classes "are for" something: not so much solving the problem as giving a rigorous, transferable answer to the question "why is this hard, and how much can you trust that it will stay hard".

On the practical side, this negative result is also a guide: once a problem is known to be NP-complete, the right question stops being "how do I find the exact polynomial algorithm" and becomes "what trade-off do I accept" — approximation, heuristics, restriction to special cases, or worst-case-exponential-but-practically-effective algorithms. Modern SAT solvers (the oracle used in this very project) are the concrete example of this last path: NP-complete in the worst case does not mean "unusable" — it means "no worst-case polynomial guarantee" — and there is a lot of space, exploited industrially, between the two.




## 2.5 A bridging note toward reductions

P and NP, as defined here, talk about direct computation on a single problem: an algorithm, or a verifier, for *that* language. On their own they say nothing yet about how two problems compare to each other. The tool that builds that comparison — and that gives NP its internal structure, identifying within it the problems that are "hardest of all" — is the polynomial reduction, the subject of the next section.

</br>

[§1](01-decision-problems-p-np.md) — §2 — [§3](03-polynomial-many-one-reductions.md)
