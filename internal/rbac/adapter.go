package rbac

import (
	"context"
	"strings"

	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"github.com/pkg/errors"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/ent/casbinrule"
	"github.com/sembraniteam/setetes/internal/ent/predicate"
)

const (
	rule0 = iota
	rule1
	rule2
	rule3
	rule4
	rule5
)

type (
	Adapter struct {
		client   *ent.Client
		ctx      context.Context
		filtered bool
	}

	Filter struct {
		Ptype []string
		V0    []string
		V1    []string
		V2    []string
		V3    []string
		V4    []string
		V5    []string
	}
)

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

func (a *Adapter) IsFiltered() bool {
	return a.filtered
}

func (a *Adapter) WithTx(fn func(tx *ent.Tx) error) error {
	tx, err := a.client.Tx(a.ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rollback failed: %s", rerr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "commit failed: %v", tx)
	}

	return nil
}

func (a *Adapter) LoadPolicy(m model.Model) error {
	policies, err := a.client.CasbinRule.Query().Order(ent.Asc("id")).All(a.ctx)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err := loadPolicyLine(policy, m); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) LoadFilteredPolicy(m model.Model, filter Filter) error {
	q := a.client.CasbinRule.Query()

	q = applyIn(q, filter.Ptype, casbinrule.PtypeIn)
	q = applyIn(q, filter.V0, casbinrule.V0In)
	q = applyIn(q, filter.V1, casbinrule.V1In)
	q = applyIn(q, filter.V2, casbinrule.V2In)
	q = applyIn(q, filter.V3, casbinrule.V3In)
	q = applyIn(q, filter.V4, casbinrule.V4In)
	q = applyIn(q, filter.V5, casbinrule.V5In)

	lines, err := q.All(a.ctx)
	if err != nil {
		return err
	}

	for _, line := range lines {
		if err := loadPolicyLine(line, m); err != nil {
			return err
		}
	}

	a.filtered = true

	return nil
}

func (a *Adapter) SavePolicy(m model.Model) error {
	return a.WithTx(func(tx *ent.Tx) error {
		if _, err := tx.CasbinRule.Delete().Exec(a.ctx); err != nil {
			return err
		}

		lines := make([]*ent.CasbinRuleCreate, 0)
		for ptype, ast := range m["p"] {
			for _, rule := range ast.Policy {
				line := a.savePolicyLine(tx, ptype, rule)
				lines = append(lines, line)
			}
		}

		for ptype, ast := range m["g"] {
			for _, rule := range ast.Policy {
				line := a.savePolicyLine(tx, ptype, rule)
				lines = append(lines, line)
			}
		}

		batchSize := 5000
		for i := 0; i < len(lines); i += batchSize {
			end := min(i+batchSize, len(lines))

			batch := lines[i:end]
			if _, err := tx.CasbinRule.CreateBulk(batch...).Save(a.ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *Adapter) AddPolicy(sec, ptype string, rules []string) error {
	return a.WithTx(func(tx *ent.Tx) error {
		_, err := a.savePolicyLine(tx, ptype, rules).Save(a.ctx)
		return err
	})
}

func (a *Adapter) AddPolicies(sec, ptypes string, rules [][]string) error {
	return a.WithTx(func(tx *ent.Tx) error {
		return a.createPolicies(tx, ptypes, rules)
	})
}

func (a *Adapter) RemovePolicy(sec, ptype string, rules []string) error {
	return a.WithTx(func(tx *ent.Tx) error {
		instance := a.instance(ptype, rules)
		if _, err := tx.CasbinRule.Delete().Where(
			rulePredicate(instance)...,
		).Exec(a.ctx); err != nil {
			return err
		}

		return nil
	})
}

func (a *Adapter) RemovePolicies(sec, ptype string, rules [][]string) error {
	return a.WithTx(func(tx *ent.Tx) error {
		for _, rule := range rules {
			instance := a.instance(ptype, rule)
			if _, err := tx.CasbinRule.Delete().Where(
				rulePredicate(instance)...,
			).Exec(a.ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *Adapter) RemoveFilteredPolicy(
	sec, ptype string,
	index int,
	values ...string,
) error {
	return a.WithTx(func(tx *ent.Tx) error {
		crs := buildFilteredPredicates(ptype, index, values)
		_, err := tx.CasbinRule.Delete().Where(crs...).Exec(a.ctx)
		return err
	})
}

func (a *Adapter) UpdatePolicy(
	sec, ptype string,
	oldRules, newPolicies []string,
) error {
	return a.WithTx(func(tx *ent.Tx) error {
		rule := a.instance(ptype, oldRules)
		line := tx.CasbinRule.Update().Where(
			rulePredicate(rule)...,
		)

		rule = a.instance(ptype, newPolicies)
		line.SetV0(rule.V0)
		line.SetV1(rule.V1)
		line.SetV2(rule.V2)
		line.SetV3(rule.V3)
		line.SetV4(rule.V4)
		line.SetV5(rule.V5)

		_, err := line.Save(a.ctx)
		return err
	})
}

func (a *Adapter) UpdatePolicies(
	sec, ptype string,
	oldRules, newRules [][]string,
) error {
	return a.WithTx(func(tx *ent.Tx) error {
		for _, policy := range oldRules {
			rule := a.instance(ptype, policy)
			if _, err := tx.CasbinRule.Delete().Where(
				casbinrule.PtypeEQ(rule.Ptype),
				casbinrule.V0EQ(rule.V0),
				casbinrule.V1EQ(rule.V1),
				casbinrule.V2EQ(rule.V2),
				casbinrule.V3EQ(rule.V3),
				casbinrule.V4EQ(rule.V4),
				casbinrule.V5EQ(rule.V5),
			).Exec(a.ctx); err != nil {
				return err
			}
		}

		lines := make([]*ent.CasbinRuleCreate, 0)
		for _, policy := range newRules {
			lines = append(lines, a.savePolicyLine(tx, ptype, policy))
		}

		if _, err := tx.CasbinRule.CreateBulk(lines...).Save(a.ctx); err != nil {
			return err
		}

		return nil
	})
}

func (a *Adapter) UpdateFilteredPolicies(
	sec, ptype string,
	newPolicies [][]string,
	index int,
	values ...string,
) ([][]string, error) {
	oldPolicies := make([][]string, 0)

	err := a.WithTx(func(tx *ent.Tx) error {
		crs := buildFilteredPredicates(ptype, index, values)

		rules, err := tx.CasbinRule.Query().Where(crs...).All(a.ctx)
		if err != nil {
			return err
		}

		if len(rules) == 0 {
			return a.createPolicies(tx, ptype, newPolicies)
		}

		ids := make([]int, 0, len(rules))
		for _, r := range rules {
			ids = append(ids, r.ID)
			oldPolicies = append(oldPolicies, ToStringArray(r))
		}

		if _, err := tx.CasbinRule.Delete().
			Where(casbinrule.IDIn(ids...)).
			Exec(a.ctx); err != nil {
			return err
		}

		return a.createPolicies(tx, ptype, newPolicies)
	})
	if err != nil {
		return nil, err
	}

	return oldPolicies, nil
}

func ToStringArray(rule *ent.CasbinRule) []string {
	arr := make([]string, 0)
	if rule.V0 != "" {
		arr = append(arr, rule.V0)
	}

	if rule.V1 != "" {
		arr = append(arr, rule.V1)
	}

	if rule.V2 != "" {
		arr = append(arr, rule.V2)
	}

	if rule.V3 != "" {
		arr = append(arr, rule.V3)
	}

	if rule.V4 != "" {
		arr = append(arr, rule.V4)
	}

	if rule.V5 != "" {
		arr = append(arr, rule.V5)
	}

	return arr
}

func (a *Adapter) savePolicyLine(
	tx *ent.Tx,
	ptype string,
	rules []string,
) *ent.CasbinRuleCreate {
	line := tx.CasbinRule.Create()

	line.SetPtype(ptype)
	if len(rules) > rule0 {
		line.SetV0(rules[rule0])
	}

	if len(rules) > rule1 {
		line.SetV1(rules[rule1])
	}

	if len(rules) > rule2 {
		line.SetV2(rules[rule2])
	}

	if len(rules) > rule3 {
		line.SetV3(rules[rule3])
	}

	if len(rules) > rule4 {
		line.SetV4(rules[rule4])
	}

	if len(rules) > rule5 {
		line.SetV5(rules[rule5])
	}

	return line
}

func (a *Adapter) instance(ptype string, rules []string) *ent.CasbinRule {
	instance := &ent.CasbinRule{}

	instance.Ptype = ptype
	if len(rules) > rule0 {
		instance.V0 = rules[rule0]
	}

	if len(rules) > rule1 {
		instance.V1 = rules[rule1]
	}

	if len(rules) > rule2 {
		instance.V2 = rules[rule2]
	}

	if len(rules) > rule3 {
		instance.V3 = rules[rule3]
	}

	if len(rules) > rule4 {
		instance.V4 = rules[rule4]
	}

	if len(rules) > rule5 {
		instance.V5 = rules[rule5]
	}

	return instance
}

func (a *Adapter) createPolicies(
	tx *ent.Tx,
	ptype string,
	policies [][]string,
) error {
	lines := make([]*ent.CasbinRuleCreate, 0)
	for _, policy := range policies {
		lines = append(lines, a.savePolicyLine(tx, ptype, policy))
	}

	if _, err := tx.CasbinRule.CreateBulk(lines...).Save(a.ctx); err != nil {
		return err
	}

	return nil
}

func loadPolicyLine(line *ent.CasbinRule, m model.Model) error {
	p := []string{
		line.Ptype,
		line.V0,
		line.V1,
		line.V2,
		line.V3,
		line.V4,
		line.V5,
	}

	var lt string
	if line.V5 != "" {
		lt = strings.Join(p, ", ")
	} else if line.V4 != "" {
		lt = strings.Join(p[:6], ", ")
	} else if line.V3 != "" {
		lt = strings.Join(p[:5], ", ")
	} else if line.V2 != "" {
		lt = strings.Join(p[:4], ", ")
	} else if line.V1 != "" {
		lt = strings.Join(p[:3], ", ")
	} else if line.V0 != "" {
		lt = strings.Join(p[:2], ", ")
	}

	if err := persist.LoadPolicyLine(lt, m); err != nil {
		return err
	}

	return nil
}

func applyIn[T comparable](
	q *ent.CasbinRuleQuery,
	values []T,
	where func(...T) predicate.CasbinRule,
) *ent.CasbinRuleQuery {
	if len(values) > 0 {
		return q.Where(where(values...))
	}

	return q
}

func rulePredicate(instance *ent.CasbinRule) []predicate.CasbinRule {
	return []predicate.CasbinRule{
		casbinrule.PtypeEQ(instance.Ptype),
		casbinrule.V0EQ(instance.V0),
		casbinrule.V1EQ(instance.V1),
		casbinrule.V2EQ(instance.V2),
		casbinrule.V3EQ(instance.V3),
		casbinrule.V4EQ(instance.V4),
		casbinrule.V5EQ(instance.V5),
	}
}

func buildFilteredPredicates(
	ptype string,
	index int,
	values []string,
) []predicate.CasbinRule {
	crs := []predicate.CasbinRule{
		casbinrule.PtypeEQ(ptype),
	}

	cols := []func(string) predicate.CasbinRule{
		casbinrule.V0EQ,
		casbinrule.V1EQ,
		casbinrule.V2EQ,
		casbinrule.V3EQ,
		casbinrule.V4EQ,
		casbinrule.V5EQ,
	}

	for i := range cols {
		pos := i - index
		if pos < 0 || pos >= len(values) {
			continue
		}
		if values[pos] == "" {
			continue
		}
		crs = append(crs, cols[i](values[pos]))
	}

	return crs
}
