package query

// BaseQuery is a base query, consisting of a key-value pair.
type BaseQuery struct {
	Field string
	Value string
}

// A Disjunction is an OR query.
type Disjunction struct {
	SubQueries []Query
}

// A Conjunction is an AND query.
type Conjunction struct {
	SubQueries []Query
}

// A LinkedConjunction is a linked conjunction query.j
type LinkedConjunction struct {
	SubQueries []Query
}

// A Negation is a negation of the sub-query.
type Negation struct {
	SubQuery Query
}

// A Query represents a query.
// It exists purely to facilitate polymorphism.
// It is guaranteed to be one of the types defined above in this file.
// Callers have to use GetUnderlying() and do a type-switch.
type Query struct {
	underlying interface{}
}

// GetUnderlying returns the underlying query of this query.
// It returns nil for invalid queries.
// Callers will need to do a type-switch on the returned value.
// It is guaranteed to be one of the types defined above in this file.
// (ie, one of the files that have a corresponding New function)
func (q *Query) GetUnderlying() interface{} {
	return q.underlying
}

// NewDisjunction creates a new disjunction query.
func NewDisjunction(subQueries []Query) Query {
	return Query{underlying: &Disjunction{SubQueries: subQueries}}
}

// NewConjunction creates a new conjunction query.
func NewConjunction(subQueries []Query) Query {
	return Query{underlying: &Conjunction{SubQueries: subQueries}}
}

// NewLinkedConjunction creates a new linked conjunction query.
func NewLinkedConjunction(subQueries []Query) Query {
	return Query{underlying: &LinkedConjunction{SubQueries: subQueries}}
}

// NewNegation creates a new negation query.
func NewNegation(subQuery Query) Query {
	return Query{underlying: &Negation{SubQuery: subQuery}}
}

// NewBase creates a new base query.
func NewBase(field, value string) Query {
	return Query{underlying: &BaseQuery{
		Field: field,
		Value: value,
	}}
}
