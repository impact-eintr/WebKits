package loadgenerator

// Gotickets 表示Gouroutine票池的接口
type GoTickets interface {
	Fetch()
	Return()
	Active() bool // 票池是否已经被激活
	Total() uint32
	Remainder() uint32
}
