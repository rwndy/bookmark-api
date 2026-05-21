package dto

type CreateBookmarkRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Title string `json:"title" validate:"required,min=1,max=255"`
	Tags  string `json:"tags" validate:"max=500"`
}

type UpdateBookmarkRequest struct {
	URL   string `json:"url" validate:"omitempty,url"`
	Title string `json:"title" validate:"omitempty,min=1,max=255"`
	Tags  string `json:"tags" validate:"omitempty,max=500"`
}