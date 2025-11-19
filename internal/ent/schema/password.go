package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Password holds the schema definition for the Password entity.
type Password struct {
	ent.Schema
}

// Mixin for the Password.
func (Password) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Password.
func (Password) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash").Unique().Sensitive().Comment("Hashed password using Argon2.").
			Annotations(entsql.WithComments(true)),
	}
}

// Edges of the Password.
func (Password) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("password").Unique().Required(),
	}
}
