package services

import (
	"encoding/json"
	"fmt"
	"github.com/requiem-ai/discord-dynamic-sig/server"
	"io/ioutil"
	"net/http"
	"time"
)

type DiscordService struct {
	DefaultService
	client *http.Client
}

const DISCORD_SVC = "discord_svc"

func (svc DiscordService) Id() string {
	return DISCORD_SVC
}

func (svc *DiscordService) Start() error {

	svc.client = &http.Client{
		Timeout: 1 * time.Second,
	}

	return nil
}

func (svc *DiscordService) ServerDetail(serverId string) (*server.ServerInfo, error) {
	resp, err := svc.client.Get(fmt.Sprintf("https://discord.com/api/guilds/%s/widget.json", serverId))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var det server.ServerInfo
	err = json.Unmarshal(body, &det)
	if err != nil {
		return nil, err
	}

	return &det, nil
}

func (svc *DiscordService) InviteDetail(inviteId string) (*server.DiscordInvite, error) {
	resp, err := svc.client.Get(fmt.Sprintf("https://discord.com/api/v9/invites/%s?with_counts=true&with_expiration=true", inviteId))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var inv server.DiscordInvite
	err = json.Unmarshal(body, &inv)
	if err != nil {
		return nil, err
	}

	return &inv, nil
}

func (svc *DiscordService) InviteDetailFromServer(serverId string) (*server.DiscordInvite, error) {
	detail, err := svc.ServerDetail(serverId)
	if err != nil {
		return nil, err
	}

	return svc.InviteDetail(detail.InviteCode())
}
