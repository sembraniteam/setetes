package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PMILocation holds the schema definition for the PMILocation entity.
type PMILocation struct {
	ent.Schema
}

// Mixin for the PMILocation.
func (PMILocation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the PMILocation.
func (PMILocation) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(164).NotEmpty().StructTag(`json:"name"`),
		field.Int("bed_capacities").Default(0).Positive().StructTag(`json:"bed_capacities"`),
		field.Float("lat").SchemaType(map[string]string{dialect.Postgres: "numeric(9,6)"}).StructTag(`json:"lat"`),
		field.Float("lng").SchemaType(map[string]string{dialect.Postgres: "numeric(9,6)"}).StructTag(`json:"lng"`),
		field.Text("street").NotEmpty().StructTag(`json:"street"`),
		field.String("email").MinLen(3).MaxLen(164).Unique().Nillable().StructTag(`json:"email"`),
		field.String("phone_number").MinLen(11).MaxLen(13).Unique().StructTag(`json:"phone_number"`),
		field.Int("opens_at").Positive().StructTag(`json:"opens_at"`),
		field.Int("closes_at").Positive().StructTag(`json:"closes_at"`),
	}
}

// Edges of the PMILocation.
func (PMILocation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subdistrict", Subdistrict.Type).StorageKey(edge.Column("subdistrict_id")),
	}
}
