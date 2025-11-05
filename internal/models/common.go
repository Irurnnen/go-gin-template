package models

type (
	Message struct {
		Message string `json:"message"`
	}

	HTTPError struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
)
