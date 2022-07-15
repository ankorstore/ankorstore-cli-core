package plugin

import (
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/hashicorp/go-plugin"
)

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "ANKOR_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"
	FailedPlugin     = "error"
)

type PluginConfig map[string]string
type Type interface {
	String() string
	GRPCClient() (plugin.Plugin, error)
	GRPCServer(Plugin) (plugin.Plugin, error)
}
type Plugin interface {
	PluginInfo() *util.PluginInfo
	Init() (PluginConfig, error)
	Type() Type
}
type NotFoundError struct {
	Plugin string
}
