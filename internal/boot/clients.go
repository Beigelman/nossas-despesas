package boot

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/service"
	"github.com/Beigelman/ludaapi/internal/infra/jwt"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/eon"
	"google.golang.org/api/option"
	"log/slog"
)

var ClientsModule = eon.NewModule("Clients", func(ctx context.Context, c *di.Container, lc eon.LifeCycleManager, info eon.Info) {
	var dbClient db.Database

	di.Provide(c, func() service.TokenProvider {
		return jwt.NewJWTProvider("-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQDdxJVc0eeqGD/3M2k1ePUT/6xGhFa7s1ef5iwup6n25VOTVxXL\npqNBcD1DWVlkpo25N5D9WSbUowXCmFW04inNc3mmP9rRnp3COWbHeJuOLP9YGKgc\n5uRk6sFknEiQ16mmYzuec3BW6ORYqVOPbUS+BYbWhWkRtgyXXNXM+alvqwIDAQAB\nAoGBAMrZR+4RKiBR8iCBbBi3PSU/1iriXhtunhXqijtarYLinSHGpG8VS3tN2RvD\nnJsOJdBnXT3/0B7rxxcKFEtSG/zW5GQ3BJPRmDYFoglitaghdYHOggY4tPbtoTr6\nM3P9kPFuC2XygL9k6PLJy01aXyduJDqCLqjr1EmfXVLq1m+RAkEA+cnYO2KcNM8+\n7B6bu+IZRLIrlHDDeyNmUJrWViQ+jKZVwCX4nzVg8uVScgGM5m4VztZd2kkaj0Bj\n+UvnYP6h3QJBAONIXBdgQYyjCD5X6ogolr3/1F1SNStxh8tE814xCciwqMayc9Av\n/GtdjI38VqprERMicKEtmcmeVnoQ2jGC8ycCQEUrY7luIRtumFoCT9XDUoP3YqIE\nZ91dfCOt/NR1zOxd0zkWSrarrWEVp7LyQvY8XcWdDvg3bidlCUorfrMZT/ECQCTV\nbD9JlTXykfpwiwzH7y4ZkNQS55UD0CsMIJjKP7irkJ6q+wPpUvIfdhDorS7vLRQ5\nx6EHX94B8CfWJVZSz48CQQDyjOngfBrwLjGaUfpi94bG4neLsyY72F0eVPU8K/qD\n0WmhawCe7HMaqJx2oj5t7tQDxRq43GzfCQyECr9iUmv4\n-----END RSA PRIVATE KEY-----")
	})
	di.Provide(c, func(cfg *config.Config) db.Database {
		dbClient = db.New(cfg)
		return dbClient
	})
	di.Provide(c, func() (*auth.Client, error) {
		opt := option.WithCredentialsFile("./firebaseServiceAccount.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			return nil, fmt.Errorf("firebase.NewApp: %w", err)
		}
		return app.Auth(ctx)
	})

	lc.OnDisposing(eon.HookOrders.PREPEND, func() error {
		if dbClient != nil {
			slog.Info("Closing db connection")
			if err := dbClient.Close(); err != nil {
				return fmt.Errorf("dbClient.Close: %w", err)
			}
		}
		return nil
	})
})
