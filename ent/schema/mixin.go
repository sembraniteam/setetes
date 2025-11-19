package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	now := time.Now()

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable().Unique().StructTag(`json:"id"`).Default(uuid.New),
		field.Int64("created_at").Positive().Immutable().StructTag(`json:"created_at"`).Default(now.UnixMilli()),
		field.Int64("updated_at").Positive().Nillable().StructTag(`json:"updated_at"`).Default(now.UnixMilli()).UpdateDefault(now.UnixMilli()),
		field.Int64("deleted_at").Positive().Nillable().StructTag(`json:"deleted_at"`).Comment("Represents soft delete timestamp in milliseconds."),
	}
}

func (BaseMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
	}
}
