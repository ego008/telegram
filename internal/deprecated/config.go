package main

import (
	"flag"
	"log"
	"regexp"

	botan "github.com/botanio/sdk/go"
	config "github.com/olebedev/config"
	patreon "github.com/toby3d/go-patreon"
)

const (
	version = `3.0 "Cure Aqua"`
	build   = 187

	bb = "http://beta.hentaidb.pw"
)

var (
	adm, chID, pCampaign                  int
	botToken, webSet, webListen, webServe string
	cfg                                   *config.Config
	err                                   error
	glbr                                  map[string]interface{}
	p                                     *patreon.Client
	r                                     *regexp.Regexp

	flagLocale  = flag.String("i18n", "./i18n", "load locales from custom path")
	flagConfig  = flag.String("config", "./config.yml", "load custom config file")
	flagDB      = flag.String("db", "./hentai.db", "select custom db file")
	flagDebug   = flag.Bool("debug", false, "enable debug logs")
	flagWebhook = flag.Bool("webhook", false, "enable webhooks support")
	flagCache   = flag.Int("cache", 0, "cache time in seconds for inline-search results")
)

func init() {
	flag.Parse()
	go langInit()
	go dbInit()

	r, err = regexp.Compile(`([\w]+:+[[:graph:]]+)\s?`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Load configuration
	cfg, err = config.ParseYamlFile(*flagConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Botan
	botanToken, err := cfg.String("botan.token")
	if err != nil {
		log.Fatalln(err.Error())
	}
	metrika = botan.New(botanToken)

	botToken, err = cfg.String("telegram.bot.token")
	if err != nil {
		log.Fatalln(err.Error())
	}

	webSet, err = cfg.String("telegram.webhook.set")
	if err != nil {
		log.Fatalln(err.Error())
	}

	webListen, err = cfg.String("telegram.webhook.listen")
	if err != nil {
		log.Fatalln(err.Error())
	}

	webServe, err = cfg.String("telegram.webhook.serve")
	if err != nil {
		log.Fatalln(err.Error())
	}

	chID, err = cfg.Int("telegram.channel.id")
	if err != nil {
		log.Fatalln(err.Error())
	}

	pID, err := cfg.String("patreon.client.id")
	if err != nil {
		log.Fatalln(err.Error())
	}

	pSecret, err := cfg.String("patreon.client.secret")
	if err != nil {
		log.Fatalln(err.Error())
	}

	pURI, err := cfg.String("patreon.client.redirect_uri")
	if err != nil {
		log.Fatalln(err.Error())
	}

	pCampaign, err = cfg.Int("patreon.campaign")
	if err != nil {
		log.Fatalln(err.Error())
	}

	p = patreon.NewClient(pID, pSecret, pURI)

	adm, err = cfg.Int("superadmin")
	if err != nil {
		log.Fatalln(err.Error())
	}

	glbr, err = cfg.Map("resource.0")
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("%s configuration is loaded", *flagConfig)
	log.Printf("Trying to run version %s (build %d).", version, build)
}
