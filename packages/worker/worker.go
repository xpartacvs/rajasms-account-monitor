package worker

import (
	"rajasms-account-monitor/packages/config"
	"rajasms-account-monitor/packages/logger"
	"rajasms-account-monitor/packages/webhook"
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

	s := gocron.NewScheduler(loc)
	_, err = s.Cron(config.Get().RajaSMSApiKey()).Do(checkAccount)
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
