package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/YxTiBlya/ci-monitor/pkg/models"
)

func (svc *Service) Forwarder(ch chan string) {
	for data := range ch {
		b, err := os.ReadFile(fmt.Sprintf("%s/pipeline.yaml", data))
		if err != nil {
			svc.log.Error("failed to open pipeline.yaml", zap.Error(err))
			continue
		}

		var pipeline []models.Pipeline
		if err := yaml.Unmarshal(b, &pipeline); err != nil {
			svc.log.Fatal("failed to unmarshal pipeline file", zap.Error(err))
			continue
		}

		b, err = json.Marshal(models.QSPipelineMsg{
			Repo:     data,
			Pipeline: pipeline,
		})
		if err != nil {
			svc.log.Fatal("failed to marshal msg", zap.Error(err))
			continue
		}

		msg := &amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         b,
		}

		if err := svc.Relations.QS.Publish(context.Background(), "", svc.cfg.QSName, false, false, msg); err != nil {
			svc.log.Error("failed to publish data", zap.String("data", string(b)), zap.Error(err))
			continue // if only 1 in channel thats raise error?
		}

		svc.log.Info("succesfully published data", zap.String("data", string(b)))
	}
}
