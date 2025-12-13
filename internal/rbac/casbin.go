package rbac

import (
	_ "embed"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/megalodev/setetes/internal/ent"
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

func (m *Manager) HasRole(user, role string) (bool, error) {
	has, err := m.enforcer.HasRoleForUser(user, role)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (m *Manager) AddRoleForUser(user, role string) error {
	_, err := m.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func (m *Manager) RemoveRoleForUser(user, role string) error {
	_, err := m.enforcer.DeletePermissionForUser(user, role)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func (m *Manager) AddPolicy(role, resource, action string) error {
	_, err := m.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}

func (m *Manager) RemovePolicy(role, resource, action string) error {
	_, err := m.enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		return err
	}

	return m.enforcer.SavePolicy()
}
