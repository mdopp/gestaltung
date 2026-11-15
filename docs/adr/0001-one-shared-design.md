# ADR 0001 — One shared design, in its own package

Applies when you style a page in any household program (Laib's admin
pages, Homes, Abstimmung, the archive, or a new one), when you are about to
write a colour, a font size or a radius into a program's CSS, or when you
change a token here.

- **Status:** accepted, 2026-11-15 (the owner, after the preview: „Mir reicht die Vorschau“)
- **Deciders:** the owner (mdopp), on the question where the shared design
  should live: „Eigenes Paket“. Phase 2 (migrating the programs) waits for
  the owner's look at the preview.

## Context

Every program carries its own copy of the CSS, and the copies disagree:

- Laib (`internal/pages/laib.css`) and Homes (`internal/pages/family.go`)
  share a dark scheme: `--bg:#070709`, `--panel:#0c0c0e`, `--card:#14141a`,
  `--ink:#f2f2f5`, orange accent `#f97316`. The archive (`seite.go`) copied
  the same values again.
- Abstimmung (`apps/abstimmung/page.go`) is cream with a brown accent and
  its own token names (`--tx`, `--ak`, `--fl`).
- Inside Homes, pages drift: `kontakt_frage.go` has its own light grey
  scheme, `genutzt.go` a light-only token block of its own, and many pages
  write colours straight into their markup.
- The existing light variant fails WCAG AA for text in places: the light
  accent `#ea580c` is 3.6:1 on white, `--ok` `#059669` 3.8:1, `--muted`
  `#64748b` 4.3:1 on the page ground.

A family member moving from Homes to a poll to the archive meets three
different looks for the same kind of thing. The owner's principles for
family pages are calm, first-sight clarity („could my mother use it?“), no
persona gimmicks, and „rigoros einfach“.

The layer rule (laib ADR 0035) says dependencies point only down and no
program may depend on Homes. So the design cannot live in Homes, and
Laib, at layer 0, must be able to use it too.

## Decision

1. **The design is its own building block: `mdopp/gestaltung`.** It holds
   design tokens (`tokens.css`) and base components (`components.css`):
   page, card, button, question box, list row, status pill and the quiet
   empty state. It has no runtime and depends on nothing, so every layer,
   Laib included, may depend on it.
2. **Build time, never runtime.** Go programs import the module and embed
   `gestaltung.CSS` (`gestaltung.Style()` for `html/template`). Other
   programs copy a pinned version with `scripts/sync.sh <dir> <tag>`. No
   program loads the CSS from another service: a page must render when
   every other service is down.
3. **The palette is the one already shared.** Dark first, from the
   Laib/Homes scheme unchanged. Light when the device is light, with the
   colours corrected to pass AA (`--accent:#c2410c`, `--ok:#047857`,
   `--warn:#a14a07`, `--bad:#b91c1c`, `--muted:#5b6474`). System font, 16 px
   base, one column of at most 40 rem. Abstimmung's cream scheme is not
   carried over.
4. **Programs use the tokens and carry no foreign colour literals.** A
   program's own CSS and markup refer to colours only as `var(--…)` from
   this package; no hex, `rgb()`/`rgba()` or `hsl()` values, no own colour
   token names. Type sizes, spacing and radii use the tokens too. What the
   shared components do not cover, a program builds from tokens. If a
   program needs a new colour, the token is added here first. The only
   exceptions are content that is not chrome (a user's photo, a calendar
   colour the user chose, an embedded chart's data series), and they are
   named in the program's check allowlist.
5. **The rule is enforced by a CI check in each program**, added in phase
   2: it fails the build on a colour literal outside the embedded or
   synced copy of gestaltung and outside the allowlist. This package keeps
   the same rule itself: its tests fail if `components.css` carries a
   colour literal.
6. **Every text token passes WCAG AA (4.5:1) in both themes**, on every
   ground it is meant for; `--faint` and `--edge` are for borders and
   icons only, `--faint` at 3:1. `go test` measures it; a token change
   that fails does not merge.
7. **Versions are tags.** Programs pin a tag and move to a new one in an
   ordinary PR, so a change here never restyles a program without its own
   review.

## Consequences

- One place to change the look; a fix to a contrast or a radius reaches
  every program on its next version bump.
- Phase 2 is real work per program: replace the copied CSS, rename
  divergent tokens (`--danger` → `--bad`, Abstimmung's `--tx`/`--ak`/…),
  and remove the inline colours from markup. Abstimmung changes most,
  from cream to the shared scheme.
- The light theme looks slightly darker in its accent and status colours
  than today's Laib light variant; that is the price of AA.
- A program cannot quickly add „just one colour“; it adds a token here
  first. That friction is intended.
- The preview (`go run ./cmd/vorschau`) shows every token and component in
  both themes with household content; it is the page the owner judges
  before phase 2.
