package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const (
	RhesusPositive = "POSITIVE"
	RhesusNegative = "NEGATIVE"
)

const (
	BloodA  = "A"
	BloodB  = "B"
	BloodAB = "AB"
	BloodO  = "O"
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
			Values(BloodA, BloodB, BloodAB, BloodO).
			StructTag(`json:"group"`).
			Comment("comment:The ABO blood group classification (A, B, AB, or O)."),
		field.Enum("rhesus").
			Values(RhesusPositive, RhesusNegative).
			StructTag(`json:"rhesus"`).
			Nillable().
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
