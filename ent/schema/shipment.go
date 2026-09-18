package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
)

// Shipment holds the schema definition for the Shipment entity.
type Shipment struct {
	ent.Schema
}

// Fields of the Shipment.
func (Shipment) Fields() []ent.Field {
	return []ent.Field{
		field.String("tracking_number").
			Unique(),

		field.String("sender_name"),

		field.String("receiver_name"),

		field.String("origin"),

		field.String("destination"),

		field.String("status").
			Default("pending"),

		field.Float("weight"),
	}
}

func (Shipment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Schema: "logistics",
		},
	}
}

// Edges of the Shipment.
func (Shipment) Edges() []ent.Edge {
	return nil
}
