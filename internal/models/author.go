package models

type Author struct {
	ID      int
	NameFa  string `validate:"required,min=1,max=100"`
	NameEn  string `validate:"omitempty,min=1,max=100"`
	Country string `validate:"required,min=1,max=50"`
	Bio     string `validate:"omitempty,min=1,max=500"`
}
