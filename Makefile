development:
	go build

webhook-release:
	go build -tags "webhook"

production:
	go build -tags "easterEggs"

webhook-production:
	go build -tags "easterEggs webhook"