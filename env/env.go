package env

import (
	"fmt"

	"github.com/newde36524/goctl2/pkg/env"
	"github.com/spf13/cobra"
)

func write(_ *cobra.Command, args []string) error {
	if len(sliceVarWriteValue) > 0 {
		return env.WriteEnv(sliceVarWriteValue)
	}
	fmt.Println(env.Print(args...))
	return nil
}
