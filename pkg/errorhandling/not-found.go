package errorhandling

import (
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/core/util"
)

type NotFoundError struct {
	Plugin *util.PluginName
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(
		"could not find plugin '%s' of type %s",
		e.Plugin.Name,
		e.Plugin.Type,
	)
}
