package service

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/YxTiBlya/ci-core/rabbitmq"
)

type QueryService interface {
	Publish(ctx context.Context, exchange, key string, mandatory bool, immediate bool, msg *amqp.Publishing) error
	AddMigrates(migates ...rabbitmq.Migrate) error
}
