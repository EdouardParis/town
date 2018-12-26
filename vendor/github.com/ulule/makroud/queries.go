package makroud

import (
	"bytes"
	"strings"

	lkb "github.com/ulule/loukoum/builder"
)

// Queries is a list of Query instances.
type Queries []Query

func (q Queries) String() string {
	buffer := &bytes.Buffer{}
	for i := range q {
		query := q[i].String()
		buffer.WriteString(query)
		if !strings.HasSuffix(query, ";") {
			buffer.WriteString(";")
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// Query is a query generated by loukoum and makroud.
type Query struct {
	Raw   string
	Query string
	Args  map[string]interface{}
}

// String returns query statement.
func (q Query) String() string {
	return q.Raw
}

// NewQuery creates a new Query instance from given loukoum builder.
func NewQuery(builder lkb.Builder) Query {
	raw := builder.String()
	query, args := builder.NamedQuery()
	return Query{
		Raw:   raw,
		Query: query,
		Args:  args,
	}
}

// NewRawQuery creates a new Query instance from given query.
func NewRawQuery(raw string) Query {
	return Query{
		Raw: raw,
	}
}