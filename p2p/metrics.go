package p2p

import "github.com/iost-official/Go-IOS-Protocol/metrics"

var (
	neighborCountGauge = metrics.NewGauge("iost_p2p_neighbor_count", nil)
	routingCountGauge  = metrics.NewGauge("iost_p2p_routing_count", nil)
	byteOutSummary     = metrics.NewSummary("iost_p2p_byte_out", []string{"mtype"})
	packetOutCounter   = metrics.NewCounter("iost_p2p_packet_out", []string{"mtype"})
	byteInSummary      = metrics.NewSummary("iost_p2p_byte_in", []string{"mtype"})
	packetInCounter    = metrics.NewCounter("iost_p2p_packet_in", []string{"mtype"})
)
