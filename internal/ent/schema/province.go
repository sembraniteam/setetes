package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Province holds the schema definition for the Province entity.
type Province struct {
	ent.Schema
}

// Mixin for the Province.
func (Province) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Province.
func (Province) Fields() []ent.Field {
	return []ent.Field{
		field.String("bps_code").MaxLen(2).Unique().Annotations(entsql.IndexType("HASH")).StructTag(`json:"bps_code"`).
			SchemaType(map[string]string{dialect.Postgres: "char(2)"}),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the Province.
func (Province) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("city", City.Type).Ref("province").Unique().Required(),
	}
}
