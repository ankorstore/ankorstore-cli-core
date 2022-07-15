package plugin

import (
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/core/util"
	error_handler "github.com/ankorstore/ankorstore-cli-core/pkg/errorhandling"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
	"path/filepath"
)

type Manager struct {
	pluginTypes map[string]Type
	plugins     map[string]map[string]Plugin
}

func Init() Manager {
	//config.Configure()
	//
	//if os.Getenv(MagicCookieKey) == MagicCookieValue {
	//	log.SetDefault(log.New(config.GetPluginLoggerOptions()))
	//} else {
	//	log.SetDefault(log.New(config.GetLoggerOptions()))
	//}

	return Manager{
		pluginTypes: map[string]Type{},
		plugins:     map[string]map[string]Plugin{},
	}
}

func (m Manager) RegisterPluginTypes(ts ...Type) {
	for _, t := range ts {
		m.pluginTypes[t.String()] = t
	}
}

func (m Manager) IncludePlugins(ps ...Plugin) {
	for _, p := range ps {
		name := p.PluginInfo().Name

		if m.plugins[name.Type] == nil {
			m.plugins[name.Type] = map[string]Plugin{}
		}

		m.plugins[name.Type][name.Name] = p
	}
}

func (m Manager) ServePlugins(plugins ...Plugin) error {
	pluginMap := map[string]plugin.Plugin{}

	for _, p := range plugins {
		s, err := p.Type().GRPCServer(p)
		if err != nil {
			return fmt.Errorf("could not instantiate GRPC Server: %w", err)
		}

		pluginMap[p.PluginInfo().Name.Type] = s
	}

	plugin.Serve(&plugin.ServeConfig{
		GRPCServer: plugin.DefaultGRPCServer,
		Plugins:    pluginMap,
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  ProtocolVersion,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
	})

	return nil
}

func (m Manager) GetPlugins(pluginType string) map[string]Plugin {
	return m.plugins[pluginType]
}

func (m Manager) LoadPlugins() {
	m.findPlugins()
}

func (Manager) UnloadPlugins() {
	plugin.CleanupClients()
}

func (m Manager) findPlugins() {
	matches, err := filepath.Glob(util.PathByName("*"))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.Default().With("path", path)

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0o111 == 0 {
			continue
		}

		for n, t := range m.pluginTypes {
			logger.Debug("attempt loading plugin", "type", n)

			grpcClient, err := t.GRPCClient()
			if err != nil {
				logger.Debug("error loading plugin", "error", err)

				continue
			}

			p, err := loadGRPCPlugin(path, n, grpcClient)
			if err != nil {
				logger.Debug("could not load plugin over grpc", "error", err)

				continue
			}

			name := p.PluginInfo().Name
			if m.plugins[name.Type] == nil {
				m.plugins[name.Type] = map[string]Plugin{}
			}

			m.plugins[name.Type][name.Name] = p
		}
	}
}

func (m Manager) initPlugin(p Plugin) {
	info := p.PluginInfo()
	logger := log.L().With("name", info.Name.Name, "type", info.Name.Type)
	logger = augmentLogger(logger, info.Fields)
	if config, err := p.Init(); err != nil {
		logger.Warn("could not load plugin", "error", err)
		delete(m.plugins[info.Name.Type], info.Name.Name)
	} else {
		augmentLogger(logger, config).Info("loaded plugin")
	}
}

func loadGRPCPlugin(path, pluginType string, grpcPlugin plugin.Plugin) (Plugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		Managed:          true,
		Plugins:          map[string]plugin.Plugin{pluginType: grpcPlugin},
		Cmd:              exec.Command(path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		//Logger:           log.New(config.GetLoggerOptions()),
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  ProtocolVersion,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
	})

	conn, err := client.Client()
	if err != nil {
		client.Kill()

		return nil, fmt.Errorf("error getting plugin client: %w", err)
	}

	raw, err := conn.Dispense(pluginType)
	if err != nil {
		client.Kill()

		return nil, fmt.Errorf("error dispensing plugin: %w", err)
	}

	if p, ok := raw.(Plugin); ok {
		return p, nil
	}

	client.Kill()

	return nil, error_handler.InvalidError{
		Plugin:  &util.PluginName{Type: pluginType, Name: path},
		Message: "does not implement Plugin interface",
	}
}

func augmentLogger(logger log.Logger, fields map[string]string) log.Logger {
	fs := []interface{}{}

	for k, v := range fields {
		fs = append(fs, k, v)
	}

	return logger.With(fs...)
}
