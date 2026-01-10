package rbac

import (
	_ "embed"
	"errors"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/util"
	"github.com/sembraniteam/setetes/internal/ent"
)

const (
	args2Len = 2
	args4Len = 4
)

//go:embed model.conf
var modelFile string

type Manager struct {
	enforcer *casbin.SyncedEnforcer
}

func New(client *ent.Client) (*Manager, error) {
	a, err := NewAdapter(client)
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromString(modelFile)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewSyncedEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	e.AddFunction("abacMatch", abacMatch)
	e.AddFunction("domMatch", domainMatch)

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Manager{enforcer: e}, nil
}

func (m *Manager) GetEnforcer() casbin.IEnforcer {
	return m.enforcer
}

func (m *Manager) LoadPolicy() error {
	return m.enforcer.LoadPolicy()
}

func (m *Manager) HasRole(user, role string, domain ...string) (bool, error) {
	has, err := m.enforcer.HasRoleForUser(user, role, domain...)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (m *Manager) AddRoleForUser(user, role string, domain ...string) error {
	_, err := m.enforcer.AddRoleForUser(user, role, domain...)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func (m *Manager) RemoveRoleForUser(user, role string, domain ...string) error {
	_, err := m.enforcer.DeleteRoleForUser(user, role, domain...)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func (m *Manager) AddPolicy(role, domain, resource, action string) error {
	_, err := m.enforcer.AddPolicy(role, domain, resource, action)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func (m *Manager) RemovePolicy(role, domain, resource, action string) error {
	_, err := m.enforcer.RemovePolicy(role, domain, resource, action)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func abacMatch(args ...any) (any, error) {
	if len(args) != args4Len {
		return false, errors.New("abacMatch expects 4 arguments")
	}

	return true, nil
}

func domainMatch(args ...any) (any, error) {
	if len(args) != args2Len {
		return false, errors.New("domainMatch expects 2 arguments")
	}

	domain, ok := args[0].(string)
	if !ok {
		return false, errors.New("domainMatch expects string request")
	}

	policy, ok := args[1].(string)
	if !ok {
		return false, errors.New("domainMatch expects string policy")
	}

	if policy == "" || policy == "*" {
		return true, nil
	}

	return util.KeyMatch2(domain, policy), nil
}
