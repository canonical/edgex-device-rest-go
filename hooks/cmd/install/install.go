// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"

	hooks "github.com/canonical/edgex-snap-hooks"
)

var cli *hooks.CtlCli = hooks.NewSnapCtl()

// installProfiles copies the profile configuration.toml files from $SNAP to $SNAP_DATA.
func installConfig() error {
	var err error

	path := "/config/device-rest-go/res/configuration.toml"
	destFile := hooks.SnapData + path
	srcFile := hooks.Snap + path

	if err = os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
		return err
	}

	if err = hooks.CopyFile(srcFile, destFile); err != nil {
		return err
	}

	return nil
}

// TODO: merge into the above function...
func installDevProfiles() error {
	var err error

	profs := [...]string{"image", "numeric", "json"}

	for _, v := range profs {
		path := fmt.Sprintf("/config/device-rest-go/res/sample-%s-device.yaml", v)
		destFile := hooks.SnapData + path
		srcFile := hooks.Snap + path

		if err := os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return err
		}

		if err = hooks.CopyFile(srcFile, destFile); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var err error

	if err = hooks.Init(false, "edgex-device-rest"); err != nil {
		fmt.Println(fmt.Sprintf("edgex-device-rest::install: initialization failure: %v", err))
		os.Exit(1)

	}

	err = installConfig()
	if err != nil {
		hooks.Error(fmt.Sprintf("edgex-device-rest:install: %v", err))
		os.Exit(1)
	}

	err = installDevProfiles()
	if err != nil {
		hooks.Error(fmt.Sprintf("edgex-device-rest:install: %v", err))
		os.Exit(1)
	}

	// disable the service and handle the autostart logic in the configure hook
	// as default snap configuration is not available when the install hook runs
	// TODO: update the service name to drop the "-go"
	err = cli.Stop("device-rest-go", true)
	if err != nil {
		hooks.Error(fmt.Sprintf("Can't stop service - %v", err))
		os.Exit(1)
	}
}
