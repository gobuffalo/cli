package build

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b", "bill", "install"},
		Short:   "Build the application binary, including bundling of webpack assets",
		RunE:    runE,
	}

	cmd.Flags().StringVarP(&buildOptions.bin, "output", "o", buildOptions.Bin, "set the name of the binary")
	cmd.Flags().StringVarP(&buildOptions.Tags, "tags", "t", "", "compile with specific build tags")
	cmd.Flags().BoolVarP(&buildOptions.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	cmd.Flags().BoolVarP(&buildOptions.SkipAssets, "skip-assets", "k", false, "skip running webpack and building assets")
	cmd.Flags().BoolVarP(&buildOptions.SkipBuildDeps, "skip-build-deps", "", false, "skip building dependencies")
	cmd.Flags().BoolVarP(&buildOptions.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	cmd.Flags().StringVar(&buildOptions.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	cmd.Flags().BoolVarP(&buildOptions.Verbose, "verbose", "v", false, "print debugging information")
	cmd.Flags().BoolVar(&buildOptions.DryRun, "dry-run", false, "runs the build 'dry'")
	cmd.Flags().BoolVar(&buildOptions.SkipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	cmd.Flags().BoolVar(&buildOptions.CleanAssets, "clean-assets", false, "will delete public/assets before calling webpack")
	cmd.Flags().StringVarP(&buildOptions.Environment, "environment", "", "development", "set the environment for the binary")
	cmd.Flags().StringVar(&buildOptions.Mod, "mod", "", "-mod flag for go build")

	return cmd
}
