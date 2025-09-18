package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct{ ent.Schema }

func (Wallet) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user_id", uuid.UUID{}).Unique(),
		field.Other("balance", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(20,2)",
			}),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Wallet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("wallet").Field("user_id").Required().Unique(),
		edge.To("transactions", Transaction.Type),
	}
}

func (Wallet) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id").Unique(),
	}
}
