package option

import "github.com/deadblue/elevengo/plugin"

type AgentOptions struct {
	// Underlying HTTP client which is used to perform HTTP request.
	HttpClient plugin.HttpClient

	// Minimum delay in milliseconds after last API calling.
	CooldownMinMs uint

	// Maximum delay in milliseconds after last API calling.
	CooldownMaxMs uint

	// Custom user-agent.
	Name string

	// Custom app version.
	Version string
}

func (o *AgentOptions) WithHttpClient(hc plugin.HttpClient) *AgentOptions {
	o.HttpClient = hc
	return o
}

func (o *AgentOptions) WithCooldown(minMs, maxMs uint) *AgentOptions {
	o.CooldownMinMs = minMs
	o.CooldownMaxMs = maxMs
	return o
}

func (o *AgentOptions) WithName(name string) *AgentOptions {
	o.Name = name
	return o
}
func (o *AgentOptions) WithVersion(version string) *AgentOptions {
	o.Version = version
	return o
}

func Agent() *AgentOptions {
	return &AgentOptions{}
}
