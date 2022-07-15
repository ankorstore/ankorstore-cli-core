package errorhandling

import (
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/core/util"
)

type InvalidError struct {
	Plugin  *util.PluginName
	Message string
}

func (e InvalidError) Error() string {
	if e.Plugin == nil {
		return fmt.Sprintf("invalid unknown plugin encountered: %s", e.Message)
	}

	return fmt.Sprintf(
		"invalid plugin '%s' of type %s: %s",
		e.Plugin.Name,
		e.Plugin.Type,
		e.Message,
	)
}
