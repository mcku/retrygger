package retrygger

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mcku/retrygger/logrpc"
	"github.com/mcku/retrygger/managedcrons"
	"github.com/mcku/retrygger/modules/grpc/reconpb/jobmgmt"
)

type cronMgr interface {
	GetCron(cronName string) *managedcrons.MC
}
type txnJobServer struct {
	service      string
	logSvcAddr   string
	cronMgr      cronMgr
	buildVersion string
}

// GetStatus implements jobmgmt.JobmgmtAdminApiServiceServer.
func (s *txnJobServer) GetStatus(context.Context, *jobmgmt.GetStatusRequest) (*jobmgmt.GetStatusResponse, error) {
	return nil, fmt.Errorf("get status: not implemented")
}

// TriggerJob implements jobmgmt.JobmgmtAdminApiServiceServer.
func (s *txnJobServer) TriggerJob(ctx context.Context, req *jobmgmt.TriggerJobRequest,
) (*jobmgmt.TriggerJobResponse, error) {
	if req.Service == "" || req.Job == "" {
		return nil, fmt.Errorf("invalid job/service")
	}
	if req.Service != s.service {
		return nil, fmt.Errorf("triggerJob: services do not match: incoming %s vs %s", req.Service, s.service)
	}
	// lookup cron jobs for job
	cronEntry := s.cronMgr.GetCron(req.Job)
	if cronEntry == nil {
		return nil, fmt.Errorf("triggerJob: cron entry is nil for %s", s.service)
	}
	trigger := cronEntry.GetTrigger()
	if trigger == nil {
		return nil, fmt.Errorf("triggerJob: trigger is nil for %s", s.service)
	}
	log, err := trigger()

	if s.logSvcAddr != "" {

		status := jobmgmt.LogRecord_STATUS_SUCCESS
		if err != nil {
			status = jobmgmt.LogRecord_STATUS_FAILED
		}
		logRecord := &jobmgmt.LogRecord{
			RecordId:     uuid.NewString(),
			Timestamp:    time.Now().UnixNano(),
			Message:      log,
			Service:      req.Service,
			Job:          req.Job,
			Status:       status,
			Params:       req.Params,
			AckStatus:    false,
			BuildVersion: s.buildVersion,
			Initiator:    "manual",
		}
		logrpc.WriteLogRpc(ctx, s.logSvcAddr, []*jobmgmt.LogRecord{
			logRecord,
		})

	}

	if err != nil {
		return &jobmgmt.TriggerJobResponse{
			Status:  jobmgmt.JobStatus_JOB_STATUS_HAS_ERRORS,
			Message: fmt.Sprintf("Error: %s\n Log: %s", err, log),
		}, nil
	}
	return &jobmgmt.TriggerJobResponse{
		Status:  jobmgmt.JobStatus_JOB_STATUS_SUCCESS,
		Message: log,
	}, nil
}

func NewTxnJobServer(
	service string,
	logSvcAddr string,
	cronMgr cronMgr,
	buildVersion string,
) *txnJobServer {

	return &txnJobServer{
		service:      service,
		logSvcAddr:   logSvcAddr,
		cronMgr:      cronMgr,
		buildVersion: buildVersion,
	}
}
