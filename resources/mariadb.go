package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const MariadbResource = "Mariadb"

func init() {
	registry.Register(&registry.Registration{
		Name:     MariadbResource,
		Scope:    nuke.Account,
		Resource: &Mariadb{},
		Lister:   &MariadbLister{},
	})
}

type MariadbLister struct{}

func (l *MariadbLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.Mariadb.DefaultAPI.ListInstances(ctx, opts.ProjectID).Execute()
	if err != nil {
		return resources, err
	}

	for _, db := range resp.Instances {
		resources = append(resources, &Server{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			ID:        db.InstanceId,
			Name:      &db.Name,
			Status:    db.Status,
		})
	}

	return resources, nil
}

type Mariadb struct {
	client    *stackitclient.Client
	projectID string
	ID        *string
	Name      *string
	Status    *string `description:"Current db status"`
}

func (r *Mariadb) Remove(ctx context.Context) error {
	err := r.client.Mariadb.DefaultAPI.DeleteInstance(ctx, r.projectID, deref(r.ID)).Execute()
	return err
}

func (r *Mariadb) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Mariadb) String() string               { return safeName(r.Name, r.ID) }
