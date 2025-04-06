package authModule

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Credential struct {
		Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		PlayerId     string             `json:"player_id" bson:"player_id"`
		RoleCode     int32              `json:"role_code" bson:"role_code"`
		RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
		AccessToken  string             `json:"access_token" bson:"access_token"`
		CreatedAt    *time.Time         `json:"created_at"`
		UpdatedAt    *time.Time         `json:"updated_at"`
	}

	Role struct {
		Id    primitive.ObjectID `json:"id" bson:"id"`
		Title string             `json:"title" bson:"title"`
		Code  int                `json:"code" bson:"code"`
	}
)
