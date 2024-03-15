package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"go-shc/shc"
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack shell scripts",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		script, _ := cmd.Flags().GetString("script")
		if script != "" && len(args) > 0 {
			return errors.New("input files have been specified, but with a script parameter")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		opt := shc.PackOption{}
		opt.Shell, _ = cmd.Flags().GetString("shell")
		opt.Osarch, _ = cmd.Flags().GetString("osarch")
		opt.Output, _ = cmd.Flags().GetString("output")
		opt.OutputDir, _ = cmd.Flags().GetString("output-dir")
		opt.TempDir, _ = cmd.Flags().GetString("temp-dir")
		opt.TrimPath, _ = cmd.Flags().GetBool("trim-path")
		opt.UseTempFile, _ = cmd.Flags().GetBool("use-temp-file")
		opt.GoCompiler, _ = cmd.Flags().GetString("compiler")
		script, _ := cmd.Flags().GetString("script")
		glob, _ := cmd.Flags().GetBool("glob")
		if len(args) > 0 {
			if glob {
				args = getGlobFiles(args...)
			}
			var buf []byte
			for _, v := range args {
				bs, err := os.ReadFile(v)
				exitWithError(err)
				buf = append(buf, bs...)
			}
			script = string(buf)
		}

		onlySource, _ := cmd.Flags().GetBool("only-source")
		if onlySource {
			if opt.Shell == "" {
				_, shell, index := shc.ParseShebang(script)
				if index > 0 {
					opt.Shell = shell
				} else {
					opt.Shell = "bash"
				}
			}
			var r io.Writer = os.Stdout
			if opt.Output != "" {
				f, err := os.OpenFile(opt.Output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
				exitWithError(err)
				r = f
				defer f.Close()
			}
			exitWithError(shc.GenerateGoCodes(script, r, opt.GenerateOption))
			if opt.Output != "" {
				fmt.Println("go source code has been written to", opt.Output)
			}
			return
		}
		if f, err := shc.PackShellScript(script, opt); err == nil {
			f.Close()
			fmt.Println(f.Name())
		} else {
			exitWithError(err)
		}
	},
}

func init() {
	packCmd.Flags().Bool("glob", false, "use glob pattern")
	packCmd.Flags().Bool("only-source", false, "Only generate source code and write to stdout or output file")
	packCmd.Flags().StringP("script", "s", "", "Specify the input script")
	packCmd.Flags().String("shell", "", "Specify the shell")
	packCmd.Flags().String("osarch", "", "Specify the os and arch, e.g. linux/amd64")
	packCmd.Flags().StringP("output", "o", "", "Specify the output file")
	packCmd.Flags().String("output-dir", "./", "Specify the output directory")
	packCmd.Flags().String("temp-dir", "", "Specify the temporary directory of the generated source code")
	packCmd.Flags().Bool("trim-path", false, "Trim source code path")
	packCmd.Flags().BoolP("use-temp-file", "t", false, "Run by creating temporary files")
	packCmd.Flags().String("compiler", "go", "Specify the compiler")
	rootCmd.AddCommand(packCmd)
}
