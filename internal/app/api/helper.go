package api

import (
	"register/storage"

	"github.com/sirupsen/logrus"
)

var (
	prefix string = "/api/v1"
)

func (a *API) configreLoggerField() error {
	log_level, err := logrus.ParseLevel(a.config.LoggerLevel)
	if err != nil {
		return err
	}
	a.logger.SetLevel(log_level)
	return nil
}
func (a *API) configreRouterField() {
	a.router.HandleFunc(prefix+"/users", a.GetAllUsers).Methods("GET")
	a.router.HandleFunc(prefix+"/users/register", a.RegisterUser).Methods("POST")
}

func (a *API) configureStorageField() error {
	storage := storage.New(a.config.Storage)

	if err := storage.Open(); err != nil {
		return err
	}
	a.storage = storage
	return nil
}
