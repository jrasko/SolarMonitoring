package graph

//go:generate go run github.com/99designs/gqlgen generate

import "pv-service/service"

type Resolver struct {
	Processor processing.Service
}
