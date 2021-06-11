package webhook

import (
	"errors"
	"fmt"
	"rajasms-account-monitor/packages/config"
	"rajasms-account-monitor/packages/logger"
	"regexp"
	"sync"
	"time"

	"github.com/leekchan/accounting"
	"github.com/xpartacvs/go-dishook"
)

type webhook dishook.Payload

type Webhook interface {
	AddReminder(margin, balance, grace uint, expiry time.Time) Webhook
	Send(url string) error
}

var (
	w    Webhook
	once sync.Once
)

func GetInstance() Webhook {
	once.Do(func() {
		w = &webhook{
			Username:  config.Get().DishookBotName(),
			AvatarUrl: dishook.Url(config.Get().DishookBotAvatarURL()),
			Content:   config.Get().DishookBotMessage(),
		}
	})
	return w
}

func (p *webhook) AddReminder(margin, balance, grace uint, expiry time.Time) Webhook {
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
	desc := fmt.Sprintf("Saldo kurang dari %s Segera lakukan topup atau SMS tidak bisa terkirim.", moneyMargin)

	if remainingDays <= grace {
		title = "Mendekati Tanggal Kedaluarsa"
		desc = "Masa aktif saldo akun hampir berakhir. Segera lakukan topup, atau saldo hangus."
	}

	embed := dishook.Embed{
		Color:       14327864,
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
				Value:  fmt.Sprintf("%d hari", remainingDays),
				Inline: true,
			},
		},
	}

	p.Embeds = nil
	p.Embeds = append(p.Embeds, embed)

	return p
}

func (p *webhook) Send(url string) error {
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
