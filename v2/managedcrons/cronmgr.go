package managedcrons

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/mcku/retrygger/v2/modules/grpc/reconpb/jobmgmt"
	"github.com/robfig/cron/v3"
)

type CronMgr interface {
	GetEngine() *cron.Cron
	GetCron(cronName string) *MC
	SetCron(cronName string, cron *MC)
	StartEngine()
	// StartCron(cronEntry cron.EntryID)
	// StopCron(cronEntry cron.EntryID)
	RemoveCron(cronEntry cron.EntryID)
	List() []*MC
	AutoReconf(mcItem *MC) error
	RunNow(cronEntry cron.EntryID)
}

type cronManager struct {
	lock       sync.RWMutex
	cronEngine *cron.Cron
	crons      map[string]*MC
}

func NewCronManager() *cronManager {

	loc, _ := time.LoadLocation("Europe/Istanbul")
	c := cron.New(cron.WithLocation(loc), cron.WithSeconds())
	return &cronManager{
		cronEngine: c,
		crons:      make(map[string]*MC),
	}
}

func (s *cronManager) GetCron(cronName string) *MC {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.crons[cronName]
}

func (s *cronManager) SetCron(cronName string, cron *MC) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.crons[cronName] = cron
}

func (s *cronManager) List() (list []*MC) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, c := range s.crons {
		list = append(list, c)
	}
	return list
}

func (s *cronManager) AutoReconf(mcItem *MC) error {
	engine := s.cronEngine
	glog.Infof("autoReconf: fetching config for job %s", mcItem.GetName())
	jc, err := mcItem.GetConfigFetcher()()
	if err != nil {
		return fmt.Errorf("autoReconf: %s", err.Error())
	}
	glog.Infof("autoReconf: job %s isEnabled: %s schedule: %s", mcItem.GetName(), jc.Enabled, jc.CronSchedule)
	//  mcItem.SetConfig(jc)

	if jc.Enabled != jobmgmt.YesNoStatus_YES_NO_STATUS_YES {
		// disabled remotely
		if mcItem.currentConfig == nil ||
			mcItem.currentConfig.Enabled != jobmgmt.YesNoStatus_YES_NO_STATUS_YES {
			// do nothing, disabled both locally and remotely
			mcItem.SetConfig(jc)
			return nil
		}
		// remote disabled, local enabled -> disable locally
		glog.Infof("autoReconf: DISABLING job %s isEnabled: %s schedule: %s", mcItem.GetName(), jc.Enabled, jc.CronSchedule)
		jobCron := mcItem.GetCron()
		if jobCron != 0 {
			// glog.Infof("autoReconf: stopping cron %s", jc.Job)
			engine.Remove(jobCron)
		}
		mcItem.SetConfig(jc)
		return nil
	}
	// remote enabled
	if mcItem.currentConfig == nil ||
		mcItem.currentConfig.Enabled != jobmgmt.YesNoStatus_YES_NO_STATUS_YES {
		// remote enabled, local disabled -> enable locally (with up to date remote params)
		glog.Infof("autoReconf: ENABLING job %s isEnabled: %s schedule: %s", mcItem.GetName(), jc.Enabled, jc.CronSchedule)
		mcItem.SetConfig(jc)
		err := addToCron(engine, mcItem)
		if err != nil {
			return fmt.Errorf("autoReconf: %s", err.Error())
		}
		return nil
	}
	// remote & local enabled: do nothing but check the schedule and params
	mcc := mcItem.currentConfig
	configChanged := mcc.ConfigDbParams != jc.ConfigDbParams
	scheduleChanged := mcc.CronSchedule != jc.CronSchedule
	if configChanged || scheduleChanged {

		if configChanged {
			glog.Infof("autoReconf: PARAMS updated for job %s isEnabled: %s schedule: %s", mcItem.GetName(), jc.Enabled, jc.CronSchedule)
		}
		if scheduleChanged {
			glog.Infof("autoReconf: SCHEDULE updated for job %s isEnabled: %s schedule: %s", mcItem.GetName(), jc.Enabled, jc.CronSchedule)
		}

		engine.Remove(mcItem.GetCron())
		mcItem.SetConfig(jc)
		addToCron(engine, mcItem)
		return nil

	}
	return nil

}
func addToCron(engine *cron.Cron, mcItem *MC) error {
	cronEntryID, err := engine.AddFunc(mcItem.currentConfig.CronSchedule, func() {
		log, err := mcItem.trigger("")
		if err != nil {
			mcItem.logWriter(log, jobmgmt.LogRecord_STATUS_FAILED, "cron",
				"")
			return
		}
		mcItem.logWriter(log, jobmgmt.LogRecord_STATUS_SUCCESS, "cron",
			"")
	})
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	mcItem.SetCron(cronEntryID)
	return nil
}
func (s *cronManager) Initialize() error {
	// s.StopCron()
	for _, mc := range s.crons {
		s.AutoReconf(mc)
	}
	return nil
}

func (s *cronManager) GetEngine() *cron.Cron {
	return s.cronEngine
}
func (s *cronManager) RemoveCron(cronEntry cron.EntryID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cronEngine.Remove(cronEntry)
}

func (s *cronManager) RunNow(cronEntry cron.EntryID) {
	s.cronEngine.Entry(cronEntry).WrappedJob.Run()
}

func (s *cronManager) StartEngine() {
	s.cronEngine.Start()
}
