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

// District holds the schema definition for the District entity.
type District struct {
	ent.Schema
}

// Fields of the District.
func (District) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			StructTag(`json:"id"`).
			Annotations(entsql.DefaultExpr("uuid_generate_v4()")),
		field.String("bps_code").
			MaxLen(7).
			Unique().
			StructTag(`json:"bps_code"`).
			SchemaType(map[string]string{dialect.Postgres: "char(7)"}),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the District.
func (District) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("city", City.Type).
			Ref("district").
			Unique().
			Required(),
		edge.To("subdistrict", Subdistrict.Type).
			StorageKey(edge.Column("district_id")).
			Annotations(entsql.OnDelete(entsql.NoAction)),
	}
}

// Annotations of the District.
func (District) Annotations() []schema.Annotation {
	return []schema.Annotation{
		&entsql.Annotation{
			Checks: map[string]string{
				"bps_code": "length(bps_code) = 7",
			},
		},
	}
}
