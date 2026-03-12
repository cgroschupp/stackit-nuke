package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const RedisResource = "Redis"

func init() {
	registry.Register(&registry.Registration{
		Name:      RedisResource,
		Scope:     nuke.Account,
		Resource:  &Redis{},
		Lister:    &RedisLister{},
		DependsOn: []string{VolumeResource},
	})
}

type RedisLister struct{}

func (l *RedisLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.Redis.DefaultAPI.ListInstances(ctx, opts.ProjectID).Execute()
	if err != nil {
		return resources, err
	}

	for _, redis := range resp.Instances {
		resources = append(resources, &Server{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			ID:        redis.InstanceId,
			Name:      &redis.Name,
			Status:    redis.Status,
		})
	}

	return resources, nil
}

type Redis struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Status    *string `description:"Current redis status"`
}

func (r *Redis) Remove(ctx context.Context) error {
	err := r.client.Redis.DefaultAPI.DeleteInstance(ctx, r.projectID, deref(r.ID)).Execute()
	return err
}

func (r *Redis) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Redis) String() string               { return safeName(r.Name, r.ID) }
