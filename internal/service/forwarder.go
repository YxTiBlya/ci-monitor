package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/yaml.v3"

	"github.com/YxTiBlya/ci-monitor/pkg/models"
)

func (svc *Service) Forwarder(ch chan string) {
	for data := range ch {
		b, err := os.ReadFile(fmt.Sprintf("%s/pipeline.yaml", data))
		if err != nil {
			svc.log.Error().Err(err).Msg("failed to read pipeline.yaml")
			continue
		}

		var pipeline []models.Pipeline
		if err := yaml.Unmarshal(b, &pipeline); err != nil {
			svc.log.Error().Err(err).Msg("failed to unmarshal pipeline file")
			continue
		}

		b, err = json.Marshal(models.QSPipelineMsg{
			Repo:     data,
			Pipeline: pipeline,
		})
		if err != nil {
			svc.log.Error().Err(err).Msg("failed to marshal msg")
			continue
		}

		msg := &amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         b,
		}

		if err := svc.Relations.QS.Publish(context.Background(), "", svc.cfg.QSName, false, false, msg); err != nil {
			svc.log.Error().Err(err).Str("data", string(b)).Msg("failed to publish data")
			continue
		}

		svc.log.Info().Str("data", string(b)).Msg("succesfully published data")
	}
}
