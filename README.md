# Karp All-Reductions Pipeline

Una pipeline di riduzioni esplicite e tipizzate tra problemi NP-completi, verificata end-to-end contro un SAT solver reale — "reduction as compilation". Il nome è quello di Richard Karp, l'autore delle riduzioni many-one polinomiali che il progetto implementa e dei 21 problemi NP-completi che le collegano.

## Cos'è

Ogni riduzione tra due problemi NP-completi (es. 3-SAT → Independent Set) è implementata come una funzione pura `f: IstanzaA → IstanzaB`, accompagnata da una mappa inversa sui certificati `g: CertificatoB → CertificatoA`. La correttezza non è solo dichiarata: viene messa alla prova con property-based testing, confrontando su istanze generate casualmente la risposta di un SAT solver reale (per 3-SAT) contro un oracolo di riferimento per il problema ridotto, e verificando che `g` ricostruisca sempre un certificato valido.

Catena implementata (deliberatamente piccola, per restare un progetto autocontenuto):

```
3-SAT → Independent Set → Vertex Cover (→ Clique per complemento)
3-SAT → Subset Sum
```

3-SAT è preso come radice della catena — la sua NP-completezza è citata (teorema di Cook–Levin) e assunta, non ridimostrata in codice.

**Linguaggio: Go.** Deciso su R dopo un confronto diretto su tre assi (chiarezza del codice di riduzione, ergonomia del property testing, maturità del solver SAT disponibile in ciascun linguaggio) — l'ultimo è stato quello decisivo, dato che l'unico binding SAT in-process per R (`rpicosat`) risulta archiviato da CRAN dal 2022.

## Stato del progetto

- [x] Relazione teorica (sezioni 1–6)
- [ ] Implementazione Go delle riduzioni
- [ ] Property test contro l'oracolo SAT

## Relazione teorica

Il documento completo è in [`docs/relazione.md`](docs/relazione.md). Indice:

1. [Problemi di decisione, P e NP](docs/relazione.md#1-problemi-di-decisione-p-e-np) — linguaggi, la classe P, NP a verificatore+certificato, l'equivalenza con le macchine di Turing non deterministiche, e l'asimmetria NP/co-NP.
2. [A cosa servono le classi di complessità](docs/relazione.md#2-a-cosa-servono-le-classi-di-complessità) — classificano il problema non l'algoritmo, sono robuste rispetto al modello di calcolo, danno un vocabolario comune a domini diversi, trasferiscono risultati negativi.
3. [Riduzioni many-one polinomiali](docs/relazione.md#3-riduzioni-many-one-polinomiali) — la definizione di Karp, e perché il progetto usa solo quella (e non le più generali riduzioni di Turing/Cook).
4. [A cosa servono le riduzioni](docs/relazione.md#4-a-cosa-servono-le-riduzioni) — trasferimento di difficoltà in entrambe le direzioni, l'ordine indotto da `≤p`, e la nozione di completezza.
5. [Come si usano in pratica](docs/relazione.md#5-come-si-usano-in-pratica-cooklevin-come-radice-poi-tutto-per-riduzione) — il teorema di Cook–Levin come radice, la normalizzazione SAT → 3-SAT, e la ricetta a due passi usata da ogni riduzione successiva.
6. [Architettura della pipeline](docs/relazione.md#6-architettura-della-pipeline) — i componenti del codice (istanze tipizzate, riduzioni, oracoli, confine DIMACS) e la loro corrispondenza uno a uno con le sezioni precedenti.
