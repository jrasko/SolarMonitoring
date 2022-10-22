package graph

//go:generate go run github.com/99designs/gqlgen generate

import "pv-service/processing"

type Resolver struct {
	Processor *processing.Processor
}
