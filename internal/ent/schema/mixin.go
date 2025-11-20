package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable().Unique().StructTag(`json:"id"`).
			Annotations(entsql.DefaultExpr("uuid_generate_v4()")),
		field.Int64("created_at").Positive().Immutable().StructTag(`json:"created_at"`).
			Annotations(entsql.DefaultExpr("FLOOR(EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)")),
		field.Int64("updated_at").Positive().Nillable().StructTag(`json:"updated_at"`).
			UpdateDefault(time.Now().UnixMilli),
		field.Int64("deleted_at").Positive().Nillable().StructTag(`json:"deleted_at"`).
			Comment("Represents soft delete timestamp in milliseconds.").
			Annotations(entsql.WithComments(true)),
	}
}

func (BaseMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
	}
}
