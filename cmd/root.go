package cmd

import (
	"errors"
	"fmt"
	"os"

	"go-shc/shc"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "go-shc",
	Short:            "Obfuscate shell scripts",
	TraverseChildren: true,
	Version:          shc.Version,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		script, _ := cmd.Flags().GetString("script")
		if script != "" && len(args) > 0 {
			return errors.New("input files have been specified, but with a script parameter")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		opt := shc.ObfuscateOption{}
		opt.Shell, _ = cmd.Flags().GetString("shell")
		opt.UseTempFile, _ = cmd.Flags().GetBool("use-temp-file")
		opt.VarNamePrefix, _ = cmd.Flags().GetString("var-prefix")
		opt.SliceLength, _ = cmd.Flags().GetInt("slice-length")
		opt.MinVarNameLength, _ = cmd.Flags().GetInt("min-var-name-length")
		opt.ArgLengthLimit, _ = cmd.Flags().GetInt("arg-length-limit")
		opt.VarCountLimit, _ = cmd.Flags().GetInt("var-count-limit")
		script, _ := cmd.Flags().GetString("script")
		output, _ := cmd.Flags().GetString("output")
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
		result := shc.ObfuscateShellScript(script, opt)
		if output != "" {
			if err := os.WriteFile(output, []byte(result.Output), 0644); err == nil {
				fmt.Println("obfuscated shell script has been written to", output)
			} else {
				exitWithError(err)
			}
		} else {
			fmt.Println(result.Output)
		}
	},
}

func init() {
	rootCmd.Flags().Bool("glob", false, "use glob pattern")
	rootCmd.Flags().StringP("script", "s", "", "Specify the input script")
	rootCmd.Flags().StringP("output", "o", "", "Specify the output file")
	rootCmd.Flags().String("shell", "", "Specify the shell")
	rootCmd.Flags().BoolP("use-temp-file", "t", false, "Run by creating temporary files")
	rootCmd.Flags().StringP("var-prefix", "p", "_", "Specify the variable prefix")
	rootCmd.Flags().Int("slice-length", 5, "Specify the slice length")
	rootCmd.Flags().Int("min-var-name-length", 3, "Specify the min variable name length")
	rootCmd.Flags().Int("arg-length-limit", 17200, "Specify the command line argument length limit")
	rootCmd.Flags().Int("var-count-limit", 20000, "Specify the variable count limit")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
