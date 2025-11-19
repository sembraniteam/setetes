package schema

import (
	"entgo.io/ent"
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

// Blood holds the schema definition for the Blood entity.
type Blood struct {
	ent.Schema
}

// Mixin for the Blood.
func (Blood) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Blood.
func (Blood) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("group").Values(BloodA, BloodB, BloodAB, BloodO).StructTag(`json:"group"`).Comment("comment:The ABO blood group classification (A, B, AB, or O)."),
		field.Enum("rhesus").Values(RhesusPositive, RhesusNegative).Nillable().Comment("The Rhesus (Rh) factor of the blood group, either POSITIVE or NEGATIVE."),
	}
}

// Edges of the Blood.
func (Blood) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("blood").Unique().Required(),
	}
}
