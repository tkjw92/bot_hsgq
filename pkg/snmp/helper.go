package snmp

import "strings"

func findIndexesContaining(stores []string, input string) []int {
	var result []int

	for i, v := range stores {
		if strings.Contains(strings.ToLower(v), strings.ToLower(input)) {
			result = append(result, i)
		}
	}

	return result
}

func searchByName(c *SNMPClient, name string) (index []string) {
	indexes := getIndex(c)

	onu := Onu{
		Name: getOnuName(c, indexes),
	}

	idx := findIndexesContaining(onu.Name, name)
	for _, i := range idx {
		index = append(index, indexes[i])
	}

	return
}

func searchByMac(c *SNMPClient, mac string) (index []string) {
	indexes := getIndex(c)

	onu := Onu{
		Mac: getOnuMac(c, indexes),
	}

	idx := findIndexesContaining(onu.Mac, mac)
	for _, i := range idx {
		index = append(index, indexes[i])
	}

	return
}
