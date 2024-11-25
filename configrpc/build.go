package configrpc

import (
	"context"
	"time"

	"github.com/mcku/retrygger/modules/grpc/reconpb/jobmgmt"
)

func BuildRpcConfigReader(providerAddr, serviceName, jobName string) func() (*jobmgmt.JobConfig, error) {

	return func() (*jobmgmt.JobConfig, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return ReadConfigRpc(ctx, providerAddr, serviceName, jobName)
	}
}
