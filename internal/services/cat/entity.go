package cat

const (
	CatStatusAvailable CatStatus = "available"
	CatStatusAdopted   CatStatus = "adopted"
	CatStatusSleeping  CatStatus = "sleeping"
	CatStatusPlaying   CatStatus = "playing"
	CatStatusSick      CatStatus = "sick"
)

type (
	CatStatus string

	CatCreate struct {
		Name           string
		Breed          string
		BirthTimestamp int64
		Color          string
		Gender         string
		Weight         int
		Description    *string
	}

	CatCreateRepo struct {
		ID             string
		Name           string
		Breed          string
		Status         CatStatus
		BirthTimestamp int64
		Color          string
		Gender         string
		Weight         int
		Description    *string
		CreatedAt      int64
		UpdatedAt      int64
	}

	CatID struct {
		ID string
	}

	CatSearchParams struct {
		Name   *string
		Breed  *string
		Age    *int
		Limit  int
		Offset *int
	}

	Cat struct {
		ID          string
		Name        string
		Status      CatStatus
		Color       string
		Gender      string
		Weight      int
		Description string
		Breed       string
		Age         int
		CreatedAt   int64
		UpdatedAt   int64
	}

	CatStatusStruct struct {
		Status CatStatus
	}

	ToyCreate struct {
		Name     string
		Type     string
		Color    string
		Material string
	}

	ToyCreateRepo struct {
		ID        string
		Name      string
		Type      string
		Color     string
		Material  string
		CreatedAt int64
		UpdatedAt int64
	}

	ToyID struct {
		ID string
	}

	Toy struct {
		ID        string
		Name      string
		Type      string
		Color     string
		Material  string
		CreatedAt int64
		UpdatedAt int64
	}

	ToySearchParams struct {
		Name     *string
		Type     *string
		Color    *string
		Material *string
		Limit    int
		Offset   *int
	}
)
