package service

import (
	"github.com/BinaryHexer/go-deptrac/examples/simple-mvc/repository"
)

type Service interface{}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return &service{
		repository: repository,
	}
}
