package logrpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcku/retrygger/v2/modules/grpc/reconpb/jobmgmt"
)

type logFunc func(log string, status jobmgmt.LogRecord_Status, initiator, runtimeParams string) error

func BuildRpcLogWriter(serviceName, jobName, buildVersion, logWriterAddr string) logFunc {
	return func(log string, status jobmgmt.LogRecord_Status, initiator, runtimeParams string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logRecord := &jobmgmt.LogRecord{
			RecordId:     uuid.NewString(),
			Timestamp:    time.Now().UnixNano(),
			Message:      log,
			Service:      serviceName,
			Job:          jobName,
			Status:       status,
			Params:       runtimeParams,
			AckStatus:    false,
			BuildVersion: buildVersion,
			Initiator:    initiator,
		}
		return WriteLogRpc(ctx, logWriterAddr, []*jobmgmt.LogRecord{logRecord})
	}
}
