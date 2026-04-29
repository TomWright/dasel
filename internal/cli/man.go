package cli

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/alecthomas/kong"
	"github.com/tomwright/dasel/v3/internal"
)

type ManCmd struct{}

type manFlag struct {
	Name  string
	Short string
	Help  string
}

type manSubcommand struct {
	Name  string
	Help  string
	Flags []manFlag
}

type manData struct {
	Name        string
	Description string
	Version     string
	Date        string
	Subcommands []manSubcommand
	GlobalFlags []manFlag
	QueryFlags  []manFlag
}

const dateFormat = "2006-01-02"

func getReproducibleDate() (string, error) {
	sourceDateEpochEnv, ok := os.LookupEnv("SOURCE_DATE_EPOCH")

	// Fallback to current date if SOURCE_DATE_EPOCH is unset or empty.
	if !ok || sourceDateEpochEnv == "" {
		return time.Now().Format(dateFormat), nil
	}

	sde, err := strconv.ParseInt(sourceDateEpochEnv, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid SOURCE_DATE_EPOCH: %q", sourceDateEpochEnv)
	}

	return time.Unix(sde, 0).UTC().Format(dateFormat), nil
}

func extractManData(k *kong.Kong) manData {
	app := k.Model

	data := manData{
		Name:        strings.ToUpper(app.Name),
		Description: app.Help,
		Version:     internal.Version,
	}

	for _, flag := range app.Flags {
		if flag.Hidden {
			continue
		}
		mf := manFlag{
			Name: flag.Name,
			Help: flag.Help,
		}
		if flag.Short != 0 {
			mf.Short = string(flag.Short)
		}
		data.GlobalFlags = append(data.GlobalFlags, mf)
	}

	for _, child := range app.Children {
		if child.Hidden {
			continue
		}
		if child.Type != kong.CommandNode {
			continue
		}
		sub := manSubcommand{
			Name: child.Name,
			Help: child.Help,
		}
		for _, flag := range child.Flags {
			if flag.Hidden {
				continue
			}
			mf := manFlag{
				Name: flag.Name,
				Help: flag.Help,
			}
			if flag.Short != 0 {
				mf.Short = string(flag.Short)
			}
			sub.Flags = append(sub.Flags, mf)
		}
		if child.Name == "query" {
			data.QueryFlags = sub.Flags
		}
		data.Subcommands = append(data.Subcommands, sub)
	}

	return data
}

func (c *ManCmd) Run(ctx *Globals) error {
	data := extractManData(ctx.Kong)

	date, err := getReproducibleDate()
	if err != nil {
		return fmt.Errorf("error extracting man page data: %w", err)
	}
	data.Date = date

	tmpl, err := template.New("man").Funcs(template.FuncMap{
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	}).Parse(manPageTmpl)
	if err != nil {
		return fmt.Errorf("error parsing man page template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing man page template: %w", err)
	}

	_, err = ctx.Stdout.Write(buf.Bytes())
	return err
}

const manPageTmpl = `.TH {{.Name}} 1 "{{.Date}}" "{{.Version}}" "{{.Name}} Manual"
.SH NAME
{{.Name | toLower}} \- {{.Description}}
.SH SYNOPSIS
.B {{.Name | toLower}}
[\fIflags\fR]
[\fIquery\fR]
.br
.B {{.Name | toLower}}
\fIcommand\fR
[\fIflags\fR]
.SH DESCRIPTION
.B {{.Name | toLower}}
is a command-line tool for querying and modifying data structures.
It supports multiple data formats including JSON, YAML, TOML, CSV, and XML.
.SH COMMANDS
{{range .Subcommands}}.TP
.B {{.Name}}
{{.Help}}
{{end}}
.SH OPTIONS
The following options are available for the default \fBquery\fR command:
{{range .QueryFlags}}.TP
{{if .Short}}\fB\-{{.Short}}\fR, {{end}}\fB\-\-{{.Name}}\fR
{{.Help}}
{{end}}
.SH GLOBAL OPTIONS
{{range .GlobalFlags}}.TP
{{if .Short}}\fB\-{{.Short}}\fR, {{end}}\fB\-\-{{.Name}}\fR
{{.Help}}
{{end}}{{range .Subcommands}}{{if .Flags}}
.SH {{.Name | toUpper}} OPTIONS
{{range .Flags}}.TP
{{if .Short}}\fB\-{{.Short}}\fR, {{end}}\fB\-\-{{.Name}}\fR
{{.Help}}
{{end}}{{end}}{{end}}
.SH EXAMPLES
.TP
Query JSON from stdin:
echo '{"name": "Tom"}' | {{.Name | toLower}} 'name'
.TP
Convert JSON to YAML:
echo '{"name": "Tom"}' | {{.Name | toLower}} -i json -o yaml
.TP
Query with compact output:
echo '{"name": "Tom"}' | {{.Name | toLower}} -i json -o json --compact
.SH SEE ALSO
.UR https://daseldocs.tomwright.me
Dasel documentation
.UE
`
