package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// City holds the schema definition for the City entity.
type City struct {
	ent.Schema
}

// Fields of the City.
func (City) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			StructTag(`json:"id"`).
			Annotations(entsql.DefaultExpr("uuid_generate_v4()")),
		field.String("bps_code").
			MaxLen(4).
			Unique().
			Annotations(entsql.IndexType("HASH")).
			StructTag(`json:"bps_code"`).
			SchemaType(map[string]string{dialect.Postgres: "char(4)"}),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the City.
func (City) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("province", Province.Type).
			Ref("city").
			Unique().
			Required(),
		edge.To("district", District.Type).
			StorageKey(edge.Column("city_id")).
			Annotations(entsql.OnDelete(entsql.NoAction)),
	}
}
