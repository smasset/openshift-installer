package vpc

import (
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

func ResourceNetworkingVIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingVIPAssociateV2Create,
		Update: resourceNetworkingVIPAssociateV2Update,
		Read:   resourceNetworkingVIPAssociateV2Read,
		Delete: resourceNetworkingVIPAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceNetworkingVIPAssociateV2Import,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vip_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingPortIDs(d *schema.ResourceData) []string {
	rawPortIDs := d.Get("port_ids").(*schema.Set).List()
	portids := make([]string, len(rawPortIDs))
	for i, raw := range rawPortIDs {
		portids[i] = raw.(string)
	}
	return portids
}

func updateNetworkingVIPAssociate(client *golangsdk.ServiceClient, vipID string, portIDs []string) error {
	allAddrs := make([]string, len(portIDs))

	// check the port id
	for i, portid := range portIDs {
		port, err := ports.Get(client, portid).Extract()
		if err != nil {
			return fmt.Errorf("error fetching port %s: %s", portid, err)
		}

		if len(port.FixedIPs) > 0 {
			allAddrs[i] = port.FixedIPs[0].IPAddress
		} else {
			return fmt.Errorf("port %s has no ip address, Error associate it", portid)
		}
	}

	// construct allowed address pairs
	allowedPairs := make([]ports.AddressPair, len(allAddrs))
	for i, addr := range allAddrs {
		allowedPairs[i] = ports.AddressPair{
			IPAddress: addr,
		}
	}
	// associate vip to port
	associateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
	}
	log.Printf("[DEBUG] VIP Associate %s with options: %#v", vipID, associateOpts)
	_, err := ports.Update(client, vipID, associateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error associate vip: %s", err)
	}

	// Update the allowed-address-pairs of the port to 1.1.1.1/0
	// to disable the source/destination check
	portpairs := make([]ports.AddressPair, 1)
	portpairs[0] = ports.AddressPair{
		IPAddress: "1.1.1.1/0",
	}
	portUpdateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &portpairs,
	}

	for _, portid := range portIDs {
		_, err = ports.Update(client, portid, portUpdateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error update port %s: %s", portid, err)
		}
	}

	return nil
}

func resourceNetworkingVIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	// check the vip port
	vipID := d.Get("vip_id").(string)
	_, err = ports.Get(networkingClient, vipID).Extract()
	if err != nil {
		return fmt.Errorf("error fetching vip %s: %s", vipID, err)
	}

	portids := resourceNetworkingPortIDs(d)
	if err = updateNetworkingVIPAssociate(networkingClient, vipID, portids); err != nil {
		return err
	}

	// set id
	d.SetId(hashcode.Strings(portids))
	return resourceNetworkingVIPAssociateV2Read(d, meta)
}

func resourceNetworkingVIPAssociateV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	// check the vip port
	vipID := d.Get("vip_id").(string)
	_, err = ports.Get(networkingClient, vipID).Extract()
	if err != nil {
		return fmt.Errorf("error fetching vip %s: %s", vipID, err)
	}

	portids := resourceNetworkingPortIDs(d)
	if err = updateNetworkingVIPAssociate(networkingClient, vipID, portids); err != nil {
		return err
	}

	return resourceNetworkingVIPAssociateV2Read(d, meta)
}

func resourceNetworkingVIPAssociateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	// check the vip port
	vipID := d.Get("vip_id").(string)
	vip, err := ports.Get(networkingClient, vipID).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "vip")
	}

	var allPorts []string
	var allAddrs []string
	// check the port still exists
	portids := resourceNetworkingPortIDs(d)
	for _, portid := range portids {
		p, err := ports.Get(networkingClient, portid).Extract()
		if err != nil {
			log.Printf("[WARN] failed to fetch port %s: %s", portid, err)
			continue
		}

		for _, ip := range p.FixedIPs {
			for _, addresspair := range vip.AllowedAddressPairs {
				if ip.IPAddress == addresspair.IPAddress {
					allPorts = append(allPorts, portid)
					allAddrs = append(allAddrs, ip.IPAddress)
					break
				}
			}
		}
	}

	// if no port is associated
	if len(allPorts) == 0 {
		log.Printf("[WARN] no port is associated with vip %s", vipID)
		d.SetId("")
		return nil
	}

	// Set the attributes pulled from the composed resource ID
	d.Set("vip_id", vipID)
	d.Set("vip_subnet_id", vip.FixedIPs[0].SubnetID)
	d.Set("vip_ip_address", vip.FixedIPs[0].IPAddress)
	d.Set("port_ids", allPorts)
	d.Set("ip_addresses", allAddrs)

	return nil
}

func resourceNetworkingVIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	// check the vip port
	vipID := d.Get("vip_id").(string)
	_, err = ports.Get(networkingClient, vipID).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "vip")
	}

	// disassociate all allowed address pairs
	allowedPairs := make([]ports.AddressPair, 0)
	disassociateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
	}
	log.Printf("[DEBUG] Disassociate all ports with %s", vipID)
	_, err = ports.Update(networkingClient, vipID, disassociateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error disassociate vip: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceNetworkingVIPAssociateV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <vip_id>/<port_id>," +
			" and at least 1 port_id must be provided")
	}

	portIds := parts[1:]
	d.SetId(hashcode.Strings(portIds))
	d.Set("vip_id", parts[0])
	d.Set("port_ids", portIds)

	return []*schema.ResourceData{d}, nil
}
