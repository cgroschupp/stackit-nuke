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

const ServerResource = "Server"

func init() {
	registry.Register(&registry.Registration{
		Name:      ServerResource,
		Scope:     nuke.Account,
		Resource:  &Server{},
		Lister:    &ServerLister{},
		DependsOn: []string{},
	})
}

type ServerLister struct{}

func (l *ServerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.IaaS.DefaultAPI.ListServers(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, server := range resp.Items {
		resources = append(resources, &Server{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			ID:        server.Id,
			Name:      &server.Name,
			// Status:    stringPtr(string(server.Status)),
			Created: parseTimePtr(server.CreatedAt),
		})
	}

	return resources, nil
}

type Server struct {
	client    *stackitclient.Client
	projectID string
	region    string
	ID        *string
	Name      *string
	Status    *string    `description:"Current server status"`
	Created   *time.Time `description:"Creation timestamp"`
}

func (r *Server) Remove(ctx context.Context) error {
	err := r.client.IaaS.DefaultAPI.DeleteServer(ctx, r.projectID, r.region, deref(r.ID)).Execute()
	return err
}

func (r *Server) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Server) String() string               { return safeName(r.Name, r.ID) }
