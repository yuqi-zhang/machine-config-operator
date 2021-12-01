package main

import (
	"flag"
	"fmt"
	"os"

	daemon "github.com/openshift/machine-config-operator/pkg/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var firstbootCompleteMachineconfig = &cobra.Command{
	Use:                   "firstboot-complete-machineconfig",
	DisableFlagsInUseLine: true,
	Short:                 "Complete the host's initial boot into a MachineConfig",
	Args:                  cobra.MaximumNArgs(0),
	Run:                   executeFirstbootCompleteMachineConfig,
}

var layeredImage bool

// init executes upon import
func init() {
	rootCmd.AddCommand(firstbootCompleteMachineconfig)
	startCmd.PersistentFlags().BoolVarP(&layeredImage, "layered", "", false, "Bootstrap with layered image")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func runFirstBootCompleteMachineConfig(_ *cobra.Command, _ []string) error {
	flag.Set("logtostderr", "true")
	flag.Parse()

	exitCh := make(chan error)
	defer close(exitCh)

	dn, err := daemon.New(daemon.NewNodeUpdaterClient(), exitCh)
	if err != nil {
		return err
	}

	if layeredImage {
		return dn.RunFirstbootCompleteMachineconfigLayered()
	}
	return dn.RunFirstbootCompleteMachineconfig()
}

func executeFirstbootCompleteMachineConfig(cmd *cobra.Command, args []string) {
	err := runFirstBootCompleteMachineConfig(cmd, args)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
