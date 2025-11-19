package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subdistrict holds the schema definition for the Subdistrict entity.
type Subdistrict struct {
	ent.Schema
}

// Mixin for the Subdistrict.
func (Subdistrict) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Subdistrict.
func (Subdistrict) Fields() []ent.Field {
	return []ent.Field{
		field.String("bps_code").MaxLen(10).Unique().Annotations(entsql.IndexType("HASH")).StructTag(`json:"bps_code"`),
		field.String("postal_code").MaxLen(5).Unique().StructTag(`json:"postal_code"`),
		field.String("name").StructTag(`json:"name"`),
	}
}

// Edges of the Subdistrict.
func (Subdistrict) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("district", District.Type).StorageKey(edge.Column("district_id")).Annotations(entsql.OnDelete(entsql.NoAction)),
		edge.From("pmi_location", PMILocation.Type).Ref("subdistrict").Unique().Required(),
	}
}
