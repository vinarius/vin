package models

import "time"

type (
	PrimaryKey struct {
		Pk string `json:"-" dynamo:"pk"`
		Sk string `json:"-" dynamo:"sk"`
	}
	Base struct {
		PrimaryKey
		Type      string `json:"type" dynamo:"type"`
		CreatedAt string `json:"createdAt" dynamo:"createdAt"`
	}
	Input struct {
		Pk        string
		Sk        string
		Type      string
		CreatedAt string
	}
)

// HashKey returns the partition key - added to satisfy interface of dynamo.Keyed
func (primaryKey PrimaryKey) HashKey() any {
	return primaryKey.Pk
}

// RangeKey returns the sort key - added to satisfy interface of dynamo.Keyed
func (primaryKey PrimaryKey) RangeKey() any {
	return primaryKey.Sk
}

func New(input Input) Base {
	createdAt := time.Now().UTC().Format(time.RFC3339)
	if input.CreatedAt != "" {
		createdAt = input.CreatedAt
	}

	return Base{
		PrimaryKey: PrimaryKey{
			input.Pk,
			input.Sk,
		},
		Type:      input.Type,
		CreatedAt: createdAt,
	}
}
