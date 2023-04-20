package service

import "context"

type Logger interface {
	Printf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Service struct {
	Version     string
	Name        string
	Description string
	Logger      Logger
}

func (s Service) Proceed(payload func(context.Context)) {
	s.entryPoint(payload)
}
