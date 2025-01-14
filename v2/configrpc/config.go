package configrpc

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/mcku/retrygger/v2/modules/grpc/reconpb/jobmgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ReadConfigRpc(ctx context.Context, providerAddr, service, job string) (*jobmgmt.JobConfig, error) {
	gCli, err := grpc.Dial(providerAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) // grpc.WithChainUnaryInterceptor(UnaryClientLoggingInterceptor)

	if err != nil {
		return nil, fmt.Errorf("readConfig: grpc client %s", err.Error())
	}
	defer gCli.Close()
	cli := jobmgmt.NewJobmgmtConfigProviderServiceClient(gCli)

	req := &jobmgmt.ReadConfigRequest{
		Job:     job,
		Service: service,
	}
	resp, err := cli.ReadConfig(ctx, req)
	if err != nil {
		glog.Warningf("readConfig: read from %s err: %s, req: %v", providerAddr, err.Error(), req)
		return nil, fmt.Errorf("readConfig: read err: %s, req: %v", err.Error(), req)
	}
	return resp.Config, nil
}
