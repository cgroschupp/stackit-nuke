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

const PublicIPResource = "PublicIP"

func init() {
	registry.Register(&registry.Registration{
		Name:     PublicIPResource,
		Scope:    nuke.Account,
		Resource: &PublicIP{},
		Lister:   &PublicIPLister{},
	})
}

type PublicIPLister struct{}

func (l *PublicIPLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.IaaS.DefaultAPI.ListPublicIPs(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, ip := range resp.Items {
		resources = append(resources, &PublicIP{
			client:    opts.Client,
			projectID: opts.ProjectID,
			ID:        ip.Id,
			Address:   ip.Ip,
		})
	}

	return resources, nil
}

type PublicIP struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Address   *string    `description:"Allocated public IP address"`
	Created   *time.Time `description:"Creation timestamp"`
}

func (r *PublicIP) Remove(ctx context.Context) error {
	err := r.client.IaaS.DefaultAPI.DeletePublicIP(ctx, r.projectID, r.region, deref(r.ID)).Execute()
	return err
}

func (r *PublicIP) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *PublicIP) String() string               { return safeName(r.Name, r.ID) }
