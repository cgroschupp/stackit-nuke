package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const ObjectStorageResource = "ObjectStorage"

func init() {
	registry.Register(&registry.Registration{
		Name:     ObjectStorageResource,
		Scope:    nuke.Account,
		Resource: &ObjectStorage{},
		Lister:   &ObjectStorageLister{},
	})
}

type ObjectStorageLister struct{}

func (l *ObjectStorageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.ObjectStorage.DefaultAPI.ListBuckets(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, bucket := range resp.Buckets {
		resources = append(resources, &Server{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			Name:      &bucket.Name,
		})
	}

	return resources, nil
}

type ObjectStorage struct {
	client    *stackitclient.Client
	projectID string
	region    string
	Name      *string
}

func (r *ObjectStorage) Remove(ctx context.Context) error {
	_, err := r.client.ObjectStorage.DefaultAPI.DeleteBucket(ctx, r.projectID, r.region, deref(r.Name)).Execute()
	return err
}

func (r *ObjectStorage) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *ObjectStorage) String() string               { return safeName(r.Name, &r.projectID) }
