package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

type Zap struct {
	Cfg Command
}

type Command struct {
	Cmd    string
	Args   []string
	Parser func(stdout []byte) ([]Process, error)
}

type Process struct {
	PID      int
	Name     string
	Port     int
	Protocol string
	Address  string
}

func (z *Zap) ListProcesses() ([]Process, error) {
	if z.Cfg.Cmd == "" {
		return nil, fmt.Errorf("unsupported operating system")
	}
	osCmd := exec.Command(z.Cfg.Cmd, z.Cfg.Args[:]...)
	stdout, err := osCmd.Output()
	if err != nil {
		return nil, err
	}
	return z.Cfg.Parser(stdout)
}

func (z *Zap) SelectProcess() (*Process, error) {
	processes, _ := z.ListProcesses()
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }} (PID: {{ .PID | red }}) :{{ .Port | green }}",
		Inactive: "  {{ .Name | cyan }} (PID: {{ .PID | red }}) :{{ .Port | green }}",
		Selected: "Killing {{ .Name | cyan }} (PID: {{ .PID | red }})...",
		Details: `
--------- Process ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "PID:" | faint }}	{{ .PID }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "Protocol:" | faint }}	{{ .Protocol }}
{{ "Address:" | faint }}	{{ .Address }}`,
	}

	searcher := func(input string, index int) bool {
		process := processes[index]
		name := strings.ToLower(process.Name)
		input = strings.ToLower(input)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select a process to kill",
		Items:     processes,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &processes[index], nil
}

func (z *Zap) KillProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Kill(); err != nil {
		return err
	}
	return nil
}
