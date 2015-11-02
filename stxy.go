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
	app.Version = "0.0.3"
	app.Usage = "haproxy stats to statsd"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "haproxy-url",
			Value: "localhost:22002/;csv",
			Usage: "host:port of haproxy server",
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
			Usage: "time in milliseconds",
			Value: "10000",
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
			if records == nil {
				log.Fatal("Unable to read stats from HaProxy")
			}
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range records {
				go send_stat(client, v, "scur", 4)
				go send_stat(client, v, "smax", 5)
				go send_stat(client, v, "ereq", 12)
				go send_stat(client, v, "econ", 13)
				go send_stat(client, v, "rate", 33)
				go send_stat(client, v, "bin", 8)
				go send_stat(client, v, "bout", 9)
				go send_stat(client, v, "hrsp_1xx", 39)
				go send_stat(client, v, "hrsp_2xx", 40)
				go send_stat(client, v, "hrsp_3xx", 41)
				go send_stat(client, v, "hrsp_4xx", 42)
				go send_stat(client, v, "hrsp_5xx", 43)
				go send_stat(client, v, "qtime", 58)
				go send_stat(client, v, "ctime", 59)
				go send_stat(client, v, "rtime", 60)
				go send_stat(client, v, "ttime", 61)
			}
			color.Cyan("-------------------")
			interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
	app.Run(os.Args)
}

func send_stat(client statsd.Statter, v []string, name string, position int64) {
	if v[1] == "BACKEND" {
		stat := fmt.Sprint(v[0], ".", name)
		value, _ := strconv.ParseInt(v[position], 10, 64)
		fmt.Println(fmt.Sprint(stat, ":", value, "|g"))
		client.Gauge(stat, value, 1.0)
	}
}
