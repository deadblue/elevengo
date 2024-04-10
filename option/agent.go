package option

import "github.com/deadblue/elevengo/plugin"

type AgentOption interface {
	isAgentOption()
}

type AgentCooldownOption struct {
	// Minimum cooldown duration in millisecond
	Min uint
	// Maximum cooldown duration in millisecond
	Max uint
}

func (o AgentCooldownOption) isAgentOption() {}

// AgentHttpOption allows developer customize underlying HTTP client.
type AgentHttpOption struct {
	Client plugin.HttpClient
}

func (o *AgentHttpOption) isAgentOption() {}

type AgentNameOption string

func (o AgentNameOption) isAgentOption() {}

// AgentVersionOption allows developer lock app version for agent.
type AgentVersionOption string

func (o AgentVersionOption) isAgentOption() {}
