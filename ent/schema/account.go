package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const (
	Female = "FEMALE"
	Male   = "MALE"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Mixin for the Account.
func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("national_id_hash").MaxLen(64).NotEmpty().Unique().Sensitive().Comment("SHA256 hash of a user's national identity number (e.g., KTP). Stored securely to avoid saving raw identity numbers."),
		field.String("full_name").MinLen(3).MaxLen(164).StructTag(`json:"full_name"`),
		field.Enum("gender").Values(Female, Male).StructTag(`json:"gender"`),
		field.String("email").MinLen(3).MaxLen(164).Unique().StructTag(`json:"email"`),
		field.String("country_iso_code").MaxLen(2).NotEmpty().StructTag(`json:"country_iso_code"`).Comment("ISO 3166-1 alpha-2 country code representing the user's country (e.g., ID for Indonesia, US for United States)."),
		field.String("dial_code").StructTag(`json:"dial_code"`).Comment("International dialing code of the user's country (e.g., +62 for Indonesia, +1 for United States). Used for constructing complete phone numbers."),
		field.String("phone_number").MinLen(11).MaxLen(13).Unique().StructTag(`json:"phone_number"`),
		field.Bool("activated").Default(false).StructTag(`json:"activated"`),
		field.Bool("locked").Default(false).StructTag(`json:"locked"`).Comment("Permanently locked by this account."),
		field.Int64("temp_locked_at").Positive().Nillable().StructTag(`json:"temp_locked_at"`).Comment("Temporary locked by this account based on time milliseconds."),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("blood", Blood.Type).StorageKey(edge.Column("blood_id")).Unique().Annotations(entsql.OnDelete(entsql.NoAction)),
		edge.To("password", Password.Type).StorageKey(edge.Column("password_id")).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
