package cli

import (
	"github.com/i4erkasov/go-ddd-cqrs/internal/infrastructure/api/rest"
	"github.com/i4erkasov/go-ddd-cqrs/pkg/zap"
	"github.com/spf13/cobra"
	zapLogger "go.uber.org/zap"
)

const HttpServerCommand = "http-server"
const VersionHttpServer = "1.0.0"

var httpServer = &cobra.Command{
	Use:     HttpServerCommand,
	Short:   "Start http server",
	Version: VersionHttpServer,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cfg = cfg.Sub("app")

		var log *zapLogger.Logger
		if log, err = zap.New(cfg); err != nil {
			return err
		}

		var server *rest.Server
		if server, err = rest.New(cfg.Sub("api.rest"), log); err != nil {
			return err
		}

		return server.Start(cmd.Context())
	},
}

func init() {
	cmd.AddCommand(httpServer)
}
