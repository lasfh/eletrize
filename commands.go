package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/schema"
	"github.com/lasfh/eletrize/watcher"
	"github.com/spf13/cobra"
)

var version string

func execute() error {
	var schema []uint

	rootCmd := &cobra.Command{
		Use:   "eletrize [filename]",
		Short: "Live reload tool for Go and generic projects",
		Long: `Eletrize is a live reload utility designed for Go projects and generic applications. 
		It monitors changes in the specified directory and automatically triggers a reload, allowing for a dynamic and efficient development workflow.
		Specify the [filename] argument to define the configuration file for Eletrize.`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				eletrize *Eletrize
				err      error
			)

			if len(args) == 0 {
				eletrize, err = NewEletrizeFromWD()
			} else {
				eletrize, err = NewEletrizeFromPath(args[0])
			}

			if err != nil {
				fmt.Printf("eletrize: %s\n", err.Error())
				os.Exit(1)
			}

			if len(schema) > 0 {
				eletrize.Start(args, schema...)
				os.Exit(0)
			}

			if eletrize.launch {
				fmt.Printf("eletrize: using default schema 1\n\n")
				fmt.Println(".vscode/launch.json:")

				for index, schema := range eletrize.Schema {
					fmt.Printf("\t %d -> %s\n", index+1, schema.Label.Label)
				}

				fmt.Println("\nto use a specific schema:")
				fmt.Printf("\teletrize --schema N\n\n")

				eletrize.StartOne()
				os.Exit(0)
			}

			eletrize.Start(args)
		},
	}

	rootCmd.Flags().UintSliceVarP(&schema, "schema", "s", []uint{}, "Execute a specific schema")
	rootCmd.AddCommand(
		runCommand(),
		versionCommand(),
	)

	return rootCmd.Execute()
}

func runCommand() *cobra.Command {
	var (
		label      string
		path       string
		recursive  bool
		extensions []string
		envFile    string
		workdir    string
	)

	cmd := &cobra.Command{
		Use:   "run [run] [build]",
		Short: "Run a simple execution and/or compilation command",
		Long: `The “run” command allows you to execute a command directly without the need for a configuration file.
		You can include optional [run] and [build] arguments to customize the build process and specify the command to run.`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				build *command.Command
				run   command.Command
			)

			if len(args) == 2 {
				// build command
				parts := strings.Fields(args[1])

				build = &command.Command{
					Method: strings.TrimSpace(
						parts[0],
					),
					Args: []string{},
				}

				if len(parts) > 1 {
					build.Args = parts[1:]
				}
			}

			// run command
			parts := strings.Fields(args[0])

			run.Method = parts[0]

			if len(parts) > 1 {
				run.Args = parts[1:]
			}

			eletrize := &Eletrize{
				Schema: []schema.Schema{
					{
						Label: &output.Label{
							Label: label,
						},
						Workdir: workdir,
						EnvFile: envFile,
						Watcher: watcher.Options{
							Path:       path,
							Recursive:  recursive,
							Extensions: extensions,
						},
						Commands: command.Commands{
							Build: build,
							Run:   []command.Command{run},
						},
					},
				},
			}

			eletrize.StartOne()
		},
	}

	cmd.PersistentFlags().StringVarP(&label, "label", "l", "", "Set the identification label")
	cmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "Set the path to watch for changes")
	cmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Enable recursive mode for watching")
	cmd.PersistentFlags().StringSliceVarP(&extensions, "ext", "e", []string{}, "Set file extensions to watch")
	cmd.PersistentFlags().StringVarP(&envFile, "env", "", "", "Set the path to the environment file")
	cmd.PersistentFlags().StringVarP(&workdir, "workdir", "", "", "Sets the working directory")

	return cmd
}

func versionCommand() *cobra.Command {
	var debugInfo bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version and debug information.",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := debug.ReadBuildInfo()

			fmt.Printf("Eletrize, version: %s (%s)\n", getVersion(info), runtime.Version())

			if debugInfo && ok {
				fmt.Println("\nDebug info:")
				fmt.Println(info)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&debugInfo, "info", "i", false, "Debugging information")

	return cmd
}

func getVersion(info *debug.BuildInfo) string {
	if version == "" {
		if info != nil {
			return info.Main.Version
		}

		return "(unknown)"
	}

	return version
}
