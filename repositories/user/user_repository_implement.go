package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	model "neptune/backend/models/user"
)

type userRepositoryImplement struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) UserRepository {
	return &userRepositoryImplement{Db: Db}
}

func (u userRepositoryImplement) Create(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepositoryImplement) FindByUsername(username string) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepositoryImplement) FindById(id uuid.UUID) (model.User, error) {
	//TODO implement me
	panic("implement me")
}
