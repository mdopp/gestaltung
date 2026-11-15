// Package gestaltung is the household's shared look: design tokens and base
// components as CSS (ADR 0001). A program embeds it at build time and puts
// CSS into its page's <style>; nothing loads it from another service at
// runtime.
package gestaltung

import (
	_ "embed"
	"html/template"
)

// Tokens are the custom properties: colours for dark and light, type,
// spacing and radius.
//
//go:embed tokens.css
var Tokens string

// Components are the base components, written only in tokens.
//
//go:embed components.css
var Components string

// CSS is Tokens followed by Components: what a page puts in its <style>.
var CSS = Tokens + "\n" + Components

// Style returns CSS typed for html/template, so it can be placed inside a
// <style> element without being escaped.
func Style() template.CSS { return template.CSS(CSS) }
