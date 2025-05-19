package bug

import (
	"github.com/newde36524/goctl2/internal/cobrax"
	"github.com/spf13/cobra"
)

// Cmd describes a bug command.
var Cmd = cobrax.NewCommand("bug", cobrax.WithRunE(runE), cobrax.WithArgs(cobra.NoArgs))
