package repos

import "github.com/italoservio/serviosoftwareusers/internal/modules/users/models"

type UsersRepo interface {
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	UpdateByID(id string, user *models.User) (*models.User, error)
	DeleteByID(id string) error
	List(input *ListInput) (*ListOutput, error)
}

type ListInput struct {
	Limit    int64     `validate:"min=1,max=100"`
	Page     int64     `validate:"min=1"`
	SortBy   *string   `validate:"omitempty,omitnil,oneof=fullName email createdAt"`
	Order    *string   `validate:"omitempty,omitnil,oneof=asc desc"`
	Email    *[]string `validate:"omitempty,omitnil,dive,email"`
	FullName *[]string `validate:"omitempty,omitnil,dive,min=2,max=200"`
	Roles    *[]string `validate:"omitempty,omitnil,dive,oneofrole"`
}

type ListOutput struct {
	Total int64         `bson:"total" json:"total"`
	Items []models.User `bson:"items" json:"items"`
}
