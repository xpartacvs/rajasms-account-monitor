FROM golang:1.16-alpine AS builder
WORKDIR /builder/src
COPY . .
RUN mkdir -p /builder/bin
RUN go build -ldflags="-s -w" -o /builder/bin/rajasms-monitor main.go

FROM alpine:latest
LABEL maintainer="xpartacvs@gmail.com"
ENV TZ=Asia/Jakarta
ENV DISCORD_BOT_MESSAGE=Reminder\ akun\ RajaSMS
ENV LOGMODE=disabled
ENV RAJASMS_LOWBALANCE=100000
ENV RAJASMS_GRACEPERIOD=7
ENV SCHEDULE=0\ 0\ *\ *\ *
WORKDIR /usr/local/bin
RUN apk update
RUN apk add --no-cache tzdata
COPY --from=builder /builder/bin/rajasms-monitor .
CMD ["rajasms-monitor"]
