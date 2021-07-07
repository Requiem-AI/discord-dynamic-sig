package services

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cloakd/common/context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
)

type HttpService struct {
	DefaultService
	BaseURL string
	Port    int
	ssvc    *SignatureService
	dsvc    *DiscordService
}

var ErrUnauthorized = errors.New("unauthorized")

func (svc HttpService) Id() string {
	return "http"
}

func (svc *HttpService) Configure(ctx *context.Context) error {
	urlFlag := flag.String("url", "http://discord.cloakd.co.uk", "base url of service")
	portFlag := flag.Int("port", 9099, "port to serve http on")
	flag.Parse()

	svc.BaseURL = *urlFlag
	svc.Port = *portFlag
	logrus.WithField("base_url", svc.BaseURL).Info("BaseURL")

	return svc.DefaultService.Configure(ctx)
}

func (svc *HttpService) Start() error {
	r := gin.New()

	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	r.Static("static", "static")

	r.GET("/ping", svc.ping)
	r.GET("/server/:id/image.png", svc.discordServerImage)
	r.GET("/invite/:id/image.png", svc.discordInviteImage)

	svc.dsvc = svc.ctx.Service(DISCORD_SVC).(*DiscordService)
	svc.ssvc = svc.ctx.Service(SIGNATURE_SVC).(*SignatureService)

	return r.Run(fmt.Sprintf(":%v", svc.Port))

	//return autotls.Run(r, "ai.cloakd.co.uk")
}

func (svc *HttpService) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (svc *HttpService) discordServerImage(c *gin.Context) {
	serverId := c.Param("id")
	logrus.Printf("Looking up server: %s", serverId)
	details, err := svc.dsvc.InviteDetailFromServer(serverId)
	if err != nil {
		log.Printf("Server not found")
		_ = c.AbortWithError(404, err)
		return
	}

	img, err := svc.ssvc.Generate(details, svc.ssvc.DiscordDefaultSignature())
	if err != nil {
		log.Printf("Unable to generate media")
		_ = c.AbortWithError(400, err)
		return
	}

	c.Data(200, "image/gif", img)
}

func (svc *HttpService) discordInviteImage(c *gin.Context) {
	inviteId := c.Param("id")
	logrus.Printf("Looking up invite: %s", inviteId)
	details, err := svc.dsvc.InviteDetail(inviteId)
	if err != nil {
		log.Printf("Invite not found")
		_ = c.AbortWithError(404, err)
		return
	}

	img, err := svc.ssvc.Generate(details, svc.ssvc.DiscordDefaultSignature())
	if err != nil {
		log.Printf("Unable to generate media")
		_ = c.AbortWithError(400, err)
		return
	}

	c.Data(200, "image/gif", img)
}
