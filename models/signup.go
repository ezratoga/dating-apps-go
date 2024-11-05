package models

type SignUpUser struct {
	Username	string		`json:"username"  validate:"required"`
	Password	string		`json:"password"  validate:"required"`
	Email		string		`json:"email"	  validate:"required,email"`
	Name		string 		`json:"name"	  validate:"required,min=1"`
	Birthday    string	 	`json:"birthday"  validate:"required,min=10"`
	Gender		string		`json:"gender"	  validate:"oneof=male female"`
}