package config

type ConfigOptsions struct {
	Name         string
	Namespace    string
	Replicas     int
	Token        string
	ImageVersion string
	Volume       struct {
		Enable bool
		Size   int
	}
	Services struct {
		Enable bool
		Type   string
	}
	Persistence struct {
		Dir                string
		CleanOnStart       bool
		MaxQueues          int
		MaxSubscribers     int
		MaxMessages        int
		MaxSize            int
		MaxRetention       int
		MaxInactivityPurge int `name:"Max" desc:"Sets KubeMQ delete channel/queue due to inactivity time in minutes, 0 = no purging"`
	}
}
