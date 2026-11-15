package gestaltung

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/mdopp/gestaltung/internal/tokens"
)

// A text pair: foreground token on a ground; the ground may itself be a
// translucent tint over another ground ("ok-soft/card").
type pair struct {
	fg, ground string
	min        float64
}

// textPairs is every way a text token is meant to be used. WCAG AA: 4.5:1
// for text, 3:1 for borders and icons (--faint).
var textPairs = func() []pair {
	var p []pair
	for _, fg := range []string{"ink", "muted", "accent", "ok", "warn", "bad"} {
		for _, g := range []string{"bg", "panel", "card"} {
			p = append(p, pair{fg, g, 4.5})
		}
	}
	p = append(p,
		pair{"accent-ink", "accent", 4.5},
		pair{"ink", "accent-soft/card", 4.5},
		pair{"muted", "accent-soft/card", 4.5},
		pair{"ok", "ok-soft/card", 4.5},
		pair{"warn", "warn-soft/card", 4.5},
		pair{"bad", "bad-soft/card", 4.5},
		pair{"muted", "line/card", 4.5},
		pair{"faint", "card", 3},
		pair{"faint", "bg", 3},
	)
	return p
}()

func themes(t *testing.T) tokens.Themes {
	t.Helper()
	th, err := tokens.Parse(Tokens)
	if err != nil {
		t.Fatal(err)
	}
	return th
}

func color(t *testing.T, theme map[string]string, ref string) tokens.RGBA {
	t.Helper()
	parts := strings.Split(ref, "/")
	var c tokens.RGBA
	for i := len(parts) - 1; i >= 0; i-- {
		v, ok := theme[parts[i]]
		if !ok {
			t.Fatalf("token --%s missing", parts[i])
		}
		x, err := tokens.ParseColor(v)
		if err != nil {
			t.Fatal(err)
		}
		if i == len(parts)-1 {
			if x.A != 1 {
				t.Fatalf("ground --%s must be opaque", parts[i])
			}
			c = x
			continue
		}
		c = x.Over(c)
	}
	return c
}

func TestContrastAA(t *testing.T) {
	th := themes(t)
	all := map[string]map[string]string{"dark": th.Dark}
	for i, b := range th.LightBlocks {
		all[[]string{"light (media)", "light (forced)"}[i%2]] = b
	}
	for name, theme := range all {
		for _, p := range textPairs {
			got := tokens.Contrast(color(t, theme, p.fg), color(t, theme, p.ground))
			if got < p.min {
				t.Errorf("%s: --%s on %s is %.2f:1, needs %.1f:1", name, p.fg, p.ground, got, p.min)
			}
		}
	}
}

func keys(m map[string]string) []string {
	var k []string
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

func TestThemesDefineTheSameTokens(t *testing.T) {
	th := themes(t)
	if d, l := strings.Join(keys(th.Dark), ","), strings.Join(keys(th.Light), ","); d != l {
		t.Fatalf("dark and light differ:\n dark  %s\n light %s", d, l)
	}
	if len(th.LightBlocks) != 2 {
		t.Fatalf("want the media light block and the forced one, got %d", len(th.LightBlocks))
	}
	a, b := th.LightBlocks[0], th.LightBlocks[1]
	for _, k := range keys(a) {
		if a[k] != b[k] {
			t.Errorf("light blocks disagree on --%s: %s vs %s", k, a[k], b[k])
		}
	}
	if len(a) != len(b) {
		t.Errorf("light blocks have %d and %d tokens", len(a), len(b))
	}
}

var literal = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(`)

// The components keep the rule every program will follow: tokens only.
func TestComponentsUseOnlyTokens(t *testing.T) {
	if m := literal.FindAllString(Components, -1); m != nil {
		t.Fatalf("components.css carries colour literals: %v", m)
	}
	th := themes(t)
	defined := map[string]bool{}
	for k := range th.Dark {
		defined[k] = true
	}
	for k := range th.Scale {
		defined[k] = true
	}
	for _, m := range regexp.MustCompile(`var\(--([a-z0-9-]+)\)`).FindAllStringSubmatch(Components, -1) {
		if !defined[m[1]] {
			t.Errorf("components.css uses undefined token --%s", m[1])
		}
	}
}

func TestEveryComponentIsThere(t *testing.T) {
	for _, c := range []string{".page", ".card", ".button", ".question", ".rows", ".row", ".pill", ".empty"} {
		if !strings.Contains(Components, c+"{") && !strings.Contains(Components, c+" ") {
			t.Errorf("component %s missing", c)
		}
	}
	if !strings.HasPrefix(CSS, Tokens) || !strings.HasSuffix(CSS, Components) {
		t.Fatal("CSS must be Tokens then Components")
	}
	if string(Style()) != CSS {
		t.Fatal("Style() must return CSS")
	}
}
