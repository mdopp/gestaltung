// Command vorschau writes a self-contained German preview page of every
// token and component, light and dark side by side:
//
//	go run ./cmd/vorschau > vorschau.html
package main

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"

	"github.com/mdopp/gestaltung"
	"github.com/mdopp/gestaltung/internal/tokens"
)

// What each colour is for, in the words the owner reads.
var colourUse = map[string]string{
	"bg": "Seitengrund", "panel": "Leisten", "card": "Karten", "edge": "Rand",
	"line": "Trennlinie", "ink": "Text", "muted": "Nebentext", "faint": "Linien, Symbole (kein Text)",
	"accent": "Hauptaktion, Links", "accent-ink": "Text auf Akzent", "accent-soft": "Fragen-Hintergrund",
	"ok": "In Ordnung", "ok-soft": "In Ordnung, Grund", "warn": "Bald handeln", "warn-soft": "Bald handeln, Grund",
	"bad": "Gestört", "bad-soft": "Gestört, Grund",
}

var order = []string{"bg", "panel", "card", "edge", "line", "ink", "muted", "faint", "accent", "accent-ink", "accent-soft", "ok", "ok-soft", "warn", "warn-soft", "bad", "bad-soft"}

type swatch struct{ Name, Value, Use string }

type theme struct {
	Key, Title string
	Colours    []swatch
}

const page = `<!doctype html>
<html lang="de">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Gestaltung Vorschau</title>
<style>
{{.CSS}}
/* preview layout only, tokens only */
body{margin:0;background:var(--bg);color:var(--ink);font:var(--text-m)/var(--leading) var(--font)}
.kopf{max-width:84rem;margin:0 auto;padding:var(--space-5) var(--space-4) 0}
.kopf h1{font-size:var(--text-xl);margin:0 0 var(--space-1)}
.kopf p{margin:0;color:var(--muted)}
.seiten{display:grid;grid-template-columns:1fr;gap:var(--space-4);max-width:84rem;margin:0 auto;padding:var(--space-4)}
@media (min-width:60rem){.seiten{grid-template-columns:1fr 1fr}}
.seiten>section{border:1px solid var(--edge);border-radius:var(--radius-m);overflow:hidden}
.seiten .page{max-width:none;padding-top:var(--space-1)}
.thema{font-size:var(--text-s);color:var(--muted);margin:var(--space-4) 0 0}
.farben{display:grid;grid-template-columns:repeat(auto-fill,minmax(9.5rem,1fr));gap:var(--space-2)}
.farbe{display:flex;gap:var(--space-2);align-items:center;font-size:var(--text-s);min-width:0}
.farbe i{flex:none;width:36px;height:36px;border-radius:var(--radius-s);border:1px solid var(--edge)}
.farbe code{display:block;color:var(--ink);font-size:.8rem}
.farbe span{color:var(--muted);line-height:1.25}
.skala{display:flex;flex-wrap:wrap;gap:var(--space-3);align-items:flex-end}
.skala .s{display:flex;flex-direction:column;align-items:center;gap:var(--space-1);font-size:.75rem;color:var(--muted)}
.skala .abstand{background:var(--accent);height:16px}
.skala .ecke{width:44px;height:44px;background:var(--card);border:2px solid var(--faint)}
.knoepfe{display:flex;flex-wrap:wrap;gap:var(--space-2)}
</style>
</head>
<body>
<header class="kopf">
<h1>So sehen die Programme im Haus künftig aus</h1>
<p>Eine Gestaltung für Laib, Homes, Abstimmung und Archiv. Links hell, rechts dunkel; jedes Gerät zeigt, was es eingestellt hat.</p>
</header>
<main class="seiten">
{{range .Themes}}
<section data-theme="{{.Key}}">
<div class="page">
<p class="thema">{{.Title}}</p>

<h2>Farben</h2>
<div class="farben">
{{range .Colours}}<div class="farbe"><i style="background:var(--{{.Name}})"></i><div><code>--{{.Name}}</code><span>{{.Use}}<br>{{.Value}}</span></div></div>
{{end}}</div>

<h2>Schrift</h2>
<div class="card">
<p style="font-size:var(--text-xl);font-weight:700;line-height:1.2">Guten Morgen, Lena</p>
<p style="font-size:var(--text-l);font-weight:600">Heute im Haus</p>
<p style="font-size:var(--text-m)">Die Waschmaschine ist fertig. Bitte bis Mittag aufhängen.</p>
<p style="font-size:var(--text-s)" class="muted">Systemschrift, Grundgröße 16 px</p>
</div>

<h2>Abstände und Ecken</h2>
<div class="skala">
{{range $.Spaces}}<div class="s"><div class="abstand" style="width:var(--{{.}})"></div>{{.}}</div>{{end}}
</div>
<div class="skala" style="margin-top:var(--space-3)">
{{range $.Radii}}<div class="s"><div class="ecke" style="border-radius:var(--{{.}})"></div>{{.}}</div>{{end}}
</div>

<h2>Seite und Karte</h2>
<h1>Aufgaben</h1>
<div class="card">
<p><b>Altpapier rausstellen</b></p>
<p class="muted">Morgen früh, bis 7 Uhr · Lena</p>
<div class="knoepfe"><button class="button">Erledigt</button><button class="button quiet">Später</button></div>
</div>

<h2>Knöpfe</h2>
<div class="knoepfe">
<button class="button">Speichern</button>
<button class="button quiet">Abbrechen</button>
<a class="button quiet" href="#">Alle Aufgaben</a>
<button class="button" disabled>Wird gesendet …</button>
</div>

<h2>Frage</h2>
<div class="question">
<p class="ask">Wohin geht der Ausflug am Sonntag?</p>
<p class="why">Abstimmung bis Freitag · 3 von 4 haben gewählt</p>
<div class="answers"><button class="button quiet">Zoo</button><button class="button quiet">Badesee</button><button class="button quiet">Wandern</button></div>
</div>

<h2>Liste</h2>
<ul class="rows">
<li><a class="row" href="#"><span class="what"><span class="name">Drucker im Arbeitszimmer</span><span class="note">Toner 60 %</span></span><span class="end"><span class="pill ok">Bereit</span></span></a></li>
<li><a class="row" href="#"><span class="what"><span class="name">Drucker im Flur</span><span class="note">Toner fast leer, Ersatz bestellen</span></span><span class="end"><span class="pill warn">Bald leer</span></span></a></li>
<li><a class="row" href="#"><span class="what"><span class="name">Scanner im Keller</span><span class="note">Nicht erreichbar seit 14:20</span></span><span class="end"><span class="pill bad">Gestört</span></span></a></li>
<li><a class="row" href="#"><span class="what"><span class="name">Fotodrucker</span><span class="note">Ausgeschaltet</span></span><span class="end"><span class="pill">Aus</span></span></a></li>
</ul>

<h2>Status</h2>
<div class="knoepfe"><span class="pill ok">Erledigt</span><span class="pill warn">Bald fällig</span><span class="pill bad">Überfällig</span><span class="pill">Pausiert</span></div>

<h2>Nichts zu tun</h2>
<div class="card empty"><b>Heute nichts zu tun.</b>Neue Aufgaben erscheinen hier.</div>
</div>
</section>
{{end}}
</main>
</body>
</html>
`

func swatches(m map[string]string) []swatch {
	var s []swatch
	for _, k := range order {
		if v, ok := m[k]; ok {
			s = append(s, swatch{k, v, colourUse[k]})
		}
	}
	return s
}

func prefixed(m map[string]string, p string) []string {
	var s []string
	for k := range m {
		if strings.HasPrefix(k, p) {
			s = append(s, k)
		}
	}
	sort.Strings(s)
	return s
}

func run() error {
	th, err := tokens.Parse(gestaltung.Tokens)
	if err != nil {
		return err
	}
	for k := range th.Dark {
		if _, ok := colourUse[k]; !ok {
			return fmt.Errorf("token --%s has no German description in cmd/vorschau", k)
		}
	}
	t := template.Must(template.New("v").Parse(page))
	return t.Execute(os.Stdout, map[string]any{
		"CSS": gestaltung.Style(),
		"Themes": []theme{
			{"light", "Hell", swatches(th.Light)},
			{"dark", "Dunkel", swatches(th.Dark)},
		},
		"Spaces": prefixed(th.Scale, "space-"),
		"Radii":  []string{"radius-s", "radius-m", "radius-pill"},
	})
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "vorschau:", err)
		os.Exit(1)
	}
}
