package build

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	xbuildCmd.Flags().StringVarP(&buildOptions.bin, "output", "o", buildOptions.Bin, "set the name of the binary")
	xbuildCmd.Flags().StringVarP(&buildOptions.Tags, "tags", "t", "", "compile with specific build tags")
	xbuildCmd.Flags().BoolVarP(&buildOptions.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	xbuildCmd.Flags().BoolVarP(&buildOptions.SkipAssets, "skip-assets", "k", false, "skip running webpack and building assets")
	xbuildCmd.Flags().BoolVarP(&buildOptions.SkipBuildDeps, "skip-build-deps", "", false, "skip building dependencies")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	xbuildCmd.Flags().StringVar(&buildOptions.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Verbose, "verbose", "v", false, "print debugging information")
	xbuildCmd.Flags().BoolVar(&buildOptions.DryRun, "dry-run", false, "runs the build 'dry'")
	xbuildCmd.Flags().BoolVar(&buildOptions.SkipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	xbuildCmd.Flags().BoolVar(&buildOptions.CleanAssets, "clean-assets", false, "will delete public/assets before calling webpack")
	xbuildCmd.Flags().StringVarP(&buildOptions.Environment, "environment", "", "development", "set the environment for the binary")
	xbuildCmd.Flags().StringVar(&buildOptions.Mod, "mod", "", "-mod flag for go build")

	return xbuildCmd
}
