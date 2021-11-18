package webhook

import (
	"errors"
	"rajasms-account-monitor/packages/config"
	"rajasms-account-monitor/packages/logger"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/leekchan/accounting"
	"github.com/xpartacvs/go-dishook"
)

type Webhook dishook.Payload

var (
	w    *Webhook
	once sync.Once
)

func GetInstance() *Webhook {
	once.Do(func() {
		w = &Webhook{
			Username:  config.Get().DishookBotName(),
			AvatarUrl: dishook.Url(config.Get().DishookBotAvatarURL()),
			Content:   config.Get().DishookBotMessage(),
		}
	})
	return w
}

func (p *Webhook) AddReminder(margin, balance, grace uint, expiry time.Time) *Webhook {
	ac := accounting.Accounting{
		Symbol:   "Rp",
		Thousand: ".",
		Decimal:  ",",
		Format:   "%s %v",
	}
	moneyBalance := ac.FormatMoney(balance)
	moneyMargin := ac.FormatMoney(margin)
	remainingDays := uint(time.Until(expiry).Hours() / 24)
	title := "Saldo Akun Minim"
	desc := "Saldo kurang dari " + moneyMargin + "Segera lakukan topup atau SMS tidak bisa terkirim."

	if remainingDays <= grace {
		title = "Mendekati Tanggal Kedaluarsa"
		desc = "Masa aktif saldo akun hampir berakhir. Segera lakukan topup, atau saldo hangus."
	}

	embed := dishook.Embed{
		Color:       dishook.ColorWarn,
		Url:         "https://raja-sms.com/topupsaldo/",
		Title:       title,
		Description: desc,
		Fields: []dishook.Field{
			{
				Name:   "Saldo Sekarang",
				Value:  moneyBalance,
				Inline: true,
			},
			{
				Name:   "Tanggal Kedaluarsa",
				Value:  expiry.Format("_2 Jan 2006"),
				Inline: true,
			},
			{
				Name:   "Saldo Hangus Dalam",
				Value:  strconv.FormatUint(uint64(remainingDays), 10) + " hari",
				Inline: true,
			},
		},
	}

	p.Embeds = nil
	p.Embeds = append(p.Embeds, embed)

	return p
}

func (p *Webhook) Send(url string) error {
	if p.Embeds == nil {
		logger.Log().Warn().Msg("Webhook has nothing to send")
		return nil
	}

	rgxUrl := regexp.MustCompile("^https?://discord.com/api/webhooks/.*")
	if !rgxUrl.MatchString(url) {
		p.Embeds = nil
		return errors.New("invalid webhook url")
	}

	_, err := dishook.Send(url, dishook.Payload(*p))
	p.Embeds = nil
	if err != nil {
		return err
	}

	return nil
}
