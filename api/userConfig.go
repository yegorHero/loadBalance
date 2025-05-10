package api

import (
	"net/http"
)

type UserConfig struct {
	storage userConfigStorage
}

type userConfigStorage interface {
	Create()
	Read()
	Update()
	Delete()
}

func NewUserConfig(storage userConfigStorage) *UserConfig {
	return &UserConfig{
		storage: storage,
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
