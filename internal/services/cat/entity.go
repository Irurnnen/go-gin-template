package cat

import "time"

type (
	CreateCat struct {
		Name           string
		Breed          string
		BirthTimestamp time.Time
		Color          string
		Gender         string
		Weight         int
		Description    string
	}

	CatID struct {
		ID string
	}

	CatSearchParams struct {
		Name   *string
		Breed  *string
		Age    *int
		Limit  *int
		Offset *int
	}

	Cat struct {
		ID        string
		Name      string
		Status    string
		Breed     string
		Age       int
		CreatedAt int
		UpdatedAt int
	}

	CatStatus struct {
		Status string
	}

	CreateToy struct {
		Name string
	}

	ToyID struct {
		ID string
	}

	Toy struct {
		ID         string
		Name       string
		Created_at int
	}
)
