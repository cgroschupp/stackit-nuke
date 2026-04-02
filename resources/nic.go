package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const NicResource = "Nic"

func init() {
	registry.Register(&registry.Registration{
		Name:      NicResource,
		Scope:     nuke.Account,
		Resource:  &Nic{},
		Lister:    &NicListner{},
		DependsOn: []string{},
	})
}

type NicListner struct{}

func (l *NicListner) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)
	networks, err := opts.Client.IaaS.DefaultAPI.ListNetworks(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}
	for _, network := range networks.Items {
		resp, err := opts.Client.IaaS.DefaultAPI.ListNics(ctx, opts.ProjectID, opts.Region, network.Id).Execute()
		if err != nil {
			return resources, err
		}

		for _, nic := range resp.Items {
			resources = append(resources, &Nic{
				client:    opts.Client,
				projectID: opts.ProjectID,
				region:    opts.Region,
				networkID: *nic.NetworkId,
				ID:        nic.Id,
				Name:      nic.Name,
			})
		}
	}

	return resources, nil
}

type Nic struct {
	client    *stackitclient.Client
	projectID string
	networkID string
	region    string
	ID        *string
	Name      *string
}

func (r *Nic) Remove(ctx context.Context) error {
	err := r.client.IaaS.DefaultAPI.DeleteNic(ctx, r.projectID, r.region, r.networkID, deref(r.ID)).Execute()
	return err
}

func (r *Nic) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Nic) String() string               { return safeName(r.Name, r.ID) }
