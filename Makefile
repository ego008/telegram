development:
	goi18n merge -sourceLanguage en-us -outdir ./i18n/ ./i18n/source/*
	goi18n -sourceLanguage en-us -outdir ./i18n/ ./i18n/*.all.json ./i18n/*.untranslated.json
	go build

production:
	goi18n merge -sourceLanguage en-us -outdir ./i18n/ ./i18n/source/*
	goi18n -sourceLanguage en-us -outdir ./i18n/ ./i18n/*.all.json ./i18n/*.untranslated.json
	go build -tags "eggs"