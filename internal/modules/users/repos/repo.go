package repos

import "github.com/italoservio/serviosoftwareusers/internal/modules/users/models"

type UsersRepo interface {
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	DeleteByID(id string) error
	List(input ListInput) (*ListOutput, error)
}

type ListInput struct {
	Limit    int64
	Offset   int64
	SortBy   string
	Order    string
	Email    *[]string
	FullName *[]string
	Roles    *[]string
}

type ListOutput struct {
	Total int64         `bson:"total" json:"total"`
	Items []models.User `bson:"items" json:"items"`
}
