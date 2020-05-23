package connection

type NetWork string
type Port int
type Command string

const (
	TCP = "tcp"
	UDP = "udp"

	VIDEO_PORT Port = 40921
	AUDIO_PORT Port = 40922
	CTRL_PORT  Port = 40923
	PUSH_PORT  Port = 40924
	EVENT_PORT Port = 40925
	IP_PORT    Port = 40926

	VIDEO_NETWORK NetWork = TCP
	AUDIO_NETWORK NetWork = TCP
	CTRL_NETWORK  NetWork = TCP
	PUSH_NETWORK  NetWork = UDP
	EVENT_NETWORK NetWork = TCP

	CommandSeparator         = ";"
	EnableVideo      Command = "stream on"
	DisableVideo     Command = "stream off"
	EnableAudio      Command = "audio on"
	DisableAudio     Command = "audio off"
)
