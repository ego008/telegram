build:
	go build

production:
	make translation
	goi18n -sourceLanguage en-us -outdir ./i18n/ ./i18n/*.all.json ./i18n/*.untranslated.json
	go build -tags="eggs"

translation:
	goi18n merge -sourceLanguage en-us -outdir ./i18n/ ./i18n/source/*

development:
	make translation
	goi18n -sourceLanguage en-us -outdir ./i18n/ ./i18n/*.all.json ./i18n/*.untranslated.json
	go build -tags="debug"