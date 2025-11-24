package snmp

import (
	"github.com/gosnmp/gosnmp"
)

type SNMPClient struct {
	conn *gosnmp.GoSNMP
}

type Onu struct {
	Status   []string
	Name     []string
	Mac      []string
	Distance []string
	Rx       []string
	Tx       []string
	Vendor   []string
}

const (
	RxDir = "4"
	TxDir = "5"
)
