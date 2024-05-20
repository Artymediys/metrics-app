package email

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Artymediys/metrics-app/internal/storage/models"
	"github.com/Artymediys/metrics-app/internal/storage/psql"
	"github.com/Artymediys/metrics-app/internal/utils/html_table_generator"
	"github.com/Artymediys/metrics-app/internal/utils/time_scheduler"
)

func NotifyOfCertificateErrors(
	ctx context.Context,
	ticker *time.Ticker,
	s *psql.Storage,
	eServer *Server,
	contour string,
	addresses *models.EmailAddresses,
) {
	var (
		certErrs    []models.CertificateError
		err         error
		mailBody    string
		mailSubject string
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if certErrs, err = s.FetchDataForCertificateErrors(); err != nil {
				log.Println(fmt.Errorf("services.email.NotifyOfCertificateErrors -> s.FetchDataForCertificateErrors: %w", err))
				continue
			}

			if len(certErrs) == 0 {
				continue
			}

			if mailBody, err = html_table_generator.GenerateHTML(certErrs); err != nil {
				log.Println(fmt.Errorf("services.email.NotifyOfCertificateErrors -> html_table_generator.GenerateHTML: %w", err))
				continue
			}

			mailSubject = fmt.Sprintf(
				"[%s] %d requests with a certificate error in the last 24 hours",
				contour,
				len(certErrs),
			)

			if err = eServer.Send(addresses.From, addresses.To, mailSubject, mailBody); err != nil {
				log.Println(fmt.Errorf("services.email.NotifyOfCertificateErrors -> eServer.Send: %w", err))
			}
		}
	}
}

func NotifyOfUnissuedPolicies(
	ctx context.Context,
	tc *time_scheduler.TimerController,
	s *psql.Storage,
	eServer *Server,
	contour string,
	scheduledTimes []string,
	addresses *models.EmailAddresses,
) {
	var (
		policies    []models.UnissuedPolicy
		err         error
		mailBody    string
		mailSubject string

		funcName = "NotifyOfUnissuedPolicies"
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tc.Timer.C:
			if policies, err = s.FetchDataForUnissuedPolicies(); err != nil {
				log.Println(fmt.Errorf("services.email.NotifyOfUnissuedPolicies -> s.FetchDataForUnissuedPolicies: %w", err))
				executeScheduler(tc, scheduledTimes, funcName)
				continue
			}

			if len(policies) == 0 {
				executeScheduler(tc, scheduledTimes, funcName)
				continue
			}

			if mailBody, err = html_table_generator.GenerateHTML(policies); err != nil {
				log.Println(fmt.Errorf("services.email.NotifyOfUnissuedPolicies -> html_table_generator.GenerateHTML: %w", err))
				executeScheduler(tc, scheduledTimes, funcName)
				continue
			}

			mailSubject = fmt.Sprintf("[%s] Paid but unissued policies", contour)

			if err = eServer.Send(addresses.From, addresses.To, mailSubject, mailBody); err != nil {
				log.Println(fmt.Errorf("services.email.NotifyOfCertificateErrors -> eServer.Send: %w", err))
				executeScheduler(tc, scheduledTimes, funcName)
				continue
			}

			executeScheduler(tc, scheduledTimes, funcName)
		}
	}
}

func executeScheduler(tc *time_scheduler.TimerController, scheduledTimes []string, funcName string) {
	tc.Timer.Stop()

	duration, err := time_scheduler.DurationUntilNextScheduledTime(scheduledTimes)
	if err != nil {
		log.Println(fmt.Errorf("services.email.%s -> time_scheduler.DurationUntilNextScheduledTime: %w", funcName, err))
		tc.Timer = time.NewTimer(3 * time.Hour)
	} else {
		tc.Timer = time.NewTimer(duration)
	}
}
