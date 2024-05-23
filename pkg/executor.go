package api

import (
	"context"
	"fmt"
)

type Executor struct {
	Name string
}

func (e Executor) String() string {
	return fmt.Sprintf("Executor<%s>", e.Name)
}

type Broker struct {
	Name     string
	Location string
}

type key int

var executorKey key

func NewExecutorContext(ctx context.Context, executor *Executor) context.Context {
	return context.WithValue(ctx, executorKey, executor)
}

func FromExecutorContext(ctx context.Context) (*Executor, bool) {
	e, ok := ctx.Value(executorKey).(*Executor)
	return e, ok
}

func (executor *Executor) GetAvailableBrokers() []Broker {

	brokers := make([]Broker, 1)
	brokers[0] = Broker{
		Name:     "Morgan Stanley",
		Location: "US",
	}

	return brokers

}
