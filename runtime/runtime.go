// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package runtime

import (
	"context"
	"time"

	"github.com/srl-labs/containerlab/types"
)

const (
	DockerRuntime = "docker"
	IgniteRuntime = "ignite"
)

type ContainerRuntime interface {
	// Initializes the Container runtime struct
	Init(...RuntimeOption) error
	// Adds custom configuration items to the container runtime struct
	WithConfig(*RuntimeConfig)
	// Set the network management details (generated by the config.go)
	WithMgmtNet(*types.MgmtNet)
	// Instructs the runtime not to delete the mgmt network on destroy
	WithKeepMgmtNet()
	// Create container (bridge) network
	CreateNet(context.Context) error
	// Delete container (bridge) network
	DeleteNet(context.Context) error
	// Pull container image if not present
	PullImageIfRequired(context.Context, string) error
	// Create container returns an extra interface that can be used to receive signals
	// about the container life-cycle after it was created, e.g. for post-deploy tasks
	CreateContainer(context.Context, *types.NodeConfig) (interface{}, error)
	// Start pre-created container by its name
	StartContainer(context.Context, string) error
	// Stop running container by its name
	StopContainer(context.Context, string) error
	// List all containers matching labels
	ListContainers(context.Context, []*types.GenericFilter) ([]types.GenericContainer, error)
	// Get a netns path using the pid of a container
	GetNSPath(context.Context, string) (string, error)
	// Executes cmd on container identified with id and returns stdout, stderr bytes and an error
	Exec(context.Context, string, []string) ([]byte, []byte, error)
	// ExecNotWait executes cmd on container identified with id but doesn't wait for output nor attaches stdout/err
	ExecNotWait(context.Context, string, []string) error
	// Delete container by its name
	DeleteContainer(context.Context, string) error
	// Getter for runtime config options
	Config() RuntimeConfig
	GetName() string
}

type Initializer func() ContainerRuntime

type RuntimeOption func(ContainerRuntime)

type RuntimeConfig struct {
	Timeout          time.Duration
	GracefulShutdown bool
	Debug            bool
	KeepMgmtNet      bool
}

var ContainerRuntimes = map[string]Initializer{}

func Register(name string, initFn Initializer) {
	ContainerRuntimes[name] = initFn
}

func WithConfig(cfg *RuntimeConfig) RuntimeOption {
	return func(r ContainerRuntime) {
		r.WithConfig(cfg)
	}
}

func WithMgmtNet(mgmt *types.MgmtNet) RuntimeOption {
	return func(r ContainerRuntime) {
		r.WithMgmtNet(mgmt)
	}
}

func WithKeepMgmtNet() RuntimeOption {
	return func(r ContainerRuntime) {
		r.WithKeepMgmtNet()
	}
}
