package snmp

import (
	"fmt"
	"log"
	"strings"

	"github.com/gosnmp/gosnmp"
)

func getAllOnu(c *SNMPClient) Onu {
	indexes := getIndex(c)

	return Onu{
		Status:   getOnuStatus(c, indexes),
		Name:     getOnuName(c, indexes),
		Mac:      getOnuMac(c, indexes),
		Distance: getOnuDistance(c, indexes),
		Rx:       getOnuOptical(c, RxDir, indexes),
		Tx:       getOnuOptical(c, TxDir, indexes),
	}
}

func getOnu(c *SNMPClient, indexes []string) Onu {
	return Onu{
		Status:   getOnuStatus(c, indexes),
		Name:     getOnuName(c, indexes),
		Mac:      getOnuMac(c, indexes),
		Distance: getOnuDistance(c, indexes),
		Rx:       getOnuOptical(c, RxDir, indexes),
		Tx:       getOnuOptical(c, TxDir, indexes),
	}
}

func getIndex(c *SNMPClient) (onu []string) {
	result, err := c.conn.BulkWalkAll("1.3.6.1.4.1.50224.3.3.2.1.1")
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range result {
		if res.Type == gosnmp.Gauge32 {
			val := res.Value.(uint)
			onu = append(onu, fmt.Sprintf("%d", val))
		}
	}

	return
}

func getOnuName(c *SNMPClient, indexes []string) (onu []string) {
	oids := []string{}

	for _, idx := range indexes {
		oids = append(oids, "1.3.6.1.4.1.50224.3.3.2.1.2."+idx)
	}

	result, err := c.conn.Get(oids)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range result.Variables {
		if res.Type == gosnmp.NoSuchInstance {
			onu = append(onu, "-")
			continue
		}

		if res.Type == gosnmp.OctetString {
			val := res.Value.([]byte)
			onu = append(onu, string(val))
		}
	}

	return
}

func getOnuMac(c *SNMPClient, indexes []string) (onu []string) {
	oids := []string{}

	for _, idx := range indexes {
		oids = append(oids, "1.3.6.1.4.1.50224.3.3.2.1.7."+idx)
	}

	result, err := c.conn.Get(oids)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range result.Variables {
		if res.Type == gosnmp.NoSuchInstance {
			onu = append(onu, "-")
			continue
		}

		if res.Type == gosnmp.OctetString {
			val := res.Value.([]byte)
			var mac []string
			for _, v := range val {
				mac = append(mac, fmt.Sprintf("%02x", v))
			}
			onu = append(onu, strings.Join(mac, ":"))
		}
	}

	return
}

func getOnuStatus(c *SNMPClient, indexes []string) (onu []string) {
	oids := []string{}

	for _, idx := range indexes {
		oids = append(oids, "1.3.6.1.4.1.50224.3.3.2.1.9."+idx)
	}

	result, err := c.conn.Get(oids)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range result.Variables {
		if res.Type == gosnmp.NoSuchInstance {
			onu = append(onu, "-")
			continue
		}

		if res.Type == gosnmp.Integer {
			val := res.Value.(int)
			if val == 1 {
				onu = append(onu, "up")
			} else {
				onu = append(onu, "down")
			}
		}
	}

	return
}

func getOnuOptical(c *SNMPClient, power_dir string, indexes []string) (onu []string) {
	oids := []string{}

	base_oid := fmt.Sprintf("1.3.6.1.4.1.50224.3.3.3.1.%s.", power_dir)
	for _, idx := range indexes {
		oids = append(oids, base_oid+idx+".0.0")
	}

	result, err := c.conn.Get(oids)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range result.Variables {
		if res.Type == gosnmp.NoSuchInstance {
			onu = append(onu, "-")
			continue
		}

		if res.Type == gosnmp.Integer {
			val_int := res.Value.(int)
			val_float := float64(val_int) * 0.01
			val := fmt.Sprintf("%.2f", val_float)

			onu = append(onu, val)
		}
	}

	return
}

func getOnuDistance(c *SNMPClient, indexes []string) (onu []string) {
	oids := []string{}

	for _, idx := range indexes {
		oids = append(oids, "1.3.6.1.4.1.50224.3.3.2.1.15."+idx)
	}

	result, err := c.conn.Get(oids)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range result.Variables {
		if res.Type == gosnmp.NoSuchInstance {
			onu = append(onu, "-")
			continue
		}

		if res.Type == gosnmp.Integer {
			val := res.Value.(int)
			onu = append(onu, fmt.Sprintf("%d", val))
		}
	}

	return
}
