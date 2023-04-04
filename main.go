package main

import (
	"encoding/gob"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/padok-team/yatas/plugins/commons"
)

type YatasPlugin struct {
	logger hclog.Logger
}

type GPTPlugin struct {
	apiKey      string
	model       string
	prompt      string
	forEachTest bool
	onlyFailed  bool
}

// Don't remove this function
func (g *YatasPlugin) Run(c *commons.Config) []commons.Tests {
	g.logger.Debug("message from Yatas GPT Plugin")
	var err error
	if err != nil {
		panic(err)
	}
	var checksAll []commons.Tests

	checks, err := runPlugin(g, c, "gpt")
	if err != nil {
		g.logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)
	return checksAll
}

func UnmarshalGPT(g *YatasPlugin, c *commons.Config) (GPTPlugin, error) {
	var gptCredentials GPTPlugin

	for _, r := range c.PluginConfig {
		gptFound := false
		for key, value := range r {

			switch key {
			case "pluginName":
				if value == "gpt" {
					gptFound = true

				}
			case "config":
				g.logger.Debug("🔎Found config")
				for _, v := range value.([]interface{}) {
					g.logger.Debug("🔎")
					g.logger.Debug("%v", v)
					for keys, values := range v.(map[string]interface{}) {
						switch keys {
						case "apiKey":
							gptCredentials.apiKey = values.(string)
						case "model":
							gptCredentials.model = values.(string)
						case "prompt":
							gptCredentials.prompt = values.(string)
						case "forEachTest":
							gptCredentials.forEachTest = values.(bool)
						case "onlyFailed":
							gptCredentials.onlyFailed = values.(bool)
						}
					}

				}

			}
		}
		if gptFound {
			g.logger.Debug("GPT configuration found")
		}

	}
	g.logger.Debug("✅")
	g.logger.Debug("%v", gptCredentials)
	return gptCredentials, nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	yatasPlugin := &YatasPlugin{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	// Name of your plugin
	var pluginMap = map[string]plugin.Plugin{
		"gpt": &commons.YatasPlugin{Impl: yatasPlugin},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// Function that runs the checks or things to dot
func runPlugin(g *YatasPlugin, c *commons.Config, plugin string) ([]commons.Tests, error) {
	var checksAll []commons.Tests

	// Run the checks here
	var gptCredentials GPTPlugin
	gptCredentials, err := UnmarshalGPT(g, c)
	if err != nil {
		g.logger.Error("Error unmarshaling GPT config", "error", err)
	}
	generateReportChat(g, gptCredentials, c)

	return checksAll, nil
}
