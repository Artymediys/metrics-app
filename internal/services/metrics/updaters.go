package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Artymediys/metrics-app/internal/storage/psql"
)

func UpdateMetricsPeriodically(ctx context.Context, ticker *time.Ticker, s *psql.Storage) {
	var (
		value float64
		err   error
	)

	if value, err = s.GetPoliciesPurchasedToday(); err != nil {
		log.Println(fmt.Errorf("services.metrics.UpdateMetricsPerdiodically -> s.GetPoliciesPurchasedToday: %w", err))
	} else {
		metrics["policies_purchased_today_gauge"].Set(value)
	}

	if value, err = s.GetAuthenticationsToday(); err != nil {
		log.Println(fmt.Errorf("services.metrics.UpdateMetricsPerdiodically -> s.GetAuthenticationsToday: %w", err))
	} else {
		metrics["authentications_today_gauge"].Set(value)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if value, err = s.GetPoliciesPurchasedToday(); err != nil {
				log.Println(fmt.Errorf("services.metrics.UpdateMetricsPerdiodically -> s.GetPoliciesPurchasedToday: %w", err))
			} else {
				metrics["policies_purchased_today_gauge"].Set(value)
			}

			if value, err = s.GetAuthenticationsToday(); err != nil {
				log.Println(fmt.Errorf("services.metrics.UpdateMetricsPerdiodically -> s.GetAuthenticationsToday: %w", err))
			} else {
				metrics["authentications_today_gauge"].Set(value)
			}
		}
	}
}
