package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// City holds the schema definition for the City entity.
type City struct {
	ent.Schema
}

// Mixin for the City.
func (City) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the City.
func (City) Fields() []ent.Field {
	return []ent.Field{
		field.String("bps_code").MaxLen(4).Unique().Annotations(entsql.IndexType("HASH")).StructTag(`json:"bps_code"`).
			SchemaType(map[string]string{dialect.Postgres: "char(4)"}),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the City.
func (City) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("province", Province.Type).StorageKey(edge.Column("province_id")).Annotations(entsql.OnDelete(entsql.NoAction)),
		edge.From("district", District.Type).Ref("city").Unique().Required(),
	}
}
