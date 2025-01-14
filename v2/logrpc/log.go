package logrpc

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/mcku/retrygger/v2/modules/grpc/reconpb/jobmgmt"
	"google.golang.org/grpc"
)

// WriteLogRpc - sends log to providerAddr
func WriteLogRpc(ctx context.Context, logSvcAddr string, logs []*jobmgmt.LogRecord) error {
	if logs == nil {
		return fmt.Errorf("writeLog: empty logs")
	}
	// gCli, err := grpc.NewClient(providerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.)

	// we still need localhost:port access for testing
	gCli, err := grpc.Dial(logSvcAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("writeLog: grpc client %s", err.Error())
	}
	defer gCli.Close()
	cli := jobmgmt.NewJobmgmtLogWriterServiceClient(gCli)

	req := &jobmgmt.CreateLogRequest{
		Logs: logs,
	}
	_, err = cli.CreateLog(ctx, req)
	if err != nil {
		glog.Warningf("writeLog: read from %s err: %s, req: %v", logSvcAddr, err.Error(), req)
		return fmt.Errorf("writeLog: read err: %s, req: %v", err.Error(), req)
	}
	return nil
}
