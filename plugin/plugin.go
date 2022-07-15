package plugin

import (
	"github.com/ankorstore/ankor-core/pkg/plugin"
	"github.com/ankorstore/ankor-core/pkg/plugin/command"
)

func IncludeMe(m plugin.Manager) {
	m.RegisterPluginTypes(command.Type)
}
