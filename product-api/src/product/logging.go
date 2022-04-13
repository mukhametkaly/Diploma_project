package product

import (
	"github.com/go-kit/kit/log"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	logger log.Logger
	Service
}

//var Loger *log.Logger
var Loger *logrus.Logger

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}
