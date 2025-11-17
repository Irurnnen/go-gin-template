package http

import "fmt"

type (
	BindingZone string

	BindingError struct {
		Err  error
		Zone BindingZone
		Code int
	}
)

const (
	BindingURI    BindingZone = "uri"
	BindingBody   BindingZone = "body"
	BindingJSON   BindingZone = "json"
	BindingHeader BindingZone = "header"
	BindingQuery  BindingZone = "query"
)

func (be *BindingError) Error() string {
	return fmt.Sprintf("failed to bind %s: %e", be.Zone, be.Err)
}

func (be *BindingError) Unwrap() error {
	return be.Err
}

func (be *BindingError) Message() string {
	return fmt.Sprintf("Failed to parse: invalid key/value in %s", be.Zone)
}

func (be *BindingError) StatusCode() int {
	return be.Code
}
