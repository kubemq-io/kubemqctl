package builder

type ClusterDTO struct {
	ID         string `json:"id"`
	Deployment struct {
		Base struct {
			ClusterName      string `json:"clusterName"`
			ClusterNamespace string `json:"clusterNamespace"`
			LicenseKey       string `json:"licenseKey"`
		} `json:"base"`
		Spec struct {
			Replicas string `json:"replicas"`
			Mode     string `json:"mode"`
		} `json:"spec"`
	} `json:"deployment"`
	Authentication *struct {
		Mode          string `json:"mode"`
		PublicKeyType string `json:"publicKeyType"`
		PublicKey     string `json:"publicKey"`
	} `json:"authentication"`
	Authorization *struct {
		Mode   string `json:"mode"`
		Policy struct {
			Rules []struct {
				ClientID    string `json:"ClientID"`
				Channel     string `json:"Channel"`
				Queues      bool   `json:"Queues"`
				Events      bool   `json:"Events"`
				EventsStore bool   `json:"EventsStore"`
				Queries     bool   `json:"Queries"`
				Commands    bool   `json:"Commands"`
				Read        bool   `json:"Read"`
				Write       bool   `json:"Write"`
			} `json:"rules"`
		} `json:"policy"`
		URL           string `json:"url"`
		FetchInterval int    `json:"fetchInterval"`
	} `json:"authorization"`
	Routing *struct {
		Mode   string `json:"mode"`
		Routes struct {
			KeyRoutes []struct {
				Key         string `json:"key"`
				Events      string `json:"events"`
				EventsStore string `json:"events_store"`
				Queues      string `json:"queues"`
			} `json:"keyRoutes"`
		} `json:"routes"`
		URL           string `json:"url"`
		FetchInterval int    `json:"fetchInterval"`
	} `json:"routing"`
	GrpcInterface *struct {
		Mode     string `json:"mode"`
		NodePort int    `json:"nodePort"`
	} `json:"grpcInterface"`
	RestInterface *struct {
		Mode     string `json:"mode"`
		NodePort int    `json:"nodePort"`
	} `json:"restInterface"`
	APIInterface *struct {
		Mode     string `json:"mode"`
		NodePort int    `json:"nodePort"`
	} `json:"apiInterface"`
	Security *struct {
		Mode string `json:"mode"`
		Cert string `json:"cert"`
		Key  string `json:"key"`
		Ca   string `json:"ca"`
	} `json:"security"`
	Image *struct {
		Image      string `json:"image"`
		PullPolicy string `json:"pullPolicy"`
	} `json:"image"`
	Volume *struct {
		Mode         string `json:"mode"`
		Size         int    `json:"size"`
		StorageClass string `json:"storageClass"`
	} `json:"volume"`
	Health *struct {
		Mode                string `json:"mode"`
		InitialDelaySeconds int    `json:"initialDelaySeconds"`
		PeriodSeconds       int    `json:"periodSeconds"`
		TimeoutSeconds      int    `json:"timeoutSeconds"`
		SuccessThreshold    int    `json:"successThreshold"`
		FailureThreshold    int    `json:"failureThreshold"`
	} `json:"health"`
	Resources *struct {
		Mode                     string `json:"mode"`
		LimitsCPU                int    `json:"limitsCpu"`
		RequestCPU               int    `json:"requestCpu"`
		LimitsMemory             int    `json:"limitsMemory"`
		RequestMemory            int    `json:"requestMemory"`
		LimitsEphemeralStorage   int    `json:"limitsEphemeralStorage"`
		RequestsEphemeralStorage int    `json:"requestsEphemeralStorage"`
	} `json:"resources"`
	Nodes *struct {
		Mode  string `json:"mode"`
		Items struct {
			Kv []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"kv"`
		} `json:"items"`
	} `json:"nodes"`
	Store *struct {
		Mode                     string `json:"mode"`
		MaxChannels              int    `json:"maxChannels"`
		MaxSubscribers           int    `json:"maxSubscribers"`
		MaxMessages              int    `json:"maxMessages"`
		MaxChannelSize           int    `json:"maxChannelSize"`
		MessagesRetentionMinutes int    `json:"messagesRetentionMinutes"`
		PurgeInactiveMinutes     int    `json:"purgeInactiveMinutes"`
	} `json:"store"`
	Queues *struct {
		Mode                      string `json:"mode"`
		MaxReceiveMessagesRequest int    `json:"maxReceiveMessagesRequest"`
		MaxReQueues               int    `json:"maxReQueues"`
		MaxExpirationSeconds      int    `json:"maxExpirationSeconds"`
		MaxDelaySeconds           int    `json:"maxDelaySeconds"`
		DefaultWaitTimeoutSeconds int    `json:"defaultWaitTimeoutSeconds"`
		MaxWaitTimeoutSeconds     int    `json:"maxWaitTimeoutSeconds"`
		DefaultVisibilitySeconds  int    `json:"defaultVisibilitySeconds"`
		MaxVisibilitySeconds      int    `json:"maxVisibilitySeconds"`
	} `json:"queues"`
}
