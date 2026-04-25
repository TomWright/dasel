package cli

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/tomwright/dasel/v3/parsing"
)

type CompletionCmd struct {
	Shell string `arg:"" enum:"bash,zsh,fish,powershell" help:"Shell type (bash, zsh, fish, powershell)"`
}

type completionFlag struct {
	Name  string
	Short string
	Help  string
	Bool  bool
}

type completionSubcommand struct {
	Name  string
	Help  string
	Flags []completionFlag
}

type completionData struct {
	Name        string
	Subcommands []completionSubcommand
	GlobalFlags []completionFlag
	QueryFlags  []completionFlag
	FormatList  string
	Formats     []string
	ShellList   string
}

func extractCompletionData(k *kong.Kong) completionData {
	app := k.Model

	data := completionData{
		Name:      app.Name,
		ShellList: "bash zsh fish powershell",
	}

	// Collect formats from registered readers, sorted for stable output.
	readers := parsing.RegisteredReaders()
	formats := make([]string, 0, len(readers))
	for _, f := range readers {
		formats = append(formats, f.String())
	}
	sort.Strings(formats)
	data.Formats = formats
	data.FormatList = strings.Join(formats, " ")

	// Global flags from the root node.
	for _, flag := range app.Flags {
		if flag.Hidden {
			continue
		}
		cf := completionFlag{
			Name: flag.Name,
			Help: flag.Help,
			Bool: flag.IsBool(),
		}
		if flag.Short != 0 {
			cf.Short = string(flag.Short)
		}
		data.GlobalFlags = append(data.GlobalFlags, cf)
	}

	// Walk children for subcommands.
	for _, child := range app.Children {
		if child.Hidden {
			continue
		}
		if child.Type != kong.CommandNode {
			continue
		}
		sub := completionSubcommand{
			Name: child.Name,
			Help: child.Help,
		}
		for _, flag := range child.Flags {
			if flag.Hidden {
				continue
			}
			cf := completionFlag{
				Name: flag.Name,
				Help: flag.Help,
				Bool: flag.IsBool(),
			}
			if flag.Short != 0 {
				cf.Short = string(flag.Short)
			}
			sub.Flags = append(sub.Flags, cf)
		}
		if child.Name == "query" {
			data.QueryFlags = sub.Flags
		}
		data.Subcommands = append(data.Subcommands, sub)
	}

	return data
}

func (c *CompletionCmd) Run(ctx *Globals) error {
	data := extractCompletionData(ctx.Kong)

	var tmplStr string
	switch c.Shell {
	case "bash":
		tmplStr = bashCompletionTmpl
	case "zsh":
		tmplStr = zshCompletionTmpl
	case "fish":
		tmplStr = fishCompletionTmpl
	case "powershell":
		tmplStr = powershellCompletionTmpl
	default:
		return fmt.Errorf("unsupported shell: %s", c.Shell)
	}

	tmpl, err := template.New("completion").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("error parsing completion template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing completion template: %w", err)
	}

	_, err = ctx.Stdout.Write(buf.Bytes())
	return err
}

const bashCompletionTmpl = `#!/bin/bash

_{{.Name}}() {
    local cur prev words cword
    _init_completion || return

    local commands="{{range $i, $s := .Subcommands}}{{if $i}} {{end}}{{$s.Name}}{{end}}"
    local formats="{{.FormatList}}"
    local shells="{{.ShellList}}"

    # Find if a subcommand has been specified
    local cmd=""
    for ((i=1; i < cword; i++)); do
        case "${words[i]}" in
            {{range $i, $s := .Subcommands}}{{if $i}}|{{end}}{{$s.Name}}{{end}})
                cmd="${words[i]}"
                break
                ;;
        esac
    done

    # Complete format names for --in/--out/-i/-o
    case "${prev}" in
        --in|--out|-i|-o)
            COMPREPLY=($(compgen -W "${formats}" -- "${cur}"))
            return
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
        local flags=""
        case "${cmd}" in
            "")
                # No subcommand yet: offer global flags + query flags (default command)
                flags="{{range .GlobalFlags}}--{{.Name}} {{end}}{{range .QueryFlags}}--{{.Name}} {{if .Short}}-{{.Short}} {{end}}{{end}}"
                ;;{{range .Subcommands}}
            {{.Name}})
                flags="{{range $.GlobalFlags}}--{{.Name}} {{end}}{{range .Flags}}--{{.Name}} {{if .Short}}-{{.Short}} {{end}}{{end}}"
                ;;{{end}}
        esac
        COMPREPLY=($(compgen -W "${flags}" -- "${cur}"))
        return
    fi

    # If no subcommand yet, complete with subcommand names
    if [[ -z "${cmd}" ]]; then
        COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
    fi
}

complete -F _{{.Name}} {{.Name}}
`

const zshCompletionTmpl = `#compdef {{.Name}}

_{{.Name}}() {
    local -a commands formats
    commands=(
{{- range .Subcommands}}
        '{{.Name}}:{{.Help}}'
{{- end}}
    )
    formats=({{range .Formats}}{{.}} {{end}})

    _arguments -C \
        '1:command:->command' \
        '*::arg:->args'

    case "${state}" in
        command)
            _describe 'command' commands
            _arguments \
{{- range .GlobalFlags}}
                '--{{.Name}}[{{.Help}}]' \
{{- if .Short}}
                '-{{.Short}}[{{.Help}}]' \
{{- end}}
{{- end}}
{{- range .QueryFlags}}
                '--{{.Name}}[{{.Help}}]' \
{{- if .Short}}
                '-{{.Short}}[{{.Help}}]' \
{{- end}}
{{- end}}

            ;;
        args)
            case "${words[1]}" in
{{- range .Subcommands}}
                {{.Name}})
                    _arguments \
{{- range .Flags}}
                        '--{{.Name}}[{{.Help}}]' \
{{- if .Short}}
                        '-{{.Short}}[{{.Help}}]' \
{{- end}}
{{- end}}
{{- range $.GlobalFlags}}
                        '--{{.Name}}[{{.Help}}]' \
{{- end}}

                    ;;
{{- end}}
            esac

            # Complete format names for --in/--out
            if [[ "${words[CURRENT-1]}" == --in || "${words[CURRENT-1]}" == --out || "${words[CURRENT-1]}" == -i || "${words[CURRENT-1]}" == -o ]]; then
                _describe 'format' formats
            fi
            ;;
    esac
}

compdef _{{.Name}} {{.Name}}
`

const fishCompletionTmpl = `# Fish completion for {{.Name}}

# Disable file completions by default
complete -c {{.Name}} -f

# Subcommands
{{range .Subcommands}}complete -c {{$.Name}} -n __fish_use_subcommand -a {{.Name}} -d '{{.Help}}'
{{end}}
# Global flags
{{range .GlobalFlags}}complete -c {{$.Name}} -l {{.Name}}{{if .Short}} -s {{.Short}}{{end}} -d '{{.Help}}'
{{end}}
# Query flags (default command, available at root level)
{{range .QueryFlags}}complete -c {{$.Name}} -n '__fish_use_subcommand' -l {{.Name}}{{if .Short}} -s {{.Short}}{{end}} -d '{{.Help}}'
{{end}}
# Per-subcommand flags
{{range .Subcommands}}{{$sub := .Name}}{{range .Flags}}complete -c {{$.Name}} -n '__fish_seen_subcommand_from {{$sub}}' -l {{.Name}}{{if .Short}} -s {{.Short}}{{end}} -d '{{.Help}}'
{{end}}{{end}}
# Format completions for --in/--out
complete -c {{.Name}} -l in -xa '{{.FormatList}}'
complete -c {{.Name}} -l out -xa '{{.FormatList}}'

# Shell completions for completion command
complete -c {{.Name}} -n '__fish_seen_subcommand_from completion' -xa '{{.ShellList}}'
`

const powershellCompletionTmpl = `Register-ArgumentCompleter -Native -CommandName {{.Name}} -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commands = @({{range .Subcommands}}'{{.Name}}', {{end}}'help')
    $formats = @({{range .Formats}}'{{.}}', {{end}}'')
    $shells = @('bash', 'zsh', 'fish', 'powershell')

    # Find the subcommand
    $subcommand = $null
    $tokens = $commandAst.ToString().Split() | Select-Object -Skip 1
    foreach ($token in $tokens) {
        if ($commands -contains $token) {
            $subcommand = $token
            break
        }
    }

    # Get the previous word for context-sensitive completion
    $previousWord = $null
    $allTokens = $commandAst.ToString().Split()
    if ($allTokens.Count -gt 1) {
        $previousWord = $allTokens[-1]
        if ($wordToComplete -ne '') {
            if ($allTokens.Count -gt 2) {
                $previousWord = $allTokens[-2]
            }
        }
    }

    # Complete format names
    if ($previousWord -in @('--in', '--out', '-i', '-o')) {
        $formats | Where-Object { $_ -ne '' -and $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
        return
    }

    # Complete shell names for completion command
    if ($subcommand -eq 'completion') {
        $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
        return
    }

    if ($wordToComplete -like '-*') {
        $flags = @()
        switch ($subcommand) {
            $null {
                # No subcommand: global + query flags
                $flags = @({{range .GlobalFlags}}'--{{.Name}}', {{end}}{{range .QueryFlags}}'--{{.Name}}', {{if .Short}}'-{{.Short}}', {{end}}{{end}}'')
            }
            {{range .Subcommands}}'{{.Name}}' {
                $flags = @({{range $.GlobalFlags}}'--{{.Name}}', {{end}}{{range .Flags}}'--{{.Name}}', {{if .Short}}'-{{.Short}}', {{end}}{{end}}'')
            }
            {{end}}
        }
        $flags | Where-Object { $_ -ne '' -and $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_)
        }
        return
    }

    # Complete subcommands
    if (-not $subcommand) {
        $commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'Command', $_)
        }
    }
}
`
