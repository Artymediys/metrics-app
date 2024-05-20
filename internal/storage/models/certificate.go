package models

import "time"

type CertificateError struct {
	ID              int64
	URL             string
	RequestDateTime time.Time
	Error           string
	Type            string
	App             string
}
