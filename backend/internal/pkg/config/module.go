package config

import (
	"context"
	"os"

	nossasdespesas "github.com/Beigelman/nossas-despesas"
	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/internal/pkg/eon"
)

var Module = eon.NewModule("Config", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	di.Provide(c, func() (*nossasdespesas.Config, error) {
		environment := env.MustParse(os.Getenv("ENV"))
		cfg, err := nossasdespesas.NewConfig(environment)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	})
})
