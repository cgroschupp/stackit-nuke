package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const ZoneResource = "Zone"

func init() {
	registry.Register(&registry.Registration{
		Name:     ZoneResource,
		Scope:    nuke.Account,
		Resource: &Zone{},
		Lister:   &ZoneListner{},
	})
}

type ZoneListner struct{}

func (l *ZoneListner) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)
	resp, err := opts.Client.Dns.DefaultAPI.ListZones(ctx, opts.ProjectID).Execute()
	if err != nil {
		return resources, err
	}

	for _, zone := range resp.Zones {
		if !*zone.Active {
			continue
		}
		resources = append(resources, &Zone{
			client:        opts.Client,
			Name:          &zone.Name,
			DnsName:       &zone.DnsName,
			ID:            &zone.Id,
			Active:        *zone.Active,
			IsReverseZone: *zone.IsReverseZone,
			// Labels:  zone.Labels,
		})
	}

	return resources, nil
}

type Zone struct {
	client        *stackitclient.Client
	ID            *string
	Name          *string
	DnsName       *string
	projectID     string
	Active        bool
	IsReverseZone bool
	// Labels    map[string]interface{}
}

func (r *Zone) Remove(ctx context.Context) error {
	_, err := r.client.Dns.DefaultAPI.DeleteZone(ctx, r.projectID, *r.ID).Execute()
	return err
}

func (r *Zone) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Zone) String() string               { return safeName(r.Name, r.ID) }
