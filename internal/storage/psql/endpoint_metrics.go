package psql

import "fmt"

func (s *Storage) GetPoliciesPurchasedToday() (float64, error) {
	const query = `
			SELECT COUNT(*)
			FROM test.order
			WHERE dttm > NOW() - INTERVAL '1 day'
			  AND status = 'SUCCESS'
	`

	var count float64
	if err := s.db.QueryRow(query).Scan(&count); err != nil {
		return -1, fmt.Errorf("storage.psql.GetPoliciesPurchasedToday: %w", err)
	}

	if count == 0 {
		return 1, nil
	}

	return count, nil
}

func (s *Storage) GetAuthenticationsToday() (float64, error) {
	const query = `
			SELECT COUNT(*)
			FROM test.example_sso2_session
			WHERE access_token IS NOT NULL
			  AND refresh_token IS NOT NULL
			  AND created_dttm > NOW() - INTERVAL '1 day'
	`

	var count float64
	if err := s.db.QueryRow(query).Scan(&count); err != nil {
		return -1, fmt.Errorf("storage.psql.GetAuthenticationsToday: %w", err)
	}

	if count == 0 {
		return 1, nil
	}

	return count, nil
}
