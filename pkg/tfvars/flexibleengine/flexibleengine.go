package flexibleengine

import (
	"encoding/json"
)

type config struct {
	AccessKey       string `json:"fe_access_key,omitempty"`
	SecretKey       string `json:"fe_secret_key,omitempty"`
	DomainName      string `json:"fe_domain_name"`
	Region          string `json:"fe_region"`
	ProjectName     string `json:"fe_project_name,omitempty"`
	VPC             string `json:"fe_vpc,omitempty"`
	VPCName         string `json:"fe_vpc_name,omitempty"`
	VPCCIDR         string `json:"fe_vpc_cidr,omitempty"`
	Subnet          string `json:"fe_subnet,omitempty"`
	SubnetName      string `json:"fe_subnet_name,omitempty"`
	SubnetCIDR      string `json:"fe_subnet_cidr,omitempty"`
	SubnetGatewayIP string `json:"fe_subnet_gateway_ip,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	AccessKey       string
	SecretKey       string
	DomainName      string
	Region          string
	ProjectName     string
	VPC             string
	VPCName         string
	VPCCIDR         string
	Subnet          string
	SubnetName      string
	SubnetCIDR      string
	SubnetGatewayIP string
}

// TFVars generates Azure-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	cfg := &config{
		AccessKey:       sources.AccessKey,
		SecretKey:       sources.SecretKey,
		DomainName:      sources.DomainName,
		Region:          sources.Region,
		ProjectName:     sources.ProjectName,
		VPC:             sources.VPC,
		VPCName:         sources.VPCName,
		VPCCIDR:         sources.VPCCIDR,
		Subnet:          sources.Subnet,
		SubnetName:      sources.SubnetName,
		SubnetCIDR:      sources.SubnetCIDR,
		SubnetGatewayIP: sources.SubnetGatewayIP,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
