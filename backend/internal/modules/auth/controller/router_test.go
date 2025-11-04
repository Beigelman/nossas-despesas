package controller_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth/controller"
)

func TestRouter(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	var called []string
	h := func(name string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			called = append(called, name)
			return nil
		}
	}

	controller.Router(app,
		h("signIn"),
		h("signUp"),
		h("google"),
		h("refresh"),
	)

	routes := app.GetRoutes()
	paths := make([]string, len(routes))
	for i, r := range routes {
		paths[i] = r.Method + " " + r.Path
	}

	assert.Contains(t, paths, "POST /api/v1/auth/sign-in/credentials")
	assert.Contains(t, paths, "POST /api/v1/auth/sign-in/google")
	assert.Contains(t, paths, "POST /api/v1/auth/sign-up/credentials")
	assert.Contains(t, paths, "POST /api/v1/auth/refresh-token")
}
