package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// District holds the schema definition for the District entity.
type District struct {
	ent.Schema
}

// Mixin for the District.
func (District) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the District.
func (District) Fields() []ent.Field {
	return []ent.Field{
		field.String("bps_code").MaxLen(7).Unique().Annotations(entsql.IndexType("HASH")).StructTag(`json:"bps_code"`).
			SchemaType(map[string]string{dialect.Postgres: "char(7)"}),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the District.
func (District) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("city", City.Type).StorageKey(edge.Column("city_id")).Annotations(entsql.OnDelete(entsql.NoAction)),
		edge.From("subdistrict", Subdistrict.Type).Ref("district").Unique().Required(),
	}
}
