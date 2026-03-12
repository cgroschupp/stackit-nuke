package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const PostgresFlexResource = "PostgresFlex"

func init() {
	registry.Register(&registry.Registration{
		Name:     PostgresFlexResource,
		Scope:    nuke.Account,
		Resource: &PostgresFlex{},
		Lister:   &PostgresFlexLister{},
	})
}

type PostgresFlexLister struct{}

func (l *PostgresFlexLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.Postgres.DefaultAPI.ListInstances(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, db := range resp.Items {
		resources = append(resources, &Server{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			ID:        db.Id,
			Name:      db.Name,
			Status:    db.Status,
		})
	}

	return resources, nil
}

type PostgresFlex struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Status    *string `description:"Current db status"`
}

func (r *PostgresFlex) Remove(ctx context.Context) error {
	err := r.client.Postgres.DefaultAPI.DeleteInstance(ctx, r.projectID, r.region, deref(r.ID)).Execute()
	return err
}

func (r *PostgresFlex) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *PostgresFlex) String() string               { return safeName(r.Name, r.ID) }
