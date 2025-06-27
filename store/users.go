package store

import (
	"context"

	"github.com/JamieLeeNZ/url-shortener/models"
)

type UserStore interface {
	GetOrCreateUser(ctx context.Context, user models.User) (models.User, error)
}
