package tokens

import (
	"math"
	"testing"
)

// The contrast maths must be able to fail: known WCAG values.
func TestContrastKnownValues(t *testing.T) {
	w, _ := ParseColor("#fff")
	b, _ := ParseColor("#000000")
	g, _ := ParseColor("#777777")
	for _, c := range []struct {
		a, b RGBA
		want float64
	}{{w, b, 21}, {w, w, 1}, {w, g, 4.48}} {
		if got := Contrast(c.a, c.b); math.Abs(got-c.want) > 0.01 {
			t.Errorf("contrast = %.3f, want %.2f", got, c.want)
		}
	}
}

func TestOver(t *testing.T) {
	half, _ := ParseColor("rgba(255,255,255,.5)")
	b, _ := ParseColor("#000000")
	if got := half.Over(b); math.Abs(got.R-127.5) > 0.01 || got.A != 1 {
		t.Fatalf("over = %+v", got)
	}
}

func TestParseRejects(t *testing.T) {
	for _, s := range []string{"#12", "red", "#gggggg"} {
		if _, err := ParseColor(s); err == nil {
			t.Errorf("%q parsed", s)
		}
	}
	if _, err := Parse(":root{--x:1px;}"); err == nil {
		t.Error("tokens without themes parsed")
	}
	if _, err := Parse(`.foo{--x:#000;}`); err == nil {
		t.Error("unknown selector parsed")
	}
}
