package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BloodType holds the schema definition for the BloodType entity.
type BloodType struct {
	ent.Schema
}

// Mixin for the BloodType.
func (BloodType) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the BloodType.
func (BloodType) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("group").
			NamedValues(
				"BloodA", "A",
				"BloodB", "B",
				"BloodAB", "AB",
				"BloodO", "O",
			).
			StructTag(`json:"group"`).
			Comment("comment:The ABO blood group classification (A, B, AB, or O)."),
		field.Enum("rhesus").
			NamedValues(
				"Positive", "POSITIVE",
				"Negative", "NEGATIVE",
			).
			StructTag(`json:"rhesus"`).
			Optional().
			Comment("The Rhesus (Rh) factor of the blood group, either POSITIVE or NEGATIVE."),
	}
}

// Edges of the BloodType.
func (BloodType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("blood_type").
			Unique().
			Required(),
	}
}

// Annotations ot the BloodType.
func (BloodType) Annotations() []schema.Annotation {
	withComment := true

	return []schema.Annotation{
		&entsql.Annotation{
			WithComments: &withComment,
		},
	}
}
