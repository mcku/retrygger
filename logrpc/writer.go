package logrpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcku/retrygger/modules/grpc/reconpb/jobmgmt"
)

func BuildRpcLogWriter(serviceName, jobName, buildVersion, logWriterAddr string) func(string, jobmgmt.LogRecord_Status, string) error {
	return func(log string, status jobmgmt.LogRecord_Status, initiator string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logRecord := &jobmgmt.LogRecord{
			RecordId:     uuid.NewString(),
			Timestamp:    time.Now().UnixNano(),
			Message:      log,
			Service:      serviceName,
			Job:          jobName,
			Status:       status,
			Params:       "",
			AckStatus:    false,
			BuildVersion: buildVersion,
			Initiator:    initiator,
		}
		return WriteLogRpc(ctx, logWriterAddr, []*jobmgmt.LogRecord{logRecord})
	}
}
