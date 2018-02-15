package notifications

type SlackConfig struct {
	Url         string `fig:"url"`
	ChannelName string `fig:"channel_name"`
	IconEmoji   string `fig:"icon_emoji"`
}
