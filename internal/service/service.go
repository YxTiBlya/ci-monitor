package service

import (
	"context"

	"github.com/pingcap/errors"
	"go.uber.org/zap"

	"github.com/YxTiBlya/ci-core/rabbitmq"
)

type Relations struct {
	QS QueryService
}

func New(cfg Config, log *zap.SugaredLogger, rel Relations) *Service {
	return &Service{
		cfg:       cfg,
		log:       log,
		Relations: rel,
	}
}

type Service struct {
	cfg Config
	log *zap.SugaredLogger // TODO: thats bad but maybe i do the interface later??
	Relations
}

func (svc *Service) Start(ctx context.Context) error {
	err := svc.Relations.QS.AddMigrates(
		rabbitmq.WithQueue(&rabbitmq.QueueConfig{Name: svc.cfg.QSName}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to migrate rabbitmq")
	}

	updatesCh := make(chan string, 10) // TODO: mov 10 to config?
	go svc.Monitor(updatesCh)
	go svc.Forwarder(updatesCh)

	return nil
}

func (svc *Service) Stop(ctx context.Context) error {

	// TODO: i have someone that need stop?

	return nil
}
