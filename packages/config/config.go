package config

import (
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type config struct {
	schedule   string
	logMode    zerolog.Level
	rjsApiUrl  string
	rjsApiKey  string
	rjsBalance uint
	rjsPeriod  uint
	disHook    string
	disBotName string
	disBotAva  string
	disBotMsg  string
}

type Config interface {
	ZerologLevel() zerolog.Level
	RajaSMSApiURL() string
	RajaSMSApiKey() string
	RajaSMSLowBalance() uint
	RajaSMSGraceDays() uint
	DishookURL() string
	DishookBotName() string
	DishookBotAvatarURL() string
	DishookBotMessage() string
	Schedule() string
}

var (
	cfg  Config
	once sync.Once
)

func (c config) Schedule() string {
	return c.schedule
}

func (c config) ZerologLevel() zerolog.Level {
	return c.logMode
}

func (c config) RajaSMSApiURL() string {
	return c.rjsApiUrl
}

func (c config) RajaSMSApiKey() string {
	return c.rjsApiKey
}

func (c config) RajaSMSLowBalance() uint {
	return c.rjsBalance
}

func (c config) RajaSMSGraceDays() uint {
	return c.rjsPeriod
}

func (c config) DishookURL() string {
	return c.disHook
}

func (c config) DishookBotName() string {
	return c.disBotName
}

func (c config) DishookBotAvatarURL() string {
	return c.disBotAva
}

func (c config) DishookBotMessage() string {
	return c.disBotMsg
}

func Get() Config {
	once.Do(func() {
		cfg = read()
	})
	return cfg
}

func read() Config {
	fang := viper.New()

	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()

	fang.SetConfigName("rajasms-account-monitor")
	fang.SetConfigType("yml")
	fang.AddConfigPath(".")

	value, available := os.LookupEnv("CONFIGDIR_PATH")
	if available {
		fang.AddConfigPath(value)
	}

	_ = fang.ReadInConfig()

	balance := fang.GetUint("rajasms.lowbalance")
	if balance == 0 {
		balance = 100000
	}

	period := fang.GetUint("rajasms.graceperiod")
	if period == 0 {
		period = 7
	}

	botMsg := fang.GetString("discord.bot.message")
	if len(strings.TrimSpace(botMsg)) == 0 {
		botMsg = "Reminder akun RajaSMS"
	}

	var logmode zerolog.Level
	switch fang.GetString("logmode") {
	case "debug":
		logmode = zerolog.DebugLevel
	case "info":
		logmode = zerolog.InfoLevel
	case "warn":
		logmode = zerolog.WarnLevel
	case "error":
		logmode = zerolog.ErrorLevel
	default:
		logmode = zerolog.Disabled
	}

	return &config{
		logMode:    logmode,
		rjsApiUrl:  strings.TrimSpace(fang.GetString("rajasms.api.url")),
		rjsApiKey:  strings.TrimSpace(fang.GetString("rajasms.api.key")),
		rjsBalance: balance,
		rjsPeriod:  period,
		disHook:    strings.TrimSpace(fang.GetString("discord.webhookurl")),
		disBotName: strings.TrimSpace(fang.GetString("discord.bot.name")),
		disBotAva:  strings.TrimSpace(fang.GetString("discord.bot.avatarurl")),
		disBotMsg:  botMsg,
	}
}
