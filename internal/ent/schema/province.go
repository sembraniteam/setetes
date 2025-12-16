package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Province holds the schema definition for the Province entity.
type Province struct {
	ent.Schema
}

// Fields of the Province.
func (Province) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			StructTag(`json:"id"`).
			Annotations(entsql.DefaultExpr("uuid_generate_v4()")),
		field.String("bps_code").
			MaxLen(2).
			Unique().
			StructTag(`json:"bps_code"`).
			SchemaType(map[string]string{dialect.Postgres: "char(2)"}),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the Province.
func (Province) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("city", City.Type).
			StorageKey(edge.Column("province_id")).
			Annotations(entsql.OnDelete(entsql.NoAction)),
	}
}

// Annotations of the Province.
func (Province) Annotations() []schema.Annotation {
	return []schema.Annotation{
		&entsql.Annotation{
			Checks: map[string]string{
				"bps_code": "length(bps_code) = 2",
			},
		},
	}
}
