.PHONY: all clean-all build cleand-deps deps ver

DATE := $(shell git log -1 --format="%cd" --date=short | sed s/-//g)
COUNT := $(shell git rev-list --count HEAD)
COMMIT := $(shell git rev-parse --short HEAD)

CHAT := $(shell echo $${CHAT:-mattermost})

CONFIGUREDIR := /etc/chattix
SERVICENAME := chattixd
WEBHOOKNAME := zabbix-to-${CHAT}

SERVICECONFIG := ${SERVICENAME}.conf
WEBHOOKCONFIG := ${WEBHOOKNAME}.conf

VERSION := "${DATE}.${COUNT}_${COMMIT}"

LDFLAGS := "-X main.version=${VERSION} -X main.definedChat=${CHAT}"

SERVICESTATUS := $(shell systemctl status chattixd)


default: all

all: clean-all deps build

ver:
	@echo ${VERSION}

clean-all: clean-deps
	@echo Clean builded binaries
	rm -rf .out/
	@echo Done

clean-deps:
	@echo Clean dependencies
	rm -rf vendor/*

deps:
	dep ensure

build: build-service build-hook


build-service:
	@echo Build ${SERVICENAME}
	go build -o .out/${SERVICENAME} -ldflags ${LDFLAGS} action_ack/*.go

build-hook:
	@echo Build ${WEBHOOKNAME}
	go build -o .out/${WEBHOOKNAME} -ldflags ${LDFLAGS} webhook/*.go

install: $(CONFIGUREDIR)
	@echo Installing service 
	cp .out/${SERVICENAME} /usr/bin/${SERVICENAME}
	@echo Installing webhook 
	cp .out/${WEBHOOKNAME} /usr/bin/${WEBHOOKNAME}
	@echo Install example service config
	cp action_ack/config.conf ${CONFIGUREDIR}/${SERVICECONFIG}
	cp webhook/config.conf ${CONFIGUREDIR}/${WEBHOOKCONFIG}
	@echo Install systemd service 
	cp systemd/chattixd.service /lib/systemd/system/ 
	systemctl daemon-reload
	@echo Done

/etc/chattix:
	test ! -d $(CONFIGUREDIR) && mkdir $(CONFIGUREDIR)

uninstall:
	@echo Stopping service
	@echo ${status}
	systemctl stop chattixd
	@echo Remove service
	rm /lib/systemd/system/chattixd.service
	systemctl daemon-reload
	@echo Remove binaries
	rm -rf /usr/bin/${SERVICENAME} /usr/bin/${WEBHOOKNAME}
	@echo Remove configuration
	rm -rf ${CONFIGUREDIR}
