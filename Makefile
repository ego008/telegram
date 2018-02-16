LANGUAGE_CODE = en
LANGUAGE_FORMAT = yaml

CONFIGS_TRANSLATIONS_FOLDER = ./configs/translations

build:
	go build

debug:
	go build -tags=debug

translation:
	goi18n merge \
	-format $(LANGUAGE_FORMAT) \
	-sourceLanguage $(LANGUAGE_CODE) \
	-outdir $(CONFIGS_TRANSLATIONS_FOLDER) \
	$(CONFIGS_TRANSLATIONS_FOLDER)src/*/*

localization:
	make translation
	goi18n \
	-format $(LANGUAGE_FORMAT) \
	-sourceLanguage $(LANGUAGE_CODE) \
	-outdir $(CONFIGS_TRANSLATIONS_FOLDER) \
	$(CONFIGS_TRANSLATIONS_FOLDER)*.all.$(LANGUAGE_FORMAT) \
	$(CONFIGS_TRANSLATIONS_FOLDER)*.untranslated.$(LANGUAGE_FORMAT)

image:
	docker-compose build