package rbac

import (
	"context"

	"github.com/megalodev/setetes/internal/ent"
)

type Adapter struct {
	client *ent.Client
	ctx    context.Context
}

func NewAdapter(client *ent.Client) (*Adapter, error) {
	a := &Adapter{
		client: client,
		ctx:    context.Background(),
	}

	if err := client.Schema.Create(a.ctx); err != nil {
		return nil, err
	}

	return a, nil
}
