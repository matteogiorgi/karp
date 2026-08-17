# 3. Polynomial Many-One Reductions

## 3.1 Definition

Given two decision problems $A, B \subseteq \Sigma^\*$, a **polynomial many-one reduction** (or **Karp reduction**) from $A$ to $B$ is a function $f: \Sigma^\* \to \Sigma^\*$ such that:

1. $f$ is **total and computable in polynomial time** (a deterministic algorithm exists that computes $f(x)$ in time $\le p(\|x\|)$ for every $x$, for some fixed polynomial $p$);
2. for every instance $x$: $x \in A \iff f(x) \in B$.

This is written $A \le_p B$ ("$A$ reduces to $B$"), and $f$ is read as a **witness** of that relation: its mere existence is the proof that $A \le_p B$ holds.

The name "many-one" describes the shape of the function: $f$ may send different instances of $A$ to the same instance of $B$ (it need not be injective), but every instance of $A$ is transformed into **exactly one** instance of $B$ — you are not allowed to construct several instances of $B$, query them, and combine the answers. This is exactly the sense of "reduction as compilation" adopted by the project: $f$ is a pure function, $A$-instance in, $B$-instance out, with no intermediate state and no combination logic.

## 3.2 Distinction from Turing/Cook reductions

The definition above is a special, more restrictive case of a more general notion: the **Turing reduction** (or Cook reduction). $A \le_T B$ means there exists a deterministic polynomial algorithm for $A$ that has access to an **oracle** for $B$ — you may query the oracle a polynomial number of times, **adaptively** (the $(i+1)$-th query may depend on the answer to the $i$-th one), and you may do whatever you want with the answers before producing the final verdict on $A$ (including negating them, combining them with AND/OR, etc.).

Every many-one reduction is also a Turing reduction (just make a single query, with $f(x)$, and return its answer unchanged), but the converse does not hold: a Turing reduction can do things a many-one reduction cannot — for instance, deciding $A$ by querying the oracle both on $f(x)$ and on its "complement" and comparing the two answers.

This project needs — and needs only — Karp's many-one notion, for a reason that goes beyond mere economy of means: **a single call, with no post-processing of the answer, is exactly the structure that makes reconstructing a certificate possible.** If $f(x) \in B$ with certificate $y_B$, and the transformation $f$ is known explicitly, one can often read $y_B$ "backward" through the structure of $f$ and obtain a certificate $y_A$ for $x \in A$ — this is precisely the role of the function $g$ introduced in the project's code ($g: \text{certificate}_B \to \text{certificate}_A$). With a generic Turing reduction, involving multiple adaptive queries and arbitrary combination logic, this reconstruction has in general no clean form: there is no single instance of $B$ to "read backward" from.

In the literature, when NP-completeness is referred to generically without specifying the type of reduction, the many-one notion is almost always the one meant — it is the notion Karp used in his 1972 paper on the 21 NP-complete problems, and it is the one adopted here.
