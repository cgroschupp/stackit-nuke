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

const LoadBalancerResource = "LoadBalancer"

func init() {
	registry.Register(&registry.Registration{
		Name:     LoadBalancerResource,
		Scope:    nuke.Account,
		Resource: &LoadBalancer{},
		Lister:   &LoadBalancerLister{},
	})
}

type LoadBalancerLister struct{}

func (l *LoadBalancerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.LoadBalancer.DefaultAPI.ListLoadBalancers(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	if resp.LoadBalancers == nil {
		return resources, nil
	}

	for _, lb := range resp.LoadBalancers {
		resources = append(resources, &LoadBalancer{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			Name:      lb.Name,
			Labels:    *lb.Labels,
		})
	}

	return resources, nil
}

type LoadBalancer struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Status    *string    `description:"Current load balancer status"`
	Created   *time.Time `description:"Creation timestamp"`
	Labels    map[string]string
}

func (r *LoadBalancer) Remove(ctx context.Context) error {
	_, err := r.client.LoadBalancer.DefaultAPI.DeleteLoadBalancer(ctx, r.projectID, r.region, deref(r.ID)).Execute()
	return err
}

func (r *LoadBalancer) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *LoadBalancer) String() string               { return safeName(r.Name, r.ID) }
