package commands

import (
	"os"
	"text/tabwriter"
	"text/template"
)

func printTemplate(t *template.Template, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 8,8,8, ' ', 0)
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
	_ = w.Flush()
}
