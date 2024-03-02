package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/YxTiBlya/ci-core/logger"
	"github.com/YxTiBlya/ci-core/rabbitmq"
	"github.com/YxTiBlya/ci-core/scheduler"

	"github.com/YxTiBlya/ci-monitor/internal/service"
)

type Config struct {
	Service  service.Config  `yaml:"montitor"`
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`
	QSName   string          `yaml:"qs_name"`
}

var cfgPath string

func init() {
	logger.Init(logger.DevelopmentConfig)
	flag.StringVar(&cfgPath, "cfg", "config.yaml", "app cfg path")
	flag.Parse()
}
func main() {
	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal("failed to open config file ", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatal("failed to unmarshal config file ", err)
	}
	cfg.Service.QSName = cfg.QSName

	rmq, err := rabbitmq.NewRabbitMQ(rabbitmq.WithConfig(cfg.RabbitMQ))
	if err != nil {
		log.Fatal("failed to create rabbitmq ", err)
	}

	svc := service.New(cfg.Service, service.Relations{
		QS: rmq,
	})

	sch := scheduler.NewScheduler(
		scheduler.NewComponent("rabbitmq", rmq),
		scheduler.NewComponent("service", svc),
	)
	sch.Start()
}
