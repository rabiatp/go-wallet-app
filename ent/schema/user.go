package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type User struct{ ent.Schema }

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").NotEmpty(),
		field.String("email").Unique(),
		field.String("password_hash"),
		field.Time("created_at").Default(time.Now),
	}
}

func (User) Edges() []ent.Edge {
  return []ent.Edge{
    edge.To("wallet", Wallet.Type).Unique(), // Wallet.user Ref("wallet") ile eşleşir
  }
}
