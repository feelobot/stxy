package main

import (
	"encoding/csv"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
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
			Name:  "prefix,p",
			Usage: "statsd namespace",
			Value: "haproxy",
		},
		cli.StringFlag{
			Name:  "interval,i",
			Usage: "time in milliseconds to periodically check redis",
			Value: "5000",
		},
	}
	app.Action = func(c *cli.Context) {
		for {
			client, err := statsd.NewClient(c.String("s"), c.String("p"))
			// handle any errors
			if err != nil {
				log.Fatal(err)
			}
			// make sure to clean up
			defer client.Close()
			stats, _ := sh.Command("curl", c.String("haproxy-url")).Output()
			r := csv.NewReader(strings.NewReader(string(stats)))
			records, err := r.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range records {
				send_stat(client, v, "scur", 4)
				send_stat(client, v, "smax", 5)
				send_stat(client, v, "ereq", 12)
				send_stat(client, v, "econ", 13)
				send_stat(client, v, "rate", 33)
				send_stat(client, v, "bin", 8)
				send_stat(client, v, "bout", 9)
				send_stat(client, v, "hrsp_1xx", 39)
				send_stat(client, v, "hrsp_2xx", 40)
				send_stat(client, v, "hrsp_3xx", 41)
				send_stat(client, v, "hrsp_4xx", 42)
				send_stat(client, v, "hrsp_5xx", 43)
				send_stat(client, v, "qtime", 58)
				send_stat(client, v, "ctime", 59)
				send_stat(client, v, "rtime", 60)
				send_stat(client, v, "ttime", 61)
			}
			color.White("-------------------")
			interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
	app.Run(os.Args)
}

func send_stat(client statsd.Statter, v []string, name string, position int64) {
	stat := fmt.Sprint(v[0], ".", v[1], ".", name)
	value, _ := strconv.ParseInt(v[position], 10, 64)
	fmt.Println(fmt.Sprint(stat, ":", value, "|g"))
	client.Gauge(stat, value, 1.0)
}
