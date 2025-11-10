package postgres

type (
	Message struct {
		Message string `db:"message"`
	}
)
