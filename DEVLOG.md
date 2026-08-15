# Development Log — Issues Found and Fixed

This is a running record of concrete bugs and design gaps found while building this project, and what was done about them. It exists because most of these were not obvious from reading the code after the fact — each one was caught by deliberately rechecking finished work rather than trusting that "it compiles" or "the test I just wrote passes" meant it was correct, and several were verified empirically (with a throwaway script or a local build) before being written down as a fact rather than an assumption. Keeping that trail is more useful than letting it disappear once the fix lands.

Organized in two parts: issues in the GitHub Pages / report rendering pipeline, then issues in the Go implementation.

## GitHub Pages / report rendering

### 1. Math formulas didn't render on GitHub Pages

GitHub.com's own file browser has a built-in LaTeX renderer for `$...$` / `$$...$$`; GitHub Pages builds the same Markdown through Jekyll and kramdown, which has no such feature on its own. Formulas that looked correct on github.com showed up as literal dollar-sign text on the published site.

**Fix**: added MathJax v3 via `_includes/head-custom.html` — the officially supported customization point for GitHub Pages' default theme — configured to recognize both `$...$` and `$$...$$`.

### 2. The Mermaid diagram didn't render on GitHub Pages

Same root cause as #1: GitHub.com renders ` ```mermaid ` fenced blocks natively; kramdown just emits them as a plain `<pre><code class="language-mermaid">` block of escaped text.

**Fix**: added Mermaid.js (ESM build) to the same `head-custom.html`, with a small script that finds every `pre > code.language-mermaid`, replaces it with a `div.mermaid` holding the raw text, and calls `mermaid.run()`.

### 3. `\Sigma^*` got corrupted by kramdown's emphasis parsing

Kramdown's GFM parser treats two unescaped `*` characters anywhere in the same paragraph as a `*emphasis*` pair — including two `*` that each belong to a separate, unrelated `$...$` formula. The sentence `` $A, B \subseteq \Sigma^*$ ... $f: \Sigma^* \to \Sigma^*$ `` has three such asterisks in one paragraph; kramdown paired two of them and rendered a chunk of the sentence in italics with the LaTeX itself broken in the middle. Reproduced and confirmed the exact mechanism with a local kramdown test harness (the same `github-pages` gem GitHub itself uses) before touching anything.

**Fix**: escaped the asterisk as `\Sigma^\*` — a Markdown-level escape, stripped by kramdown before the text reaches MathJax, so the rendered LaTeX is unaffected. Re-verified with the same test harness.

### 4. `|x|`, `|y|` (absolute-value bars) got auto-detected as Markdown tables

Kramdown's GFM table parser treats *any* line containing two or more literal `|` characters as an implicit single-row table — even without the delimiter row the GFM spec itself requires. Every formula using absolute-value notation (`|x|`, `p(|x|)`, and so on — used throughout the polynomial-time definitions, i.e. constantly) was silently split into garbled table cells on the published site. Found from a user screenshot, reproduced locally to confirm the mechanism before fixing.

**Fix**: escaped every `|` inside a formula as `\|`, for the same reason as #3 — a literal character survives to MathJax, the table parser never sees a bare pipe.

### 5. Footer navigation links resolved to raw `.md` instead of `.html` on Pages

GitHub Pages' `jekyll-relative-links` plugin rewrites `.md` links to `.html` only when they are genuine Markdown link syntax (`[text](file.md)`) that kramdown itself turns into an `<a>` tag — it does not touch a hand-written `<a href="file.md">`. The section footers had been switched to raw HTML specifically to get left/center/right alignment (see #6), which broke this rewriting without anyone intending to touch it.

**Fix**: switched the footer back to genuine Markdown table syntax. Along the way, discovered that kramdown requires at least one data row *below* a header-plus-delimiter pair to recognize a table at all — a header-only "table" (which is all a single visible row of links would naturally be) is not recognized as a table and falls back to a plain paragraph with literal pipe characters. Worked around this by using an empty header row, keeping the actual links as a real data row.

### 6. The footer table wasn't actually centered

Even as a real table, it had no explicit width and shrank to fit its content, sitting flush left on the page. Forcing `width: 100%` fixed the span but not the centering: `justify-content: space-between` (flexbox) guarantees equal *gaps* between items, not that the middle item sits on the container's true center — an empty "Previous" cell (the very first section has no previous link) pulled "README" visibly off-center.

**Fix**: switched to CSS Grid with `grid-template-columns: 1fr auto 1fr`. The middle, auto-sized column stays mathematically centered regardless of what the flanking columns contain, empty or not. Caught and fixed a second bug in the same change: the grid rule had accidentally been applied to both `tbody` and `tr`, nesting two grids and compressing the whole row into roughly half the page width.

### 7. `vendor/` directory collision between Ruby and Go

Jekyll's bundler had been installing gems into `vendor/bundle`. The moment a Go module was added to the same repository, Go's toolchain auto-detected the top-level `vendor/` directory as its own vendoring convention and every `go` command (`go build`, `go doc`, ...) failed with an "inconsistent vendoring" error.

**Fix**: relocated the Ruby gem path to `.bundle/vendor`, removing the name collision at its root instead of passing `-mod=mod` to every future Go command by hand.

## Go implementation

### 8. `IndependentSet.Verify` accepted certificates naming vertices outside the graph

`Verify` counted every distinct index in a certificate toward the required size `K` without ever checking that the index was a real vertex (`0 <= v < NumVertices`). A certificate padded with enough out-of-range "phantom" indices could reach size `K` and be accepted, because no edge in the graph ever touches a vertex that doesn't exist — so the independence check could never disqualify it. This is a genuine violation of the formal definition: an independent set is by definition a subset of `V`, not an arbitrary set of integers.

**Fix**: added the bounds check at the top of `Verify`, with a dedicated regression test (a certificate of `{0, 99}` on a 4-vertex graph, which must be rejected despite reaching the required size).

### 9. Stale doc comment

`IndependentSet`'s doc comment still said `ThreeSATToIndependentSet` was "not yet written in this package," left over from before the reduction was actually written in its own file.

**Fix**: updated the comment to point at `threesat_to_independent_set.go`.

### 10. The property-test generator produced ~97.5% satisfiable instances

With `numVars` and `numClauses` both drawn independently and uniformly from `[1,5]`, measuring 2000 generated instances showed only 2.5% were unsatisfiable — meaning the "no" side of the pipeline's invariant (the real SAT oracle and the brute-force oracle agreeing a formula is *un*satisfiable) was exercised in only 2–3 of every 100 test cases. Measured before touching the generator, not assumed from first principles.

**Fix**: after trying several ranges empirically and measuring each one, settled on `numVars ∈ [1,3]`, `numClauses ∈ [3,6]`, which raises the unsatisfiable share to roughly 8% — a meaningfully healthier mix for the same sample size.

### 11. The brute-force oracle's safety cap does not imply a fast test suite

`maxBruteForceUniverse = 30` (in `oracle.go`) exists to turn a would-be-infinite enumeration into an immediate panic — it says nothing about how *slow* an instance under that cap can still be. Measuring the real worst case (a genuinely unsatisfiable instance, which forces full enumeration of every candidate) gave 13ms at 5 clauses, 140ms at 6, 1.35s at 7, ~13s at 8, and over two minutes at 9 — worse than the raw doubling of the subset count alone would suggest, because each candidate also costs O(edges) to check, and edge count grows with clause count too. The first fix attempted for #10 (raising `numClauses` up to 10) would have stayed safely within the 30-clause panic threshold while still risking a multi-minute single test case — especially once rapid's shrinking re-runs the property function many times over while searching for a minimal counterexample.

**Fix**: capped the property-test generator's `numClauses` at 6, independently of and much stricter than the 30-clause panic threshold, keeping the measured worst case around 140ms. The reasoning and the measurements themselves are recorded directly in `genThreeSAT`'s doc comment in `property_test.go`, not just here.
