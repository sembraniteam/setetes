package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Mixin for the Role.
func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MinLen(3).
			MaxLen(164).
			StructTag(`json:"name"`),
		field.String("key").
			NotEmpty().
			MinLen(3).
			MaxLen(164).
			StructTag(`json:"key"`),
		field.String("domain").
			NotEmpty().
			MinLen(1).
			MaxLen(164).
			StructTag(`json:"domain"`),
		field.String("description").
			Optional().
			MinLen(30).
			MaxLen(300).
			StructTag(`json:"description"`),
		field.Bool("activated").Default(false).StructTag(`json:"activated"`),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("accounts", Account.Type).Ref("role"),
		edge.To("parent", Role.Type).From("children"),
	}
}

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		&entsql.Annotation{
			Checks: map[string]string{
				"name":        "length(name) >= 3 and length(name) <= 164",
				"key":         "length(key) >= 3 and length(key) <= 164",
				"domain":      "length(domain) >= 1 and length(domain) <= 164",
				"description": "length(description) >= 30 and length(description) <= 300",
			},
		},
	}
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("domain", "key").Unique(),
	}
}
