package models

import (
	"time"
)

type User struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty" validate:"omitempty,omitnil,mongodb"`
	FirstName string     `bson:"firstName" json:"firstName" validate:"required,min=2,max=100"`
	LastName  string     `bson:"lastName" json:"lastName" validate:"required,min=2,max=100"`
	FullName  string     `bson:"fullName" json:"fullName" validate:"omitempty,omitnil,required,min=2,max=200"`
	Email     string     `bson:"email" json:"email" validate:"required,email,min=5,max=200"`
	Roles     []string   `bson:"roles" json:"roles" validate:"required,min=1,dive,required,oneofrole"`
	Password  string     `bson:"password,omitempty" json:"password,omitempty" validate:"omitempty,min=8,max=100"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedUt" json:"updatedUt"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}
