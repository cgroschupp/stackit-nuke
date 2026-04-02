package stackit

import (
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/core/config"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"
	loadbalancer "github.com/stackitcloud/stackit-sdk-go/services/loadbalancer/v2api"
	mariadb "github.com/stackitcloud/stackit-sdk-go/services/mariadb/v1api"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"
	postgres "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v2api"
	redis "github.com/stackitcloud/stackit-sdk-go/services/redis/v1api"
	ske "github.com/stackitcloud/stackit-sdk-go/services/ske/v2api"
)

type Client struct {
	IaaS          *iaas.APIClient
	LoadBalancer  *loadbalancer.APIClient
	Ske           *ske.APIClient
	Redis         *redis.APIClient
	Postgres      *postgres.APIClient
	Mariadb       *mariadb.APIClient
	ObjectStorage *objectstorage.APIClient
	Dns           *dns.APIClient
}

func NewClient(region string) (*Client, error) {
	iaasClient, err := iaas.NewAPIClient()
	if err != nil {
		return nil, err
	}

	dnsClient, err := dns.NewAPIClient()
	if err != nil {
		return nil, err
	}
	lbClient, err := loadbalancer.NewAPIClient()
	if err != nil {
		return nil, err
	}
	skeClient, err := ske.NewAPIClient()
	if err != nil {
		return nil, err
	}
	redisClient, err := redis.NewAPIClient(config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to create redis client: %w", err)
	}
	postgresClient, err := postgres.NewAPIClient()
	if err != nil {
		return nil, err
	}
	mariadbClient, err := mariadb.NewAPIClient(config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	objectstorageClient, err := objectstorage.NewAPIClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		IaaS:          iaasClient,
		LoadBalancer:  lbClient,
		Ske:           skeClient,
		Redis:         redisClient,
		Postgres:      postgresClient,
		Mariadb:       mariadbClient,
		ObjectStorage: objectstorageClient,
		Dns:           dnsClient,
	}, nil
}
