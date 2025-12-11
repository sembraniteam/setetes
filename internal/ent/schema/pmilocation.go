package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
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
		field.String("name").
			MinLen(3).
			MaxLen(164).
			NotEmpty().
			StructTag(`json:"name"`),
		field.Int16("bed_capacities").
			Default(0).
			Positive().
			StructTag(`json:"bed_capacities"`).
			SchemaType(map[string]string{dialect.Postgres: "smallserial"}),
		field.Other("lat_lng", &GeoPoint{}).
			SchemaType(map[string]string{dialect.Postgres: "geography(Point,4326)"}).
			StructTag(`json:"lat_lng"`),
		field.Text("street").NotEmpty().StructTag(`json:"street"`),
		field.String("email").
			MinLen(3).
			MaxLen(164).
			Unique().
			Nillable().
			StructTag(`json:"email"`),
		field.String("dial_code").StructTag(`json:"dial_code"`).
			Comment("International dialing code of the user's country (e.g., +62 for Indonesia, +1 for United States). Used for constructing complete phone numbers."),
		field.String("phone_number").
			MinLen(11).
			MaxLen(13).
			Unique().
			StructTag(`json:"phone_number"`),
		field.Time("opens_at").
			SchemaType(map[string]string{dialect.Postgres: "time"}).
			StructTag(`json:"opens_at"`),
		field.Time("closes_at").
			SchemaType(map[string]string{dialect.Postgres: "time"}).
			StructTag(`json:"closes_at"`),
	}
}

// Edges of the PMILocation.
func (PMILocation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subdistrict", Subdistrict.Type).
			StorageKey(edge.Column("subdistrict_id")),
	}
}

// Annotations of the PMILocation.
func (PMILocation) Annotations() []schema.Annotation {
	withComment := true

	return []schema.Annotation{
		&entsql.Annotation{
			WithComments: &withComment,
			Checks: map[string]string{
				"name":         "length(name) >= 3 and length(name) <= 164",
				"email":        "length(email) >= 3 and length(email) <= 164",
				"phone_number": "length(phone_number) >= 11 and length(phone_number) <= 13",
			},
		},
	}
}
