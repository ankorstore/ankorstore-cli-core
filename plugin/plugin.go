package plugin

import (
	"github.com/ankorstore/ankorstore-cli-core/core/plugin"
	"github.com/ankorstore/ankorstore-cli-core/core/plugin/command"
)

func IncludeMe(m plugin.Manager) {
	m.RegisterPluginTypes(command.Type)
}
