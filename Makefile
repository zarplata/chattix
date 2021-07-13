DATE := $(shell git log -1 --format="%cd" --date=short | sed s/-//g)
COUNT := $(shell git rev-list --count HEAD)
COMMIT := $(shell git rev-parse --short HEAD)

MESSENGER := $(shell echo $${MESSENGER:-mattermost})

CONFIGUREDIR := /etc/chattix
INSTALLPREFIX := /usr/local

SERVICENAME := chattixd
WEBHOOKNAME := zabbix-to-${MESSENGER}

SERVICECONFIG := ${SERVICENAME}.conf
WEBHOOKCONFIG := ${WEBHOOKNAME}.conf

VERSION := "${DATE}.${COUNT}_${COMMIT}"

LDFLAGS := "-X main.version=${VERSION} -X main.definedMessenger=${MESSENGER}"

SERVICESTATUS := $(shell systemctl status chattixd)


default: all

.PHONY: all
all: clean-all build

.PHONY: service
service: clean-all deps build-service

.PHONY: ver
ver:
	@echo ${VERSION}

.PHONY: clean-all
clean-all:
	@echo Clean builded binaries
	rm -rf .out/
	@echo Done

.PHONY: build
build: build-service build-hook


.PHONY: build-service
build-service:
	@echo Build ${SERVICENAME}
	CGO_ENABLED=0 go build -o .out/${SERVICENAME} -ldflags ${LDFLAGS} action_ack/*.go

.PHONY: build-hook
build-hook:
	@echo Build ${WEBHOOKNAME}
	CGO_ENABLED=0 go build -o .out/${WEBHOOKNAME} -ldflags ${LDFLAGS} webhook/*.go

.PHONY: build-docker
build-docker:
	@echo Build dockerimage chattix
	docker build . -t chattix

.PHONY: install
install: $(CONFIGUREDIR)
	@echo Installing service 
	cp .out/${SERVICENAME} ${INSTALLPREFIX}/bin/${SERVICENAME}
	@echo Installing webhook 
	cp .out/${WEBHOOKNAME} ${INSTALLPREFIX}/bin/${WEBHOOKNAME}
	@echo Install example service config
	cp action_ack/config.conf ${CONFIGUREDIR}/${SERVICECONFIG}
	cp webhook/config.conf ${CONFIGUREDIR}/${WEBHOOKCONFIG}
	@echo Install systemd service 
	cp systemd/chattixd.service /lib/systemd/system/ 
	systemctl daemon-reload
	@echo Done

/etc/chattix:
	test ! -d $(CONFIGUREDIR) && mkdir $(CONFIGUREDIR)

.PHONY: uninstall
uninstall:
	@echo Stopping service
	@echo ${status}
	systemctl stop chattixd
	@echo Remove service
	rm /lib/systemd/system/chattixd.service
	systemctl daemon-reload
	@echo Remove binaries
	rm -rf /usr/bin/${SERVICENAME} ${INSTALLPREFIX}/bin/${WEBHOOKNAME}
	@echo Remove configuration
	rm -rf ${CONFIGUREDIR}
