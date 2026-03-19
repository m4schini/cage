package nix

import (
	_ "embed"
	"strings"
	"text/template"
)

type ShellNixPackages struct {
	Packages []string
	Shell    string
}

var funcMap = template.FuncMap{
	"joinLines": func(pkgs []string) string {
		return strings.Join(pkgs, "\n    ")
	},
}

//go:embed shell.nix.tmpl
var shellNix string

var shellNixTemplate = mustParse(template.
	New("shell.nix").
	Funcs(funcMap).
	Parse(shellNix),
)

func mustParse(tmpl *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}

	return tmpl
}
