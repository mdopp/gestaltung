# gestaltung

The household's shared look: one set of design tokens and base components
for every program on a Laib box (Laib's pages, Homes, Abstimmung, the
archive). Why and the rules: [ADR 0001](docs/adr/0001-one-shared-design.md).

- `tokens.css`: colours for dark and light, type, spacing, radius, as CSS
  custom properties. Dark first, light when the device is; force one with
  `data-theme="light"` or `data-theme="dark"` on any element.
- `components.css`: `.page`, `.card`, `.button` (`.quiet`), `.question`,
  `.rows` / `.row`, `.pill` (`.ok` `.warn` `.bad`), `.empty`. Tokens only.

## Use

Go (embedded at build time):

```go
import "github.com/mdopp/gestaltung"

// in an html/template: <style>{{.Style}}</style>
data := map[string]any{"Style": gestaltung.Style()}
```

Anything else: `scripts/sync.sh <dest-dir> <tag>` writes
`<dest-dir>/gestaltung.css` from a pinned version. Commit or bundle it;
never load it from another service at runtime.

## Check

`go test ./...` runs the WCAG AA contrast check on every text token in both
themes, and fails if `components.css` carries a colour literal.
`go run ./cmd/vorschau > vorschau.html` writes the German preview page.
