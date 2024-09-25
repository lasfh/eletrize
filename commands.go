package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/schema"
	"github.com/lasfh/eletrize/watcher"
)

var version = "unknown"

func execute() error {
	var schema uint16

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
				eletrize, err = NewEletrizeByFileInCW()
				if err != nil {
					fmt.Printf("eletrize: %s\n", err.Error())
					os.Exit(1)
				}
			} else {
				eletrize, err = NewEletrize(args[0])
				if err != nil {
					fmt.Printf("eletrize: %s\n", err.Error())
					os.Exit(1)
				}
			}

			if schema > 0 {
				if err = eletrize.StartFromSchema(schema); err != nil {
					fmt.Printf("eletrize: %s\n", err.Error())
					os.Exit(1)
				}

				os.Exit(0)
			}

			eletrize.Start()
		},
	}

	rootCmd.Flags().Uint16VarP(&schema, "schema", "s", 0, "Execute a specific schema")
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
						Label:   output.Label(label),
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

			eletrize.Start()
		},
	}

	cmd.PersistentFlags().StringVarP(&label, "label", "l", "APP", "Set the identification label")
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
			fmt.Printf("Eletrize, version: %s (%s)\n", version, runtime.Version())

			if !debugInfo {
				return
			}

			if info, ok := debug.ReadBuildInfo(); ok {
				fmt.Println("\nDebug info:", info.Main.Version)
				fmt.Println(info)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&debugInfo, "info", "i", false, "Debugging information")

	return cmd
}
