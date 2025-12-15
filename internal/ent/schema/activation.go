package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Activation holds the schema definition for the Activation entity.
type Activation struct {
	ent.Schema
}

// Mixin for the Activation.
func (Activation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Activation.
func (Activation) Fields() []ent.Field {
	return []ent.Field{
		field.String("token").
			MaxLen(44).
			Sensitive().
			NotEmpty().
			Unique().
			SchemaType(map[string]string{dialect.Postgres: "char(44)"}).
			Comment("The activation token is single-use and will be deleted after it is used."),
		field.Bool("is_used").Default(false).StructTag(`json:"is_used"`),
		field.Int64("expired_at").
			Positive().
			Comment("The activation link is only valid for 1 hour.").
			StructTag(`json:"expired_at"`),
	}
}

// Edges of the Activation.
func (Activation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("activation").
			Required().
			Unique(),
	}
}

// Annotations of the Activation.
func (Activation) Annotations() []schema.Annotation {
	withComment := true

	return []schema.Annotation{
		&entsql.Annotation{
			WithComments: &withComment,
			Checks: map[string]string{
				"token": "length(token) = 44",
			},
		},
	}
}
