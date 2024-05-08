package main

import (
	"os"

	"github.com/Peltoche/halium/cmd/halium/commands"
	"github.com/Peltoche/halium/internal/tools/buildinfos"
	"github.com/spf13/cobra"
)

const binaryName = "halium"

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	cmd := &cobra.Command{
		Use:     binaryName,
		Short:   "Manage your halium instance in your terminal.",
		Version: buildinfos.Version(),
	}

	// Subcommands
	cmd.AddCommand(commands.NewRunCmd(binaryName))

	err := cmd.Execute()
	if err != nil {
		return exitError
	}

	return exitOK
}
