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
		field.String("code").
			MaxLen(6).
			Sensitive().
			Unique().
			NotEmpty().
			SchemaType(map[string]string{dialect.Postgres: "char(6)"}).
			Comment("The OTP code must be 6 characters long. Will be deleted after it is used."),
		field.Enum("type").NamedValues(
			"ResetPassword", "RESET_PASSWORD",
			"Register", "REGISTER",
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
		edge.From("account", Account.Type).
			Ref("otp").
			Required().
			Unique(),
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
				"code": "length(code) = 6",
			},
		},
	}
}
