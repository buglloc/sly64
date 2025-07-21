package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/buglloc/sly64/v2/internal/config"
)

var runtime *config.Runtime
var rootArgs struct {
	Configs []string
}

var rootCmd = &cobra.Command{
	Use:           "sly64",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         `lightweight and flexible DNS server`,
}

func init() {
	cobra.OnInitialize(
		initRuntime,
		initLogger,
	)

	flags := rootCmd.PersistentFlags()
	flags.StringSliceVar(&rootArgs.Configs, "config", nil, "config file")

	rootCmd.AddCommand(
		startCmd,
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func initRuntime() {
	var err error
	runtime, err = config.NewRuntime(rootArgs.Configs...)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to create runtime: %v\n", err)
		os.Exit(1)
	}
}

func initLogger() {
	log.SetOutput(os.Stderr)
}
