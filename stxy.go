package main

import (
	"encoding/csv"
	//"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/olekukonko/tablewriter"
	//	"github.com/cactus/go-statsd-client/statsd"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// init
	app := cli.NewApp()
	app.Name = "stxy"
	app.Version = "0.0.1"
	app.Usage = "haproxy stats to statsd"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "haproxy-url",
			Value: "localhost:22002/;csv",
			Usage: "host:port of redis servier",
		},
		cli.StringFlag{
			Name:  "statsd-url, s",
			Value: "localhost:8125",
			Usage: "host:port of statsd server",
		},
		cli.StringFlag{
			Name:  "interval,i",
			Usage: "time in milliseconds to periodically check redis",
			Value: "5000",
		},
	}
	app.Action = func(c *cli.Context) {
		for {
			stats, _ := sh.Command("curl", c.String("haproxy-url")).Output()
			r := csv.NewReader(strings.NewReader(string(stats)))
			records, err := r.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(records)
			table := tablewriter.NewWriter(os.Stdout)
			for _, value := range records {
				//fmt.Println(value)
				table.Append(value)
			}
			//sendStats(c.String("statsd-host"), c.String("prefix"), gauges, counters)
			table.Render()
			color.White("-------------------")
			interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
	app.Run(os.Args)
}
