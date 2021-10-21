/*
Copyright 2020 GramLabs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"io"
	"os"
	"path/filepath"

	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

const (
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

	homeEnv              = "HOME"
	xdgConfigHomeEnv     = "XDG_CONFIG_HOME"
	xdgConfigHomeDefault = ".config"
	xdgConfigDirsEnv     = "XDG_CONFIG_DIRS"
	xdgConfigDirsDefault = "/etc/xdg"
)

// fileLoader loads a configuration from the currently configured filename
func fileLoader(cfg *OptimizeConfig) error {
	f := &file{filename: cfg.Filename}

	// If we are using a configuration file, the filename _must_ be set
	if f.filename == "" {
		f.filename, cfg.Filename = configFilename("stormforge/config")
	}

	if err := f.read(); err != nil {
		return err
	}

	cfg.Merge(&f.data)

	return nil
}

// file represents the data of a configuration file
type file struct {
	filename string
	data     Config
}

// read will decode YAML or JSON data from the specified file into this configuration file
func (l *file) read() error {
	f, err := os.Open(l.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err = yaml2.NewYAMLOrJSONDecoder(f, 4096).Decode(&l.data); err != nil && err != io.EOF {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

// write will encode YAML data from this configuration into the specified file name
func (l *file) write() error {
	output, err := yaml.Marshal(l.data)
	if err != nil {
		return err
	}

	// Create the file (and directories, if necessary). The XDG Base Dir Spec says directories should
	// be created with 0700 and the file may contain sensitive information so we use 0600 for the file.
	if err := os.MkdirAll(filepath.Dir(l.filename), 0700); err != nil {
		return err
	}
	if err := os.WriteFile(l.filename, output, 0600); err != nil {
		return err
	}
	return nil
}

// filename finds the configuration file and returns both the current file and where changes should be written
func configFilename(configFilename string) (string, string) {
	xdgConfigHome := os.Getenv(xdgConfigHomeEnv)
	if xdgConfigHome == "" {
		home := os.Getenv(homeEnv)
		if home == "" {
			home = "~" // TODO Does this work? Or do we need to error out?
		}
		xdgConfigHome = filepath.Join(home, xdgConfigHomeDefault)
	}

	xdgConfigDirs := os.Getenv(xdgConfigDirsEnv)
	if xdgConfigDirs == "" {
		xdgConfigDirs = xdgConfigDirsDefault
	}

	userConfigFilename := filepath.Join(xdgConfigHome, configFilename)
	currentConfigFilename := userConfigFilename
	for _, dir := range append([]string{xdgConfigHome}, filepath.SplitList(xdgConfigDirs)...) {
		filename := filepath.Join(dir, configFilename)
		if _, err := os.Stat(filename); err == nil {
			currentConfigFilename = filename
			break
		}
	}

	return currentConfigFilename, userConfigFilename
}
