package dto

type NewProjectRequest struct {
	Name string `json:"name" binding:"required"`
}

type FindProjectOptions struct {
	Name string
}
