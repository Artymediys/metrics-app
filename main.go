package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/Artymediys/metrics-app/internal/services/email"
	"github.com/Artymediys/metrics-app/internal/services/metrics"
	"github.com/Artymediys/metrics-app/internal/storage/models"
	"github.com/Artymediys/metrics-app/internal/storage/psql"
	"github.com/Artymediys/metrics-app/internal/utils/time_scheduler"
)

func main() {
	// ============================================
	// ========= ENVIRONMENT SETTINGS =============
	// ============================================

	var (
		wg      sync.WaitGroup
		logFile *os.File
		contour string
		err     error
	)

	// LOGGER
	logFile, err = os.OpenFile("./log/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(fmt.Errorf("main -> os.OpenFile: %w", err))
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// CONFIGURATION
	flag.StringVar(&contour, "contour", "demo", "Choose the app contour (local/demo/preprod/prod)")
	flag.Parse()

	viper.SetConfigName("config." + contour)
	viper.AddConfigPath("./config/")
	if err = viper.ReadInConfig(); err != nil {
		log.Panicln(fmt.Errorf("main -> viper.ReadInConfig: %w", err))
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	// ============================================
	// ========= APPLICATION PREPARATIONS =========
	// ============================================

	db, err := psql.New(viper.GetString("postgres_conn"))
	if err != nil {
		log.Panicln(fmt.Errorf("main -> psql.New: %w", err))
	}
	defer db.Stop()

	eServer := email.NewServer(viper.GetString("smtp_host"), viper.GetString("smtp_port"))

	metrics.Register()

	e := echo.New()
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	ticker1H, ticker3H := time.NewTicker(1*time.Hour), time.NewTicker(3*time.Hour)
	defer func() {
		ticker1H.Stop()
		ticker3H.Stop()
	}()

	scheduledTimes := []string{"08:50", "11:50", "14:50", "17:50"}
	initialDuration, err := time_scheduler.DurationUntilNextScheduledTime(scheduledTimes)
	if err != nil {
		log.Panicln(fmt.Errorf("main -> time_scheduler.DurationUntilNextScheduledTime: %w", err))
	}

	tc := &time_scheduler.TimerController{Timer: time.NewTimer(initialDuration)}
	defer func() { tc.Timer.Stop() }()

	// ======================================
	// ========= APPLICATION LAUNCH =========
	// ======================================
	wg.Add(4)

	go func() {
		defer wg.Done()
		metrics.UpdateMetricsPeriodically(ctx, ticker1H, db)
	}()

	go func() {
		defer wg.Done()

		email.NotifyOfCertificateErrors(
			ctx, ticker3H, db, eServer, contour,
			&models.EmailAddresses{
				From: viper.GetString("email_sender"),
				To:   viper.GetStringSlice("email_certs_receivers"),
			},
		)
	}()

	go func() {
		defer wg.Done()

		email.NotifyOfUnissuedPolicies(
			ctx, tc, db, eServer, contour, scheduledTimes,
			&models.EmailAddresses{
				From: viper.GetString("email_sender"),
				To:   viper.GetStringSlice("email_policy_receivers"),
			},
		)
	}()

	go func() {
		defer wg.Done()

		err = e.Start(":8084")
		if err != nil && err != http.ErrServerClosed {
			log.Panicln(fmt.Errorf("main -> e.Start: %w", err))
		}
	}()

	<-exitSignal
	cancel()

	ticker1H.Stop()
	ticker3H.Stop()
	tc.Timer.Stop()

	if err = e.Shutdown(ctx); err != nil {
		log.Panicln(fmt.Errorf("main -> e.Shutdown: %w", err))
	}

	wg.Wait()
}
