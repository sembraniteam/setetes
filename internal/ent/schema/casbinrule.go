package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CasbinRule holds the schema definition for the CasbinRule entity.
type CasbinRule struct {
	ent.Schema
}

// Fields of the CasbinRule.
func (CasbinRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("Ptype").Optional().Sensitive(),
		field.String("V0").Optional().Sensitive(),
		field.String("V1").Optional().Sensitive(),
		field.String("V2").Optional().Sensitive(),
		field.String("V3").Optional().Sensitive(),
		field.String("V4").Optional().Sensitive(),
		field.String("V5").Optional().Sensitive(),
	}
}

// Edges of the CasbinRule.
func (CasbinRule) Edges() []ent.Edge {
	return nil
}

// Index of the CasbinRule.
func (CasbinRule) Index() []ent.Index {
	return []ent.Index{
		index.Fields("Ptype", "V0", "V1", "V2", "V3", "V4", "V5").Unique(),
	}
}

// Annotations of the CasbinRule.
func (CasbinRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		&entsql.Annotation{
			Table: "casbin_rule",
		},
	}
}
