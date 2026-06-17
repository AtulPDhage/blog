package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        string             `bson:"-" json:"_id"`
	BsonID    primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Image     string             `bson:"image" json:"image"`
	Instagram string             `bson:"instagram" json:"instagram"`
	Linkedin  string             `bson:"linkedin" json:"linkedin"`
	Facebook  string             `bson:"facebook" json:"facebook"`
	Bio       string             `bson:"bio" json:"bio"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// SyncIDs ensures that the string ID matches the MongoDB ObjectID
func (u *User) SyncIDs() {
	if u.ID != "" {
		if oid, err := primitive.ObjectIDFromHex(u.ID); err == nil {
			u.BsonID = oid
		}
	} else if !u.BsonID.IsZero() {
		u.ID = u.BsonID.Hex()
	}
}

type Claims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}
