package agentUtils

import "sync"

// UserAgentDetail is a struct that holds the details of a user agent
// It has a mutex to lock and unlock the struct.
type UserAgentDetail struct {
	UserAgent string
	SecChUa   string
	Platform  string
	Browser   string
	Device    string
	OS        string
	Engine    string

	mutex *sync.Mutex
}
