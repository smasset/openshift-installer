package flexibleengine

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// AccessKey specifies the access key used to connect to FE APIs.
	AccessKey string `json:"accessKey"`

	// SecretKey specifies the secret key used to connect to FE APIs.
	SecretKey string `json:"secretKey"`

	// DomainName specifies the FE domain where the cluster will be created.
	DomainName string `json:"domain"`

	// Region specifies the FE region where the cluster will be created.
	Region string `json:"region"`

	// ProjectName specifies the FE project where the cluster will be created.
	ProjectName string `json:"project,omitempty"`
}
