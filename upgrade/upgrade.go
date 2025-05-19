package upgrade

import (
	"fmt"
	"runtime"

	"github.com/newde36524/goctl2/rpc/execx"
	"github.com/spf13/cobra"
)

// upgrade gets the latest goctl by
// go install github.com/newde36524/goctl2@latest
func upgrade(_ *cobra.Command, _ []string) error {
	cmd := `go install github.com/newde36524/goctl2@latest`
	if runtime.GOOS == "windows" {
		cmd = `go install github.com/newde36524/goctl2@latest`
	}
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
