package plugin

import (
	"github.com/ankorstore/ankorstore-cli-core/pkg/plugin"
	"github.com/ankorstore/ankorstore-cli-core/pkg/plugin/command"
)

func IncludeMe(m plugin.Manager) {
	m.RegisterPluginTypes(command.Type)
}
