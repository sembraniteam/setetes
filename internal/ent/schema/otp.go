package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OTP holds the schema definition for the OTP entity.
type OTP struct {
	ent.Schema
}

// Mixin for the OTP.
func (OTP) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the OTP.
func (OTP) Fields() []ent.Field {
	return []ent.Field{
		field.String("code_hash").
			MaxLen(64).
			Sensitive().
			Unique().
			NotEmpty().
			SchemaType(map[string]string{dialect.Postgres: "char(64)"}).
			Comment("Hashed OTP code. Will be deleted after it is used."),
		field.Enum("type").NamedValues(
			"Activation", "ACTIVATION",
			"ResetPassword", "RESET_PASSWORD",
			"ChangePassword", "CHANGE_PASSWORD",
		).StructTag(`json:"type"`),
		field.Int64("expired_at").
			Positive().
			Comment("The OTP code is only valid for 30 minutes.").
			StructTag(`json:"expired_at"`),
	}
}

// Edges of the OTP.
func (OTP) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).
			Unique().
			Required().
			StorageKey(edge.Column("account_id")).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Annotations of the OTP.
func (OTP) Annotations() []schema.Annotation {
	withComment := true

	return []schema.Annotation{
		&entsql.Annotation{
			Table:        "otps",
			WithComments: &withComment,
			Checks: map[string]string{
				"code_hash": "length(code_hash) = 64",
			},
		},
	}
}
