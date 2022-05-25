package cmd

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func RenderTable(data []rttHost) {
	var newData [][]string
	for _, item := range data {
		host, rtt := item.hostName, item.avgRtt
		newData = append(newData, []string{host, rtt.String()})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Host", "Round-Trip Time"})

	for _, v := range newData {
		table.Append(v)
	}
	table.Render() // Send output
}
