package arxivlib

import (
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// A Rating is a score given to a paper by a User
type Rating struct {
	UserID  primitive.ObjectID `json:"user_id" bson:"user_id"`
	Score   int                `json:"score" bson:"score"`
	Comment string             `json:"comment" bson:"comment"`
	Date    time.Time          `json:"date" bson:"date"`
}

// RatingsService interacts with the rating-related endpoints in arxivlib's API
type RatingsService interface {
	// Add a rating to a paper
	AddRating(id primitive.ObjectID, r Rating) (added bool, err error)

	// Update a rating for a paper
	UpdateRating(id primitive.ObjectID, r Rating) (updated bool, err error)
}

// ErrRatingNotFound is a failure to find a rating for a specified paper
var ErrRatingNotFound = errors.New("rating not found")
