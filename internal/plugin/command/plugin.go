package command

import (
	error_handler "github.com/ankorstore/ankorstore-cli-core/internal/errorhandling"
	ankor "github.com/ankorstore/ankorstore-cli-core/internal/plugin"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

// TODO: Implement cobra commands via GRPC

const Type pluginType = "command"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() (plugin.Plugin, error) {
	return nil, error_handler.InvalidError{
		Message: "command plugin is not implemented via grpc",
	}
}

func (t pluginType) GRPCServer(p ankor.Plugin) (plugin.Plugin, error) {
	return nil, error_handler.InvalidError{
		Message: "command plugin is not implemented via grpc",
	}
}

type Command interface {
	ankor.Plugin
	GetCobraCommand() *cobra.Command
}
