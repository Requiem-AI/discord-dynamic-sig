package server

import (
	"fmt"
	"time"
)

type DiscordInvite struct {
	Code  string `json:"code"`
	Guild struct {
		ID                string      `json:"id"`
		Name              string      `json:"name"`
		Splash            string      `json:"splash"`
		Banner            interface{} `json:"banner"`
		Description       interface{} `json:"description"`
		Icon              string      `json:"icon"`
		Features          []string    `json:"features"`
		VerificationLevel int         `json:"verification_level"`
		VanityURLCode     interface{} `json:"vanity_url_code"`
		Nsfw              bool        `json:"nsfw"`
		NsfwLevel         int         `json:"nsfw_level"`
	} `json:"guild"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"channel"`
	ApproximateMemberCount   int       `json:"approximate_member_count"`
	ApproximatePresenceCount int       `json:"approximate_presence_count"`
	ExpiresAt                time.Time `json:"expires_at"`
	ExpiredAt                time.Time `json:"expired_at"`
}

func (d *DiscordInvite) ServerAvatar(size int) string {
	return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png?size=%v", d.Guild.ID, d.Guild.Icon, size)
}
