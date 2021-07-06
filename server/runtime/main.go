package main

import (
	"github.com/cloakd/common/context"
	"github.com/joho/godotenv"
	"github.com/requiem-ai/discord-dynamic-sig/server/services"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	ctx, err := context.NewContext(
		//&services.Datastore{},
		&services.DiscordService{},
		&services.SignatureService{},
		&services.HttpService{},
	)

	if err != nil {
		logrus.WithError(err).Fatal("FATAL")
		return
	}

	err = ctx.Run()
	logrus.WithError(err).Error("FATAL")
}
