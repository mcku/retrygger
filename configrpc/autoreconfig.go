package configrpc

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/mcku/retrygger/managedcrons"
	"github.com/robfig/cron/v3"
)

type autoReconfig struct {
	selfCron *cron.Cron
}

func NewAutoReconfig() *autoReconfig {
	return &autoReconfig{}
}

// as a cron job, auto-retrieves config for other jobs from config provider
func (s *autoReconfig) InitAutoReconfCronForJobs(cronMgr managedcrons.CronMgr) error {

	if s.selfCron != nil {
		s.selfCron.Stop()
	}
	loc, _ := time.LoadLocation("Europe/Istanbul")
	c := cron.New(cron.WithLocation(loc), cron.WithSeconds())

	// first run, autoconf and first run
	doAutoReconf(cronMgr)
	// for _, mcItem := range cronMgr.List() {
	// 	cronMgr.RunNow(mcItem.GetCron())
	// }

	// every 3 minutes on 0,3,6... minutes of the hour
	_, err := c.AddFunc("30 */1 * * * *", func() {
		doAutoReconf(cronMgr)
	})
	if err != nil {
		return fmt.Errorf("initforJobs: %w", err)
	}
	s.selfCron = c
	c.Start()
	return nil
}

func doAutoReconf(cronMgr managedcrons.CronMgr) error {

	for _, mcItem := range cronMgr.List() {

		err := cronMgr.AutoReconf(mcItem)
		if err != nil {
			glog.Infof("autoReconfCron: error: %s", err.Error())
			return err
		}

	}
	return nil
}
