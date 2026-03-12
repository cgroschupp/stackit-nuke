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

const KubernetesResource = "Kubernetes"

func init() {
	registry.Register(&registry.Registration{
		Name:      KubernetesResource,
		Scope:     nuke.Account,
		Resource:  &Volume{},
		Lister:    &KubernetesLister{},
		DependsOn: []string{ServerResource},
	})
}

type KubernetesLister struct{}

func (l *KubernetesLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	resources := make([]resource.Resource, 0)

	resp, err := opts.Client.Ske.DefaultAPI.ListClusters(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return resources, err
	}

	for _, cluster := range resp.Items {
		resources = append(resources, &Volume{
			client:    opts.Client,
			projectID: opts.ProjectID,
			region:    opts.Region,
			Name:      cluster.Name,
			Created:   cluster.Status.CreationTime,
			// Labels:    cluster.Labels,
		})
	}

	return resources, nil
}

type Kubernetes struct {
	client    *stackitclient.Client
	projectID string
	region    string
	Name      *string
	Created   *time.Time `description:"Creation timestamp"`
}

func (r *Kubernetes) Remove(ctx context.Context) error {
	_, err := r.client.Ske.DefaultAPI.DeleteCluster(ctx, r.projectID, r.region, deref(r.Name)).Execute()
	return err
}

func (r *Kubernetes) Properties() types.Properties { return types.NewPropertiesFromStruct(r) }
func (r *Kubernetes) String() string               { return safeName(r.Name, &r.projectID) }
