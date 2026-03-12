package resources

import (
	"context"
	"time"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const NetworkResource = "Network"

func init() {
	registry.Register(&registry.Registration{
		Name:      NetworkResource,
		Scope:     nuke.Account,
		Resource:  &Network{},
		Lister:    &NetworkLister{},
		DependsOn: []string{ServerResource, LoadBalancerResource},
	})
}

type NetworkLister struct{}

func (l *NetworkLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.IaaS.DefaultAPI.ListNetworks(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, network := range resp.Items {
		resources = append(resources, &Network{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			ID:        &network.Id,
			Name:      &network.Name,
			Created:   parseTimePtr(network.CreatedAt),
		})
	}

	return resources, nil
}

type Network struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Created   *time.Time `description:"Creation timestamp"`
}

func (r *Network) Remove(ctx context.Context) error {
	err := r.client.IaaS.DefaultAPI.DeleteNetwork(ctx, r.projectID, r.region, deref(r.ID)).Execute()
	return err
}

func (r *Network) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Network) String() string               { return safeName(r.Name, r.ID) }
