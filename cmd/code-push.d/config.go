package main

import (
	"errors"
	"fmt"
	"github.com/funnyecho/code-push/pkg/fs"
	"strings"
)

type serveConfig struct {
	ConfigFilePath string
	Debug          bool
	PortGrpc       int
	PortHttp       int
	BoltPath       string
}

func (c *serveConfig) Validate() error {
	var errs []string

	if c.PortGrpc == 0 {
		errs = append(errs, "Invalid Grpc Port")
	}

	if c.PortHttp == 0 {
		errs = append(errs, "Invalid Http Port")
	}

	if c.BoltPath == "" {
		errs = append(errs, "BoltPath required")
	} else {
		boltFile, boltFileErr := fs.File(fs.FileConfig{
			FilePath: c.BoltPath,
		})
		if boltFileErr != nil {
			errs = append(errs, boltFileErr.Error())
		} else {
			if dirErr := boltFile.EnsurePath(); dirErr != nil {
				errs = append(errs, dirErr.Error())
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(fmt.Sprintf("FA_CONFIG_SERVE:\n\t%s", strings.Join(errs[:], "\n\t")))
}
