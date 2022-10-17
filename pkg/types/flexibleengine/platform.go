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

	// VPC specifies the pre-existing FE VPC where the cluster will be created.
	VPC string `json:"vpc,omitempty"`

	// VPCName specifies the new FE VPC name where the cluster will be created.
	VPCName string `json:"vpcName,omitempty"`

	// VPCCIDR specifies the new FE VPC CIDR where the cluster will be created.
	VPCCIDR string `json:"vpcCIDR,omitempty"`

	// Subnet specifies the pre-existing FE Subnet where the cluster will be created.
	Subnet string `json:"subnet,omitempty"`

	// SubnetName specifies the new FE Subnet name where the cluster will be created.
	SubnetName string `json:"subnetName,omitempty"`

	// SubnetCIDR specifies the new FE Subnet CIDR where the cluster will be created.
	SubnetCIDR string `json:"subnetCIDR,omitempty"`

	// SubnetGatewayIP specifies the new FE Subnet gateway IP where the cluster will be created.
	SubnetGatewayIP string `json:"subnetGatewayIP,omitempty"`
}
