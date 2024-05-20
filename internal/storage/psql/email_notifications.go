package psql

import (
	"fmt"

	"github.com/Artymediys/metrics-app/internal/storage/models"
)

func (s *Storage) FetchDataForCertificateErrors() ([]models.CertificateError, error) {
	query := `
			WITH cert_failures AS (
			SELECT id, url, req_dt, error, type, app 
			FROM test.reqlog 
			WHERE req_dt BETWEEN NOW() - INTERVAL '1 day' AND NOW() 
			  AND error LIKE '%PKIX path building failed%'
			) SELECT * FROM cert_failures ORDER BY req_dt;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage.psql.FetchDataForCertificateErrors -> db.Query: %w", err)
	}
	defer rows.Close()

	var certificateErrors []models.CertificateError
	for rows.Next() {
		var ce models.CertificateError
		if err := rows.Scan(&ce.ID, &ce.URL, &ce.RequestDateTime, &ce.Error, &ce.Type, &ce.App); err != nil {
			return nil, fmt.Errorf("storage.psql.fetchDataForCertificateErrors -> rows.Scan: %w", err)
		}
		certificateErrors = append(certificateErrors, ce)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("storage.psql.fetchDataForCertificateErrors -> rows.Next: %w", err)
	}

	return certificateErrors, nil
}

func (s *Storage) FetchDataForUnissuedPolicies() ([]models.UnissuedPolicy, error) {
	query := `
			SELECT io.io_type, o.id as order_id, io.id, io.name, io.email, ip.policy_number, o.order_sum,
			       ip.policy_sum, ip.external_id, o.payment_id, o.insurance_company, o.dttm
			FROM test.order o 
			    JOIN test.io io ON o.id = io.order_id
			    LEFT JOIN test.io_process ip ON io.id = ip.io_id
			    LEFT JOIN LATERAL (
			        SELECT EXISTS(SELECT * FROM test.ewallet_reqlog 
			                               WHERE o.payment_id = payment_id 
			                                 AND type = 'PAYMENT_REVERSE'
			                                 AND status = 'SUCCESS') as reverted
			        ) as ew ON true 
			WHERE o.dttm BETWEEN (NOW() - '8 day'::interval)::date AND (NOW() - '20 minute'::interval) 
			  AND o.payment_id IS NOT NULL 
			  AND o.status != 'SUCCESS' 
			  AND NOT ew.reverted 
			ORDER BY o.dttm DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage.psql.fetchDataForUnissuedPolicies -> db.Query: %w", err)
	}
	defer rows.Close()

	var policies []models.UnissuedPolicy
	for rows.Next() {
		var up models.UnissuedPolicy
		if err := rows.Scan(&up.IOType, &up.OrderID, &up.ID, &up.Name, &up.Email, &up.PolicyNumber, &up.OrderSum,
			&up.PolicySum, &up.ExternalID, &up.PaymentID, &up.InsuranceCompany, &up.DateTime); err != nil {
			return nil, fmt.Errorf("storage.psql.fetchDataForUnissuedPolicies -> rows.Scan: %w", err)
		}
		policies = append(policies, up)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("storage.psql.fetchDataForUnissuedPolicies -> rows.Next: %w", err)
	}

	return policies, nil
}
