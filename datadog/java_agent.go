/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package datadog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type JavaAgent struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaAgent(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, logger bard.Logger) JavaAgent {
	contrib, _ := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return JavaAgent{LayerContributor: contrib, Logger: logger}
}

func (j JavaAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, filepath.Base(j.LayerContributor.Dependency.URI))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy artifact to %s\n%w", file, err)
		}
		layer.LaunchEnvironment.Default("BPI_DATADOG_AGENT_PATH", file)
		layer.LaunchEnvironment.Appendf("JAVA_TOOL_OPTIONS", " ", "-javaagent:%s", file)
		return layer, nil
	})
}

func (j JavaAgent) Name() string {
	return j.LayerContributor.LayerName()
}
