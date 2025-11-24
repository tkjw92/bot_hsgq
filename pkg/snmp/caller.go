package snmp

import (
	"log"

	"github.com/gosnmp/gosnmp"
)

func NewSNMPClient(target, community string) *SNMPClient {
	conn := &gosnmp.GoSNMP{
		Target:    target,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   gosnmp.Default.Timeout,
		Retries:   gosnmp.Default.Retries,
	}

	if err := conn.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return &SNMPClient{conn: conn}
}

func GetByName(c *SNMPClient, name string) Onu {
	indexes := searchByName(c, name)
	return getOnu(c, indexes)
}

func GetByMac(c *SNMPClient, name string) Onu {
	indexes := searchByMac(c, name)
	return getOnu(c, indexes)
}

func GetOnuList(c *SNMPClient) Onu {
	indexes := getIndex(c)

	return Onu{
		Name:   getOnuName(c, indexes),
		Status: getOnuStatus(c, indexes),
		Mac:    getOnuMac(c, indexes),
	}
}
