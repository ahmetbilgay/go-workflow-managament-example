package models

import (
	"github.com/qiniu/qmgo/field"
)

type User struct {
	field.DefaultField `bson:",inline"`
	Name               string `bson:"name" json:"name"`
	Email              string `bson:"email" json:"email"`
	Age                int    `bson:"age" json:"age"`
}
