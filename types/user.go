package types

type User struct {
	// `json:"id,omitempty"` will remove the field from the JSON output if the field is empty.
	// `json:"_"` will remove the field completely from the JSON output regardless of its contents.
	FirstName string `bson:"firstName" json:"firstName"`
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	LastName  string `bson:"lastName" json:"lastName"`
}
