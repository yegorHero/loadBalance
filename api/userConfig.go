package api

import (
	"net/http"
)

type UserConfig struct {
	repo userConfigRepo
}

type userConfigRepo interface {
	Create()
	Read()
	Update()
	Delete()
}

func NewUserConfig(repo userConfigRepo) *UserConfig {
	return &UserConfig{
		repo: repo,
	}
}

func (u UserConfig) Create(w http.ResponseWriter, r *http.Request) {

}

func (u UserConfig) Read(w http.ResponseWriter, r *http.Request) {

}

func (u UserConfig) Update(w http.ResponseWriter, r *http.Request) {

}

func (u UserConfig) Delete(w http.ResponseWriter, r *http.Request) {

}
