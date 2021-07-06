package server

import "strings"

type ServerInfo struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	InstantInvite string        `json:"instant_invite"`
	Channels      []interface{} `json:"channels"`
	Members       []struct {
		ID            string      `json:"id"`
		Username      string      `json:"username"`
		Discriminator string      `json:"discriminator"`
		Avatar        interface{} `json:"avatar"`
		Status        string      `json:"status"`
		AvatarURL     string      `json:"avatar_url"`
		Game          struct {
			Name string `json:"name"`
		} `json:"game,omitempty"`
	} `json:"members"`
	PresenceCount int `json:"presence_count"`
}

func (s *ServerInfo) InviteCode() string {
	return strings.ReplaceAll(s.InstantInvite, "https://discord.com/invite/", "")
}
