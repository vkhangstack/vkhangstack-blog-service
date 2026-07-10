package casbin

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/uptrace/bun"
)

// casbinRule mirrors the casbin_rules table.
type casbinRule struct {
	bun.BaseModel `bun:"table:casbin_rules,alias:cr"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Ptype         string `bun:"ptype,notnull"`
	V0            string `bun:"v0,notnull"`
	V1            string `bun:"v1,notnull"`
	V2            string `bun:"v2,notnull"`
	V3            string `bun:"v3,notnull"`
	V4            string `bun:"v4,notnull"`
	V5            string `bun:"v5,notnull"`
}

// DBAdapter implements persist.Adapter backed by PostgreSQL via bun.
type DBAdapter struct {
	db *bun.DB
}

// NewDBAdapter returns a Casbin adapter that reads/writes casbin_rules via bun.
func NewDBAdapter(db *bun.DB) *DBAdapter {
	return &DBAdapter{db: db}
}

// LoadPolicy loads all rules from the DB into the Casbin model.
func (a *DBAdapter) LoadPolicy(m model.Model) error {
	var rules []casbinRule
	if err := a.db.NewSelect().Model(&rules).Scan(context.Background()); err != nil {
		return fmt.Errorf("casbin db adapter: load: %w", err)
	}
	for _, rule := range rules {
		persist.LoadPolicyLine(
			fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", rule.Ptype, rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5),
			m,
		)
	}
	return nil
}

// SavePolicy writes all model policies back to DB (replaces all rows).
func (a *DBAdapter) SavePolicy(m model.Model) error {
	ctx := context.Background()
	var rules []casbinRule
	// Iterate "p" (policy) and "g" (grouping) sections
	for sec, ptypeMap := range map[string]string{"p": "p", "g": "g"} {
		for _, assertion := range m[sec] {
			for _, policy := range assertion.Policy {
				rules = append(rules, ruleLine(ptypeMap, policy))
			}
		}
	}
	_, err := a.db.NewDelete().Model((*casbinRule)(nil)).Where("1=1").Exec(ctx)
	if err != nil {
		return err
	}
	if len(rules) > 0 {
		_, err = a.db.NewInsert().Model(&rules).Exec(ctx)
	}
	return err
}

// AddPolicy inserts a single policy rule.
func (a *DBAdapter) AddPolicy(sec, ptype string, rule []string) error {
	r := ruleLine(ptype, rule)
	_, err := a.db.NewInsert().Model(&r).On("CONFLICT DO NOTHING").Exec(context.Background())
	return err
}

// RemovePolicy deletes a single policy rule.
func (a *DBAdapter) RemovePolicy(sec, ptype string, rule []string) error {
	r := ruleLine(ptype, rule)
	_, err := a.db.NewDelete().Model(&r).
		Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", r.Ptype, r.V0, r.V1, r.V2).
		Exec(context.Background())
	return err
}

// RemoveFilteredPolicy deletes policy rules matching a field filter.
func (a *DBAdapter) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	ctx := context.Background()
	q := a.db.NewDelete().Model((*casbinRule)(nil)).Where("ptype = ?", ptype)
	columns := []string{"v0", "v1", "v2", "v3", "v4", "v5"}
	for i, v := range fieldValues {
		if v != "" {
			q = q.Where(fmt.Sprintf("%s = ?", columns[fieldIndex+i]), v)
		}
	}
	_, err := q.Exec(ctx)
	return err
}

// ruleLine converts a ptype + values slice into a casbinRule row.
func ruleLine(ptype string, values []string) casbinRule {
	r := casbinRule{Ptype: ptype}
	cols := []*string{&r.V0, &r.V1, &r.V2, &r.V3, &r.V4, &r.V5}
	for i, v := range values {
		if i >= len(cols) {
			break
		}
		*cols[i] = v
	}
	return r
}
