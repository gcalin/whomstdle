# Traits

Proposed attributes stored per person. Each one backs a question the player can ask. Mark up what to keep, cut or change.

Design rules: every answer is short and checkable, and the set works for a global audience (no country-specific knowledge needed to ask or understand a question).

| # | Trait | Question | Answer type | Example |
| --- | --- | --- | --- | --- |
| 1 | Alive | Are they alive? | yes / no | No |
| 2 | Gender | What is their gender? | man / woman / other | Woman |
| 3 | Field | What are they known for? | film & TV, music, sport, politics & leadership, science & tech, business, literature & art, other | Music |
| 4 | Continent | Which continent are they from? | Africa, Asia, Europe, N. America, S. America, Oceania | Europe |
| 5 | Birth decade | When were they born? | decade (1920s … 2000s) | 1980s |
| 6 | Active era | In which decades were they most famous? | one or more decades | 2000s, 2010s |
| 7 | Languages | Which languages did they speak? | list of languages (open) | French, English |
| 8 | Public role | Do they perform in front of an audience (screen, stage, field)? | yes / no | Yes |
| 9 | Awards | Have they won a major international award (Oscar, Grammy, Olympic medal, Nobel…)? | yes / no | Yes |
| 10 | Hint | Give me a one-sentence hint | free text | "Known for a hit song about …" |

## Notes

- **Limit:** none. Players can ask as many questions as they like.
- **Active era** can span several decades. **Languages** is an open question: the answer lists them, and a person can have several.
- **Gender:** include "other" so the answer is never wrong, and treat it as a rare case.
- **Field and continent** are the strongest questions: they split the pool well without needing local knowledge. **Languages** and **Public role** are cheap yes/no questions that narrow it further.
- **Hint** is a free-text trait, not a category. Make it the last resort, and consider counting it as extra score.
- **Left out on purpose:** exact age, nationality (too many values, hard to guess), weight and appearance (sensitive and subjective), net worth, and relationship or family details.
- **Open questions:**
  - Without a limit, the score can't be the number of questions alone. Decide the rule later (see the plan's Phase 3).
  - Dead people: should the answer to "Birth decade" and "Active era" still apply? (Suggested: yes.)
  - Should "Awards" stay, since it depends on how well-known the person is?

## Question aliases

Players will type a question and pick from a filtered list, so each question lists other ways to ask it. Typing any alias finds the canonical question.

| Trait | Canonical question | Aliases |
| --- | --- | --- |
| Alive | Are they alive? | Are they still alive? · Are they dead? · Have they died? · Are they living? · Did they pass away? |
| Gender | What is their gender? | Are they a man or a woman? · Are they male or female? · Which sex are they? |
| Field | What are they known for? | Why are they famous? · What do they do? · What is their job? · What is their profession? · What is their field? · What made them famous? |
| Continent | Which continent are they from? | Where are they from? · What part of the world are they from? |
| Birth decade | When were they born? | Which decade were they born in? · What decade were they born? |
| Active era | In which decades were they most famous? | When were they famous? · When were they active? · When did they peak? · Which era are they from? |
| Languages | Which languages did they speak? | What language do they speak? · What is their mother tongue? · What is their native language? |
| Public role | Do they perform in front of an audience? | Are they a performer? · Do they perform live? |
| Awards | Have they won a major international award? | Have they won any awards? · Are they award-winning? |
| Hint | Give me a one-sentence hint | Can I get a hint? · Give me a clue · Tell me something about them · I'm stuck · Help |

Aliases keep the same meaning as the canonical question. Narrower questions ("Do they speak English?", "Are they European?") are not aliases.

**Later implementation: loose matching.** Match what the player types against questions and aliases by lowercasing, ignoring punctuation and matching on words, so the lists need few exact variants.
