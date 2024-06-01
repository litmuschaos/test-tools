package util

import (
	"github.com/sirupsen/logrus"
	"os"
)

func WriteToFile(engine, engineName string) error {
	if engine == "" {
		logrus.Info("No output file provided skipping write")
		return nil
	}
	f, err := os.Create(engine)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(engineName)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
