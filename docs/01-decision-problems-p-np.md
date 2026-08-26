# 1. Decision Problems, P and NP

## 1.1 Decision problems

A **decision problem** is formalized as a language $L \subseteq \Sigma^*$: a subset of the strings over a finite alphabet. An **instance** is a string $x$, and "solving" the problem on that instance means answering the question $x \in L$? with yes or no.

We work with decision problems, rather than directly with search problems ("find me a solution") or optimization problems ("find me the best solution"), for a precise technical reason: decision is the simplest possible form of the problem, and for the problems treated in this project the three forms are equivalent in difficulty via a **self-reducibility** argument — if you can decide in polynomial time "does an independent set of size $k$ exist?", you can find one explicitly with a polynomial number of calls to that same decider (fix a vertex in or out, recurse, and the yes/no answer at each step guides you). One caveat worth stating precisely: this equivalence is not a general theorem that follows automatically from NP-completeness — it is a property (called self-reducibility) proved case by case, problem by problem, with a constructive argument like the one just sketched; it holds for all the "natural" problems used here, but is not guaranteed for an arbitrary, artificially constructed NP-complete language. Studying decision remains the canonical starting point regardless.




## 1.2 The class P

A language $L$ is in **P** if there exists a deterministic algorithm that, for every instance $x$, decides $x \in L$ in a number of steps bounded by $p(\|x\|)$ for some fixed polynomial $p$.

The degree of the polynomial does not matter for the definition — only that a polynomial bound exists, uniform over all instances. This cutoff is deliberately coarse relative to practice (an $O(n^{100})$ algorithm is "efficient" by this definition but unusable) but it is the right one for theory, for a robustness reason discussed in [Section 2](02-purpose-of-complexity-classes.md).




## 1.3 The class NP: definition via verifier and certificate

The most operationally useful definition — and the one this project adopts as primary — is:

> $L \in NP$ if there exists a deterministic algorithm $V$ (the **verifier**) and a polynomial $p$ such that, for every $x$:
> $x \in L$ **if and only if** there exists a string $y$ (the **certificate**, or witness) with $\|y\| \le p(\|x\|)$ such that $V(x, y)$ accepts in time polynomial in $\|x\|$.

In words: belonging to NP does not mean "I can decide quickly", it means "if the answer is yes, a short proof of that fact exists, and that proof is quick to check". Searching for the certificate can be arbitrarily expensive; *verifying* it cannot.

Direct examples on the problems used in this project:

- **3-SAT**: certificate = a truth assignment to the variables. Verification = evaluate every clause, linear time in the size of the formula.
- **Independent Set**: certificate = the vertex set itself. Verification = check that no two chosen vertices are adjacent and that the size is $\ge k$, polynomial time in the number of chosen vertices.
- **Subset Sum**: certificate = the subset. Verification = sum its elements and compare to the target.

Note the common shape: in all three cases the certificate *is*, more or less, the object you were looking for, not an auxiliary artifact. That is not a coincidence — it is the reason why, in the project's code, the reconstruction function $g$ (which carries a certificate of the reduced problem back to a certificate of the original problem) is a natural part of the reduction rather than something bolted on: it simply makes the certificate structure of each problem explicit.




## 1.4 Equivalent definition: nondeterministic Turing machine

The more common textbook definition is alternative but equivalent: $L \in NP$ if there exists a nondeterministic Turing machine that decides $L$ in polynomial time (i.e., for every instance in $L$, an accepting computation path of polynomial length *exists*).

The equivalence with the verifier definition goes both ways with a direct argument: the nondeterministic choices along an accepting path *are* the certificate (just write them down in sequence); conversely, a deterministic verifier $V(x, y)$ is simulated by a nondeterministic machine that first "guesses" $y$ nondeterministically bit by bit, then runs $V$ deterministically.

The two definitions are theoretically interchangeable. This report and the code will almost always use the certificate one, because it is the operationally relevant one: a real SAT decider does not simulate nondeterminism — it produces a certificate (the assignment) when it answers SAT.




## 1.5 P ⊆ NP, and what remains open

$P \subseteq NP$ is immediate from the verifier definition: if $L \in P$, the verifier can ignore the certificate entirely and decide $x \in L$ on its own in polynomial time (equivalently: empty certificate, $V(x, y) :=$ "solve $x$ directly").

Whether the inclusion is strict ($P \neq NP$) is the most famous open problem in theoretical computer science, and this project makes no attempt to say anything about it. It is worth being precise about what the report *implicitly* assumes: it works under the presupposition (universally believed, never proven) that $P \neq NP$, because that presupposition is what makes it interesting for a problem to be NP-complete — if $P = NP$ turned out to be true one day, every reduction built here would remain correct, but the very notion of "intrinsic difficulty" motivating them would collapse.




## 1.6 An asymmetry worth keeping in mind

The definition of NP is **asymmetric**: it talks about short certificates for yes-instances, and says nothing about how to quickly recognize no-instances. It is not obvious (and is itself an open problem related to P vs NP) that no-instances also have short certificates of their own negativity — this property defines the class **co-NP**, and whether $NP = \text{co-NP}$ is unknown.

This asymmetry is not an erudite detail: a structurally similar one resurfaces when, in the section on reductions, the correctness of a reduction $f$ has to be argued in both directions, $x \in A \implies f(x) \in B$ and $x \notin A \implies f(x) \notin B$ — the second direction (preserving the *absence* of a solution) is typically the one that requires the real mathematical argument, for an analogous reason: there is no certificate to exhibit to prove it for a single example. The resemblance is a useful intuition pump, not a formal identity — the reduction's "no ⟹ no" direction is a universal claim about the function $f$, not itself a co-NP membership statement.

---

1 [→ 2](02-purpose-of-complexity-classes.md)
