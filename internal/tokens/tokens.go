// Package tokens reads tokens.css into per-theme maps and measures contrast
// (WCAG 2.x relative luminance). Used by the tests and the preview.
package tokens

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Themes holds each theme's colour tokens (name without "--" → value).
type Themes struct {
	Dark, Light map[string]string
	// LightBlocks counts the light blocks; the media block and the forced
	// one must say the same thing.
	LightBlocks []map[string]string
	// Scale holds the theme-independent tokens (type, spacing, radius).
	Scale map[string]string
}

var (
	blockRe = regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`)
	declRe  = regexp.MustCompile(`--([a-z0-9-]+)\s*:\s*([^;]+);`)
	comment = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// Parse reads tokens.css.
func Parse(css string) (Themes, error) {
	css = comment.ReplaceAllString(css, "")
	t := Themes{Scale: map[string]string{}}
	for _, m := range blockRe.FindAllStringSubmatch(css, -1) {
		sel := strings.TrimSpace(m[1])
		if i := strings.LastIndex(sel, "{"); i >= 0 {
			sel = strings.TrimSpace(sel[i+1:])
		}
		decls := map[string]string{}
		for _, d := range declRe.FindAllStringSubmatch(m[2], -1) {
			decls[d[1]] = strings.TrimSpace(d[2])
		}
		switch {
		case strings.Contains(sel, `[data-theme="dark"]`) && !strings.Contains(sel, ":not("):
			t.Dark = decls
		case strings.Contains(sel, `:not([data-theme="dark"])`), sel == `[data-theme="light"]`:
			t.LightBlocks = append(t.LightBlocks, decls)
			t.Light = decls
		case sel == ":root":
			for k, v := range decls {
				t.Scale[k] = v
			}
		default:
			return t, fmt.Errorf("unexpected selector %q in tokens.css", sel)
		}
	}
	if t.Dark == nil || t.Light == nil {
		return t, fmt.Errorf("tokens.css needs a dark and a light block")
	}
	return t, nil
}

// RGBA is a colour with alpha in 0..1.
type RGBA struct{ R, G, B, A float64 }

var rgbaRe = regexp.MustCompile(`^rgba?\(\s*([\d.]+)\s*,\s*([\d.]+)\s*,\s*([\d.]+)\s*(?:,\s*([\d.]+)\s*)?\)$`)

// ParseColor understands #rrggbb, #rgb and rgb()/rgba() with 0..255 channels.
func ParseColor(s string) (RGBA, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") {
		h := s[1:]
		if len(h) == 3 {
			h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
		}
		if len(h) != 6 {
			return RGBA{}, fmt.Errorf("bad hex colour %q", s)
		}
		v, err := strconv.ParseUint(h, 16, 32)
		if err != nil {
			return RGBA{}, fmt.Errorf("bad hex colour %q", s)
		}
		return RGBA{float64(v >> 16 & 255), float64(v >> 8 & 255), float64(v & 255), 1}, nil
	}
	m := rgbaRe.FindStringSubmatch(s)
	if m == nil {
		return RGBA{}, fmt.Errorf("not a colour %q", s)
	}
	c := RGBA{A: 1}
	c.R, _ = strconv.ParseFloat(m[1], 64)
	c.G, _ = strconv.ParseFloat(m[2], 64)
	c.B, _ = strconv.ParseFloat(m[3], 64)
	if m[4] != "" {
		c.A, _ = strconv.ParseFloat(m[4], 64)
	}
	return c, nil
}

// Over composites c over an opaque ground.
func (c RGBA) Over(g RGBA) RGBA {
	mix := func(a, b float64) float64 { return a*c.A + b*(1-c.A) }
	return RGBA{mix(c.R, g.R), mix(c.G, g.G), mix(c.B, g.B), 1}
}

func (c RGBA) luminance() float64 {
	lin := func(v float64) float64 {
		v /= 255
		if v <= 0.03928 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B)
}

// Contrast is the WCAG contrast ratio of two opaque colours.
func Contrast(a, b RGBA) float64 {
	la, lb := a.luminance(), b.luminance()
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}
