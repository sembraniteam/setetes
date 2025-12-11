package rbac

import (
	_ "embed"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	entadapter "github.com/casbin/ent-adapter"
	"github.com/casbin/ent-adapter/ent"
)

//go:embed model.conf
var modelFile string

type Manager struct {
	enforcer *casbin.SyncedEnforcer
}

func New(client *ent.Client) (*Manager, error) {
	a, err := entadapter.NewAdapterWithClient(client)
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromString(modelFile)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Manager{enforcer: enforcer}, nil
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
