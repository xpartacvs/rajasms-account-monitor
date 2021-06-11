package worker

import (
	"rajasms-account-monitor/packages/config"
	"rajasms-account-monitor/packages/logger"
	"rajasms-account-monitor/packages/webhook"
	"regexp"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/xpartacvs/go-rajasms"
)

var (
	once   sync.Once
	client rajasms.Client
)

func Start() error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	schedule := config.Get().Schedule()
	rgxCron := regexp.MustCompile(`(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	if !rgxCron.MatchString(schedule) {
		schedule = "0 0 * * *"
		logger.Log().Warn().Msg("Invalid schedule format. Retain to default value")
	}

	s := gocron.NewScheduler(loc)
	_, err = s.Cron(schedule).Do(checkAccount)
	if err != nil {
		return err
	}
	s.StartBlocking()

	return nil
}

func getClient() rajasms.Client {
	once.Do(func() {
		client = rajasms.NewClient(
			config.Get().RajaSMSApiURL(),
			config.Get().RajaSMSApiKey(),
		)
	})
	return client
}

func checkAccount() error {
	i, err := getClient().GetInquiry()
	if err != nil {
		logger.Log().Error().Msg(err.Error())
		return err
	}

	if (i.GetBalance() <= config.Get().RajaSMSLowBalance()) || (uint(time.Until(i.GetExpiry()).Hours()/24) <= config.Get().RajaSMSGraceDays()) {
		if err := webhook.GetInstance().AddReminder(
			config.Get().RajaSMSLowBalance(),
			i.GetBalance(),
			config.Get().RajaSMSGraceDays(),
			i.GetExpiry(),
		).Send(config.Get().DishookURL()); err != nil {
			logger.Log().Error().Msg(err.Error())
		}
	}

	return nil
}
