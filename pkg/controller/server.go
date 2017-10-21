package controller

type AgentInfo struct {
}

type Server interface {
	// Run starts the server.
	Run() error

	// RegistryHook registers.
	RegistryHook() error

	// AgentHeartbeat receives the heartbeat from agents.
	AgentHeartbeat() error

	// AgentsInfo gets the agents' info.
	AgentsInfo() ([]AgentInfo, error)
}
