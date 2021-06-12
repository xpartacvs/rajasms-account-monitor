# RajaSMS Account Monitor

Monitor status saldo dan tanggal kedaluarsa akun RajaSMS

## Cara Pakai

Yang paling enak ya pakai docker. Berikut ini step-stepnya:

```bash
# Clone repository dulu dan masuk ke directorynya
git clone git@github.com:xpartacvs/rajasms-account-monitor.git
cd rajasms-account-monitor

# Build image docker-nya
docker image build -t rajasms-account-monitor:latest .

# Run container-nya
docker container run \
    -it \
    -e DISCORD_WEBHOOKURL=... \
    -e RAJASMS_API_URL=... \
    -e RAJASMS_API_KEY=... \
    rajasms-account-monitor:latest
```

> **PENTING**: Minimal ada 3 (tiga) environment variable yang harus di-_assign_ yaitu `DISCORD_WEBHOOKURL`, `RAJASMS_API_URL`, dan `RAJASMS_API_KEY`. Lebih lengkapnya lihat bagian [konfigurasi](#konfigurasi)

## Konfigurasi

Konfigurasi aplikasi ini dapat dilakukan dengan menggunakan environment variables.

| **Variable**            | **Type**  | **Req** | **Default**             | **Description**                                                                                                                 |
| :---                    | :---      | :---:   | :---                    | :---                                                                                                                            |
| `DISCORD_WEBHOOKURL`    | `string`  | √       |                         | URL webhook Discord.                                                                                                            |
| `DISCORD_BOT_NAME`      | `string`  |         | suka-suka discord       | Nama bot yang akan muncul di channel Discord.                                                                                   |
| `DISCORD_BOT_AVATARURL` | `string`  |         | suka-suka discord       | URL ke file gambar yang akan digunakan sebagai avatar bot discord.                                                              |
| `DISCORD_BOT_MESSAGE`   | `string`  |         | `Reminder akun RajaSMS` | Pesan yang akan ditulis bot discord perihal status akun RajaSMS.                                                                |
| `LOGMODE`               | `string`  |         | `disabled`              | Mode log aplikasi: `debug`, `info`, `warn`, `error`, dan `disabled`.                                                            |
| `RAJASMS_API_URL`       | `string`  | √       |                         | URL server akun RajaSMS.                                                                                                        |
| `RAJASMS_API_KEY`       | `string`  | √       |                         | API key akun RajaSMS.                                                                                                           |
| `RAJASMS_LOWBALANCE`    | `integer` |         | `100000`                | Jika saldo <= nilai variabel ini maka alert via discord webhook akan terpicu.                                                   |
| `RAJASMS_GRACEPERIOD`   | `integer` |         | `7`                     | Jumlah hari menjelang tanggal kedaluarsa akun. Alert akan terpicu jika tanggal sekarang >= (tanggal kedaluarsa - variabel ini). |
| `SCHEDULE`              | `string`  |         | `0 0 * * *`             | Jadwal pemeriksaan status akun RajaSMS (dalam format CRON).                                                                     |

## Lisensi

MIT kok. Insya Allah _Open Source_ selamanya.

## Cara Kontribusi

- Silahkan PR saja.
- **WAJIB** menggunakan bahasa Indonesia jika ingin mengubah atau menambahkan info di `README.md`.
- Jika ingin ditiru, mohon pertimbangkan ini: _ATM lebih baik dari ATP_  (**ATM**=Amati Tiru Modifikasi, **ATP**=Amati Tiru Plek-plek)
- Oh iya aplikasi ini dibuat dengan bahasa pemrograman _GO_ ya.

## Minta tolong dibantu

- [ISSUE #1](https://github.com/xpartacvs/rajasms-account-monitor/issues/1) Ada **BUG** saat call API ke server RajaSMS. API di-_call_ lebih dari satu kali sehingga kena `rate-limit` _(meskipun tidak bikin container exit)_.
- Saat ini alertnya hanya ke discord. Monggo yang mau bantu tambah fitur slack, sms, email, dll `you are very welcome`.
