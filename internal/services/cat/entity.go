package cat

import "time"

const (
	CatStatusAvailable CatStatus = "available"
	CatStatusAdopted   CatStatus = "adopted"
	CatStatusSleeping  CatStatus = "sleeping"
	CatStatusPlaying   CatStatus = "playing"
	CatStatusSick      CatStatus = "sick"
)

type (
	CatStatus string

	CreateCat struct {
		Name           string
		Breed          string
		BirthTimestamp time.Time
		Color          string
		Gender         string
		Weight         int
		Description    string
	}

	CreateCatRepo struct {
		ID             string
		Name           string
		Breed          string
		Status         CatStatus
		BirthTimestamp time.Time
		Color          string
		Gender         string
		Weight         int
		Description    string
		CreatedAt      time.Time
		UpdatedAt      time.Time
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
		Status    CatStatus
		Breed     string
		Age       int
		CreatedAt int
		UpdatedAt int
	}

	CatStatusStruct struct {
		Status CatStatus
	}

	CreateToy struct {
		Name     string
		Type     string
		Color    string
		Material string
	}

	CreateToyRepo struct {
		ID        string
		Name      string
		Type      string
		Color     string
		Material  string
		CreatesAt time.Time
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
