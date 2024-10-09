package agentUtils

func (a *UserAgentDetail) Lock() {
	a.mutex.Lock()
}

func (a *UserAgentDetail) Unlock() {
	a.mutex.Unlock()
}
