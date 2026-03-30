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

const ImageResource = "Image"

func init() {
	registry.Register(&registry.Registration{
		Name:     ImageResource,
		Scope:    nuke.Account,
		Resource: &Image{},
		Lister:   &ImageLister{},
	})
}

type ImageLister struct{}

func (l *ImageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)
	resp, err := opts.Client.IaaS.DefaultAPI.ListImages(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, image := range resp.Items {
		if *image.Scope != "local" {
			continue
		}
		resources = append(resources, &Image{
			client:  opts.Client,
			Name:    &image.Name,
			ID:      image.Id,
			Created: parseTimePtr(image.CreatedAt),
			Labels:  image.Labels,
		})
	}

	return resources, nil
}

type Image struct {
	client    *stackitclient.Client
	ID        *string
	Name      *string
	projectID string
	region    string
	Created   *time.Time `description:"Creation timestamp"`
	Labels    map[string]interface{}
}

func (r *Image) Remove(ctx context.Context) error {
	err := r.client.IaaS.DefaultAPI.DeleteKeyPair(ctx, deref(r.Name)).Execute()
	return err
}

func (r *Image) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Image) String() string               { return deref(r.Name) }
