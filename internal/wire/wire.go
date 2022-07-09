package wire

import (
	"fmt"
	"os"

	"github.com/liampulles/proverb-gen/internal/adapter"
	"github.com/liampulles/proverb-gen/internal/driver/cli"
	"github.com/liampulles/proverb-gen/internal/usecase"
)

func Run(wd string, args []string) int {
	engine := wire()
	if err := engine.Run(wd, args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return 1
	}
	return 0
}

func wire() *cli.EngineImpl {
	proverbGenerator := usecase.NewProverbGeneratorImpl()
	gateway := adapter.NewGatewayImpl(proverbGenerator)
	return cli.NewEngineImpl(gateway)
}
