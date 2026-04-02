package resources

import (
	"context"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

const RecordSetResource = "RecordSet"

func init() {
	registry.Register(&registry.Registration{
		Name:      RecordSetResource,
		Scope:     nuke.Account,
		Resource:  &RecordSet{},
		Lister:    &RecordSetListner{},
		DependsOn: []string{},
	})
}

type RecordSetListner struct{}

func (l *RecordSetListner) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)
	resp, err := opts.Client.Dns.DefaultAPI.ListZones(ctx, opts.ProjectID).StateNeq("DELETE_SUCCEEDED").PageSize(100).Execute()
	if err != nil {
		return resources, err
	}

	for _, zone := range resp.Zones {
		if !*zone.Active {
			continue
		}
		respRS, err := opts.Client.Dns.DefaultAPI.ListRecordSets(ctx, opts.ProjectID, zone.Id).StateNeq("DELETE_SUCCEEDED").PageSize(100).Execute()
		if err != nil {
			return resources, err
		}
		for _, record := range respRS.RrSets {
			if record.Type == "NS" || record.Type == "SOA" {
				continue
			}
			resources = append(resources, &RecordSet{
				client:    opts.Client,
				Name:      &record.Name,
				projectID: opts.ProjectID,
				ZoneID:    &zone.Id,
				ID:        record.Id,
				Active:    *record.Active,
				Type:      record.Type,
				State:     record.State,
			})
		}
	}

	return resources, nil
}

type RecordSet struct {
	client    *stackitclient.Client
	ID        string
	Name      *string
	ZoneID    *string
	projectID string
	Type      string
	Active    bool
	State     string
}

func (r *RecordSet) Remove(ctx context.Context) error {
	_, err := r.client.Dns.DefaultAPI.DeleteRecordSet(ctx, r.projectID, *r.ZoneID, r.ID).Execute()
	return err
}

func (r *RecordSet) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *RecordSet) String() string               { return safeName(r.Name, &r.ID) }
