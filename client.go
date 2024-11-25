package retrygger

import (
	"errors"
	"fmt"

	"github.com/mcku/retrygger/modules/grpc/reconpb/jobmgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Manager package for transfer jobs
// Each service supports this, should implement the interface defined below

// TODO: define this as proto
// type txnJobService interface {
// 	Trigger(moduleId, taskId string) error
// 	Restart(moduleId string)
// 	Configure(cronParams)
// }

var (
	ErrNoService = errors.New("service is nil")
	ErrBadUri    = errors.New("uri is bad")
	ErrGrpcConn  = errors.New("grpc connection error")
)

type ServiceClient struct {
	service string
	client  jobmgmt.JobmgmtAdminApiServiceClient
}

func (sc *ServiceClient) GetClient() jobmgmt.JobmgmtAdminApiServiceClient {
	return sc.client
}

// dials service @ uri
// uri might be k8s cluster service or external
func NewTxnJobClient(service string, uri string) (*ServiceClient, error) {
	if service == "" {
		return nil, ErrNoService
	}
	if uri == "" {
		return nil, fmt.Errorf("%s is bad, %w", uri, ErrBadUri)
	}

	grpcConn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, fmt.Errorf("grpc error: %s, %w", err.Error(), ErrGrpcConn)
	}
	client := jobmgmt.NewJobmgmtAdminApiServiceClient(grpcConn)

	return &ServiceClient{
		service: service,
		client:  client,
	}, nil
}
