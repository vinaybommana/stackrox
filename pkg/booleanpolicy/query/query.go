package query

// An Operator denotes how to combine multiple values.
type Operator int

// This block enumerates valid operators.
const (
	Unset Operator = iota
	And
	Or
)

// FieldQuery is a base query, consisting of a query on a specific field.
// This corresponds to a PolicyGroup.
type FieldQuery struct {
	Field    string
	Values   []string
	Operator Operator
	Negate   bool
}

// A Query represents a query.
// This corresponds to a single policy section.
type Query struct {
	FieldQueries []*FieldQuery
}

// SimpleMatchFieldQuery is a convenience function that constructs a simple query
// that matches just the field and value given.
func SimpleMatchFieldQuery(field, value string) *Query {
	return &Query{FieldQueries: []*FieldQuery{
		{Field: field, Values: []string{value}},
	}}
}
