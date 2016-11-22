development:
	goi18n merge -sourceLanguage en-us -outdir ./i18n/ ./i18n/source/*
	go build

production:
	goi18n merge -sourceLanguage en-us -outdir ./i18n/ ./i18n/source/*
	go build -tags "eggs"