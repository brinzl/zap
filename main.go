package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var osCommands = map[string]Command{
	"linux": {
		Cmd:    "ss",
		Args:   []string{"-tulnp"},
		Parser: ParseSSCmd,
	},
	"darwin": {
		Cmd:    "lsof",
		Args:   []string{"-i", "-P", "-n", "-sTCP:LISTEN"},
		Parser: ParseLsofCmd,
	},
	"windows": {
		Cmd:    "netstat",
		Args:   []string{"-ano"},
		Parser: ParseNetstatCmd,
	},
}

var rootCmd = &cobra.Command{
	Use:     "zap",
	Short:   "A cross-platform zombie port killer",
	Long:    "zap helps you find and kill processes bound to network ports...",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		force, _ := cmd.Flags().GetBool("force")

		z := Zap{
			Cfg: osCommands[runtime.GOOS],
		}

		processes, err := z.ListProcesses()
		if err != nil {
			fmt.Println("Error listing processes:", err)
			os.Exit(1)
		}

		if len(processes) == 0 {
			fmt.Println("No listening processes found")
			return
		}

		if port > 0 {
			process := FindByPort(processes, port)
			if process == nil {
				fmt.Printf("No process found on port %d\n", port)
				os.Exit(1)
			}

			if !force {
				if !ConfirmKill(process) {
					fmt.Println("Aborted")
					return
				}
			}

			if err := z.KillProcess(process.PID); err != nil {
				fmt.Printf("Failed to kill process %d: %v\n", process.PID, err)
				os.Exit(1)
			}

			fmt.Printf("Killed %s (PID: %d) on port %d\n", process.Name, process.PID, process.Port)
			return
		}

		selected, err := z.SelectProcess()
		if err != nil {
			fmt.Println("\nAborted")
			return
		}

		if !force {
			if !ConfirmKill(selected) {
				fmt.Println("Aborted")
				return
			}
		}

		if err := z.KillProcess(selected.PID); err != nil {
			fmt.Printf("Failed to kill process %d: %v\n", selected.PID, err)
			os.Exit(1)
		}

		fmt.Printf("Killed %s (PID: %d) on port %d\n", selected.Name, selected.PID, selected.Port)
	},
}

func main() {
	rootCmd.Flags().IntP("port", "p", 0, "Kill process by port")
	rootCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate(fmt.Sprintf("zap %s (commit: %s, built: %s)\n", version, commit, date))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
