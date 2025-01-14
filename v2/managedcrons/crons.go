package managedcrons

import (
	"sync"

	"github.com/mcku/retrygger/v2/modules/grpc/reconpb/jobmgmt"
	"github.com/robfig/cron/v3"
)

type txnImporter interface {
	Fetch() (string, error)
}
type managedCron interface {
	GetTrigger() func() (string, error)
	SetTrigger(trigger func() (string, error))
	GetAllCrons() map[string]cron.EntryID
	SetCron(c *cron.Cron)
	StopCrons()
	GetConfigFetcher() func() (*jobmgmt.JobConfig, error)
	GetCron() cron.EntryID
	GetProviderAddr() string
	GetName() string
}
type logFunc func(log string, status jobmgmt.LogRecord_Status, initiator, runtimeParams string) error
type MC struct {
	lock          sync.RWMutex
	cronName      string
	providerAddr  string
	cronEntry     cron.EntryID
	trigger       func(paramStr string) (string, error) // returns an execution log, along with an error
	configFetcher func() (*jobmgmt.JobConfig, error)
	logWriter     logFunc
	currentConfig *jobmgmt.JobConfig
	// each cron can run with its own param set. caller needs to specify this when initializing
	runtimeParamBuilder func() string
}

func NewManagedCron(
	cronName string,
	trigger func(paramStr string) (string, error),
	configFetcher func() (*jobmgmt.JobConfig, error),
	currentConfig *jobmgmt.JobConfig,
	logWriter logFunc,
	providerAddr string,
	runtimeParamBuilder func() string,
) *MC {
	return &MC{
		cronName:            cronName,
		trigger:             trigger,
		configFetcher:       configFetcher,
		currentConfig:       currentConfig,
		logWriter:           logWriter,
		providerAddr:        providerAddr,
		runtimeParamBuilder: runtimeParamBuilder,
	}
}

func (s *MC) GetLogWriter() logFunc {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.logWriter
}

func (s *MC) GetTrigger() func(paramStr string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s == nil {
		return nil
	}
	return s.trigger
}

func (s *MC) GetConfigFetcher() func() (*jobmgmt.JobConfig, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.configFetcher
}
func (s *MC) GetCron() cron.EntryID {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.cronEntry
}
func (s *MC) SetCron(c cron.EntryID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cronEntry = c
}

func (s *MC) GetConfig() *jobmgmt.JobConfig {
	return s.currentConfig
}

func (s *MC) SetConfig(c *jobmgmt.JobConfig) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.currentConfig = c
	return nil
}

func (s *MC) GetName() string {
	return s.cronName
}
