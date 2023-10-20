package main

import (
	"flag"
	"fmt"
	"github.com/anVlad11/testapp-20231020/internal/config"
	"github.com/anVlad11/testapp-20231020/internal/server/http"
	"github.com/anVlad11/testapp-20231020/internal/services/math"
	"github.com/anVlad11/testapp-20231020/internal/services/math_handler"
	"github.com/anVlad11/testapp-20231020/pkg/interrupt_closer"
	"github.com/anVlad11/testapp-20231020/pkg/log"
	"github.com/anVlad11/testapp-20231020/pkg/metrics"
	"github.com/anVlad11/testapp-20231020/pkg/trace"
)

var configPath = flag.String(
	"config-path",
	"./config.yaml",
	"Path to the application config",
)

var ServiceName = "testapp"

func main() {
	flag.Parse()

	closer := interrupt_closer.NewCloser(interrupt_closer.DefaultSignals...)

	logger, err := log.NewZapLogger(ServiceName, log.Info)
	if err != nil {
		panic(fmt.Errorf("error: %v, %s, config", err, ServiceName))
	}

	logger.Info("service is starting", log.String("service", ServiceName))

	logger.Info("setting config")
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		logger.Panic("setting config error", log.Error(err))
	}

	logger.Info("setting tracer provider")
	err = trace.SetTracerProvider(cfg.Tracing.Enabled, cfg.Tracing.ProviderURL, ServiceName)
	if err != nil {
		logger.Panic("tracer provider setting error", log.Error(err))
	}

	mathService := math.NewService(logger)
	mathHandlerService := math_handler.NewService(logger, mathService)

	e := http.NewEchoServer(ServiceName)
	server := http.NewServer(
		cfg.HTTPServer,
		mathHandlerService,
		e,
	)

	if cfg.HTTPServer.Listen {
		go func() {
			closer.Add(server)
			logger.Info("starting http server")
			err = server.Start()
			if err != nil {
				logger.Panic("starting http server error", log.Error(err))
			}
		}()
	}

	metricsServer := metrics.NewServer(cfg.Metrics.Address, cfg.Metrics.Port)

	if cfg.Metrics.Listen {
		go func() {
			closer.Add(metricsServer)
			logger.Info("starting metrics server")
			err = metricsServer.ListenAndServe()
			if err != nil {
				logger.Panic("starting metrics server error", log.Error(err))
			}
		}()
	}

	closer.Wait()
}
