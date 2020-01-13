package network_api

type API interface {
	Info() (*Info, error)
	Peers() ([]BgpPeer, error)
	Interfaces() ([]Interface, error)
}

type Info struct {
	Uptime float64

	FreeMemory  float64
	TotalMemory float64

	Model   string
	Serial  string
	Version string
}

type BgpPeer struct {
	Up bool

	AS     int64
	Remote string

	PrefixCount float64

	MessagesReceived float64
	MessagesSent     float64
}

type Interface struct {
	Name        string
	Description string

	RxPackets float64
	RxBits    float64
	RxDrops   float64
	RxErrors  float64

	TxPackets float64
	TxBits    float64
	TxDrops   float64
	TxErrors  float64
}
