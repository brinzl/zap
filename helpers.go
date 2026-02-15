package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func FindByPort(processes []Process, port int) *Process {
	for _, p := range processes {
		if p.Port == port {
			return &p
		}
	}
	return nil
}

func ConfirmKill(process *Process) bool {
	label := fmt.Sprintf("Kill %s (PID: %d) on port %d", process.Name, process.PID, process.Port)
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err == nil
}
