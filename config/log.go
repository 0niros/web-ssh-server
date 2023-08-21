package config

import "github.com/sirupsen/logrus"

func InitLog() {
	logrus.SetLevel(logrus.TraceLevel)
}
