package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct{ ent.Schema }

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("wallet_id", uuid.UUID{}),
		field.Enum("type").Values("DEPOSIT", "WITHDRAW"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(20,2)",
			}),
		field.Time("created_at").Default(time.Now),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wallet", Wallet.Type).Ref("transactions").Field("wallet_id").Unique().Required(),
	}
}
