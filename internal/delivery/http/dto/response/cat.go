package response

type (
	CatID struct {
		ID string `json:"id" example:"a949dfa4-90ce-4c04-932c-841235d8bd85"`
	}

	Cat struct {
		ID          string `json:"id" example:"18d63d5d-a890-43a1-b109-8cc5dfb5643d"`
		Name        string `json:"name" example:"Murzik"`
		Status      string `json:"status" example:"available"`
		Color       string `json:"color" example:"white"`
		Gender      string `json:"gender" example:"male"`
		Weight      int    `json:"weight" example:"4"`
		Description string `json:"description" example:"lorem ipsum"`
		Breed       string `json:"breed" example:"savannah"`
		Age         int    `json:"age" example:"2"`
		CreatedAt   int64  `json:"created_at" example:"1763230758"`
		UpdatedAt   int64  `json:"updated_at" example:"1763230758"`
	}

	ToyID struct {
		ID string `json:"id" example:"a949dfa4-90ce-4c04-932c-841235d8bd85"`
	}

	Toy struct {
		ID        string `json:"id" example:"18d63d5d-a890-43a1-b109-8cc5dfb5643d"`
		Name      string `json:"name" example:"Teddy"`
		Type      string `json:"type" example:"bear"`
		Color     string `json:"color" example:"brown"`
		Material  string `json:"material" example:"synthetic"`
		CreatedAt int64  `json:"created_at" example:"1763256334"`
		UpdatedAt int64  `json:"updated_at" example:"1763256334"`
	}
)
