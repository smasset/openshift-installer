package flexibleengine

// Metadata contains FE metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
	DomainName  string `json:"domain"`
	Region      string `json:"region"`
	ProjectName string `json:"project"`
}
