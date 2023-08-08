package repository

import (
	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserRepository struct {
}

func (r *UserRepository) Create(ctx *gin.Context, u domain.User) {

}
