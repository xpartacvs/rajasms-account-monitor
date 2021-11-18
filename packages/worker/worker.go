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
	client *rajasms.Client
)

func Start() error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	s := gocron.NewScheduler(loc)
	_, err = s.Cron(config.Get().Schedule()).Do(checkAccount)
	if err != nil {
		return err
	}
	s.StartBlocking()

	return nil
}

func getClient() *rajasms.Client {
	once.Do(func() {
		var err error
		client, err = rajasms.NewCient(config.Get().RajaSMSApiURL(), config.Get().RajaSMSApiKey())
		if err != nil {
			logger.Log().Fatal().Msg("Unable to create RajaSMS Client")
		}
	})
	return client
}

func checkAccount() error {
	i, err := getClient().AccountInfo()
	if err != nil {
		logger.Log().Err(err).Msg("Cannot get account inquiry result")
		return err
	}

	if (i.Balance <= config.Get().RajaSMSLowBalance()) || (uint(time.Until(i.Expiry).Hours()/24) <= config.Get().RajaSMSGraceDays()) {
		if err := webhook.GetInstance().AddReminder(
			uint(config.Get().RajaSMSLowBalance()),
			uint(i.Balance),
			config.Get().RajaSMSGraceDays(),
			i.Expiry,
		).Send(config.Get().DishookURL()); err != nil {
			logger.Log().Err(err).Msg("Error while sending alert to discord channel")
		}
	}

	return nil
}
