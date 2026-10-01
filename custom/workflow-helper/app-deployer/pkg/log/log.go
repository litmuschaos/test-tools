// Package log is vendored from github.com/litmuschaos/test-tools/pkg/log at
// v1.8.0. That release is the only published version of the module, is built
// against go 1.13, and therefore pulled its whole 2020-era requirement graph
// (k8s.io/apiserver, go-restful v2, x/crypto, ...) into this module's go.sum.
package log

import (
	logrus "github.com/sirupsen/logrus"
)

// Fatalf Logs first and then calls `logger.Exit(1)`
// logging level is set to Panic.
func Fatalf(msg string, err error) {
	logrus.WithFields(logrus.Fields{}).Fatalf(msg, err)
}

// Infof log the General operational entries about what's going on inside the application
func Infof(msg string, val string) {
	logrus.WithFields(logrus.Fields{}).Infof(msg, val)
}

// Info log the General operational entries about what's going on inside the application
func Info(msg string) {
	logrus.WithFields(logrus.Fields{}).Infof(msg)
}

// InfoWithValues log the General operational entries about what's going on inside the application
// It also print the extra key values pairs
func InfoWithValues(msg string, val map[string]interface{}) {
	logrus.WithFields(val).Info(msg)
}

// Warn log the Non-critical entries that deserve eyes.
func Warn(msg string) {
	logrus.WithFields(logrus.Fields{}).Warn(msg)
}

// Errorf used for errors that should definitely be noted.
// Commonly used for hooks to send errors to an error tracking service.
func Errorf(msg string, err error) {
	logrus.WithFields(logrus.Fields{}).Errorf(msg, err)
}

// Error used for errors that should definitely be noted.
// Commonly used for hooks to send errors to an error tracking service
func Error(msg string) {
	logrus.WithFields(logrus.Fields{}).Error(msg)
}
