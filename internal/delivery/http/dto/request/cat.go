package request

type (
	CatCreate struct {
		Name           string  `json:"name" binding:"required,alphanumunicode,max=256" example:"Murzik"`
		Breed          string  `json:"breed" binding:"required,alphanumunicode,max=256" example:"savannah"`
		BirthTimestamp int64   `json:"birth_timestamp" binding:"required,positive" example:"1763250400"`
		Color          string  `json:"color" binding:"required,alphanumunicode,max=256" example:"white"`
		Gender         string  `json:"male" binding:"required,oneof=male female" example:"male"`
		Weight         int     `json:"weight" binding:"required,min=1,max=100" example:"4"`
		Description    *string `json:"description" binding:"alphanumunicode,max=1024" example:"Murzik"`
	}

	CatID struct {
		ID string `uri:"cid" binding:"required,uuid" example:"18d63d5d-a890-43a1-b109-8cc5dfb5643d"`
	}

	CatSearchParams struct {
		Name   *string `form:"name" binding:"omitempty,alphanumunicode,max=256" example:"Murzik"`
		Breed  *string `form:"breed" binding:"omitempty,alphanumunicode,max=256" example:"savannah"`
		Age    *int    `form:"age" binding:"omitempty,positive" example:"4"`
		Offset *int    `form:"offset" binding:"omitempty,positive" example:"0"`
		Limit  int     `form:"limit" binding:"required,positive,max=256" example:"10"`
	}

	CatStatus struct {
		Status string `json:"status" binding:"required,printascii" example:"available"`
	}

	ToyCreate struct {
		Name     string `json:"name" binding:"required,alphanumunicode,max=256" example:"Teddy"`
		Type     string `json:"type" binding:"required,alphanumunicode,max=256" example:"bear"`
		Color    string `json:"color" binding:"required,alphanumunicode,max=256" example:"brown"`
		Material string `json:"material" binding:"required,alphanumunicode,max=256" example:"cloth"`
	}

	ToyID struct {
		ID string `uri:"tid" binding:"required,uuid" example:"18d63d5d-a890-43a1-b109-8cc5dfb5643d"`
	}

	ToySearch struct {
		Name     *string `form:"name" binding:"omitempty,alphanumunicode,max=256" example:"Murzik"`
		Type     *string `form:"type" binding:"omitempty,alphanumunicode,max=256" example:"bear"`
		Color    *string `form:"color" binding:"omitempty,alphanumunicode,max=256" example:"brown"`
		Material *string `form:"material" binding:"omitempty,alphanumunicode,max=256" example:"synthetic"`
		Limit    int     `form:"offset" binding:"omitempty,positive" example:"0"`
		Offset   *int    `form:"limit" binding:"required,positive,max=256" example:"10"`
	}
)
