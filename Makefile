all:
	make translation
	make localization
	go build \
	-ldflags \
	"-X main.verHash=`git rev-parse --short HEAD` \
	-X main.verTimeStamp=`date -u +%Y-%m-%d.%H:%M:%S`"

build:
	go build

debug:
	go build -tags=debug

translation:
	goi18n merge \
	-format yaml \
	-sourceLanguage en \
	-outdir ./i18n/ ./i18n/src/*/*

localization:
	make translation
	goi18n \
	-format yaml \
	-sourceLanguage en \
	-outdir ./i18n/ \
	./i18n/*.all.yaml ./i18n/*.untranslated.yaml

image:
	docker-compose build