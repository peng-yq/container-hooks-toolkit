package logger

import "github.com/sirupsen/logrus"

// New returns a new logger
func New() Interface {
	return logrus.StandardLogger()
}

// NullLogger is a logger that does nothing
type NullLogger struct{}

var _ Interface = (*NullLogger)(nil)

// Debugf is a no-op for the null logger
func (l *NullLogger) Debugf(string, ...interface{}) {}

// Errorf is a no-op for the null logger
func (l *NullLogger) Errorf(string, ...interface{}) {}

// Info is a no-op for the null logger
func (l *NullLogger) Info(...interface{}) {}

// Infof is a no-op for the null logger
func (l *NullLogger) Infof(string, ...interface{}) {}

// Warning is a no-op for the null logger
func (l *NullLogger) Warning(...interface{}) {}

// Warningf is a no-op for the null logger
func (l *NullLogger) Warningf(string, ...interface{}) {}
