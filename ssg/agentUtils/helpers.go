package agentUtils

import (
	"fmt"
	"sync"
)

// GetDefaultUserAgents returns a list of default user agents
// that can be used for making requests
func GetDefaultUserAgents() []*UserAgentDetail {
	return []*UserAgentDetail{
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
			SecChUa:   "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"",
			Platform:  "\"Windows\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/536.36",
			SecChUa:   "\"Google Chrome\";v=\"116\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"116\"",
			Platform:  "\"Windows\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 15) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.555.100 Mobile Safari/537.36",
			SecChUa:   "\"Google Chrome\";v=\"125\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"125\"",
			Platform:  "\"Android\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/532.36",
			SecChUa:   "\"Google Chrome\";v=\"114\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"",
			Platform:  "\"Windows\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/531.36",
			SecChUa:   "\"Google Chrome\";v=\"113\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"113\"",
			Platform:  "\"Windows\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/520.36",
			SecChUa:   "\"Google Chrome\";v=\"112\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"112\"",
			Platform:  "\"Windows\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/510.10",
			SecChUa:   "\"Google Chrome\";v=\"111\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"111\"",
			Platform:  "\"Windows\"",
			mutex:     &sync.Mutex{},
		},
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 15) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.6668.100 Mobile Safari/537.36",
			SecChUa:   "\"Google Chrome\";v=\"129\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"129\"",
			Platform:  "\"Android\"",
			mutex:     &sync.Mutex{},
		},
	}
}

// GetAndroidUserAgents dynamically generates a list of UserAgentDetail based on the count provided
func GetAndroidUserAgents(count int) []*UserAgentDetail {
	baseUserAgents := []UserAgentDetail{
		{
			UserAgent: "Mozilla/5.0 (Linux; Android 15) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.555.100 Mobile Safari/537.36",
			SecChUa:   "\"Google Chrome\";v=\"%d\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"%d\"",
			Platform:  "\"Android\"",
		},
		{
			UserAgent: "Mozilla/5.0 (Android 10; Mobile; rv:62.0) Gecko/68.0 Firefox/%d.0",
			SecChUa:   "\"Firefox\";v=\"%d\", \"Not;A=Brand\";v=\"8\", \"Gecko\";v=\"68\"",
			Platform:  "\"Android\"",
		},
	}

	var userAgents []*UserAgentDetail
	chromeVersion := MinChromeVersion   // Starting version for Chrome
	firefoxVersion := MinFirefoxVersion // Starting version for Firefox

	for i := 0; i < count; i++ {
		// Alternate between Chrome and Firefox user agents
		if i%2 == 0 {
			userAgents = append(userAgents, &UserAgentDetail{
				UserAgent: fmt.Sprintf(baseUserAgents[0].UserAgent, chromeVersion),
				SecChUa:   fmt.Sprintf(baseUserAgents[0].SecChUa, chromeVersion, chromeVersion),
				Platform:  baseUserAgents[0].Platform,
				mutex:     &sync.Mutex{},
			})
			chromeVersion++
		} else {
			userAgents = append(userAgents, &UserAgentDetail{
				UserAgent: fmt.Sprintf(baseUserAgents[1].UserAgent, firefoxVersion),
				SecChUa:   fmt.Sprintf(baseUserAgents[1].SecChUa, firefoxVersion, firefoxVersion),
				Platform:  baseUserAgents[1].Platform,
				mutex:     &sync.Mutex{},
			})
			firefoxVersion++
		}
	}

	return userAgents
}
