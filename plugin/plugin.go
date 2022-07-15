package plugin

import (
	"github.com/ankorstore/ankorstore-cli-core/internal/plugin"
	"github.com/ankorstore/ankorstore-cli-core/internal/plugin/command"
)

func IncludeMe(m plugin.Manager) {
	m.RegisterPluginTypes(command.Type)
}
