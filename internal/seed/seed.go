package seed

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/rbac"
)

type (
	Seeder interface {
		RunAll()
	}

	seedBuilder struct {
		client *ent.Client
		rbac   *rbac.Manager
		ctx    context.Context
	}
)

func NewSeeder(i do.Injector) (Seeder, error) {
	return &seedBuilder{
		client: do.MustInvoke[*ent.Client](i),
		rbac:   do.MustInvoke[*rbac.Manager](i),
		ctx:    context.Background(),
	}, nil
}

func (s *seedBuilder) RunAll() {
	s.Role()
}
