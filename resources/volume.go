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

const VolumeResource = "Volume"

func init() {
	registry.Register(&registry.Registration{
		Name:      VolumeResource,
		Scope:     nuke.Account,
		Resource:  &Volume{},
		Lister:    &VolumeLister{},
		DependsOn: []string{ServerResource},
	})
}

type VolumeLister struct{}

func (l *VolumeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.IaaS.DefaultAPI.ListVolumes(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, volume := range resp.Items {
		resources = append(resources, &Volume{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			ID:        volume.Id,
			Name:      volume.Name,
			Size:      volume.Size,
			Created:   parseTimePtr(volume.CreatedAt),
			Labels:    volume.Labels,
		})
	}

	return resources, nil
}

type Volume struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Size      *int64     `description:"Volume size in GB"`
	Created   *time.Time `description:"Creation timestamp"`
	Labels    map[string]interface{}
}

func (r *Volume) Remove(ctx context.Context) error {
	err := r.client.IaaS.DefaultAPI.DeleteVolume(ctx, r.projectID, r.region, deref(r.ID)).Execute()
	return err
}

func (r *Volume) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Volume) String() string               { return safeName(r.Name, r.ID) }
