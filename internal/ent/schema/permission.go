package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Mixin for the Permission.
func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
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
		field.String("resource").NotEmpty().StructTag(`json:"resource"`),
		field.String("action").NotEmpty().StructTag(`json:"action"`),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).
			Ref("permissions").
			Through("role_permissions", RolePermission.Type),
	}
}

func (Permission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		&entsql.Annotation{
			Checks: map[string]string{
				"name":   "length(name) >= 3 and length(name) <= 164",
				"key":    "length(key) >= 3 and length(key) <= 164",
				"domain": "length(domain) >= 1 and length(domain) <= 164",
			},
		},
	}
}

func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("domain", "resource", "action"),
		index.Fields("domain", "key").Unique(),
	}
}
