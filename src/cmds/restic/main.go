package main

import (
	"fmt"
	"os"
	"restic"
	"restic/debug"
	"restic/options"

	"github.com/spf13/cobra"

	"restic/errors"
)

// cmdRoot is the base command when no other command has been specified.
var cmdRoot = &cobra.Command{
	Use:   "restic",
	Short: "backup and restore files",
	Long: `
restic is a backup program which allows saving multiple revisions of files and
directories in an encrypted repository stored on different backends.
`,
	SilenceErrors: true,
	SilenceUsage:  true,

	PersistentPreRunE: func(*cobra.Command, []string) error {
		// parse extended options
		opts, err := options.Parse(globalOptions.Options)
		if err != nil {
			return err
		}
		globalOptions.extended = opts

		// run the debug functions for all subcommands (if build tag "debug" is
		// enabled)
		if err := runDebug(); err != nil {
			return err
		}

		return nil
	},
	PersistentPostRun: func(*cobra.Command, []string) {
		shutdownDebug()
	},
}

func main() {
	debug.Log("main %#v", os.Args)
	err := cmdRoot.Execute()

	switch {
	case restic.IsAlreadyLocked(errors.Cause(err)):
		fmt.Fprintf(os.Stderr, "%v\nthe `unlock` command can be used to remove stale locks\n", err)
	case errors.IsFatal(errors.Cause(err)):
		fmt.Fprintf(os.Stderr, "%v\n", err)
	case err != nil:
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	var exitCode int
	if err != nil {
		exitCode = 1
	}

	Exit(exitCode)
}
