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

var http_responses []string
var initial_response_values map[string]int64

func main() {
	app := cli.NewApp()
	app.Name = "stxy"
	app.Version = "0.0.3"
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
			get_initial_values(records)
			interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
			time.Sleep(time.Duration(interval) * time.Millisecond)
			for _, v := range records {
				if v[1] == "BACKEND" {
					send_gauge(client, v, "scur", 4)
					send_gauge(client, v, "smax", 5)
					send_gauge(client, v, "ereq", 12)
					send_gauge(client, v, "econ", 13)
					send_gauge(client, v, "rate", 33)
					send_gauge(client, v, "bin", 8)
					send_gauge(client, v, "bout", 9)
					send_counter(client, v, "hrsp_1xx", 39)
					send_counter(client, v, "hrsp_2xx", 40)
					send_counter(client, v, "hrsp_3xx", 41)
					send_counter(client, v, "hrsp_4xx", 42)
					send_counter(client, v, "hrsp_5xx", 43)
					send_gauge(client, v, "qtime", 58)
					send_gauge(client, v, "ctime", 59)
					send_gauge(client, v, "rtime", 60)
					send_gauge(client, v, "ttime", 61)
				}
			}
			color.White("-------------------")
		}
	}
	app.Run(os.Args)
}

func send_gauge(client statsd.Statter, v []string, name string, position int64) {
	stat := fmt.Sprint(v[0], ".", name)
	value, _ := strconv.ParseInt(v[position], 10, 64)
	fmt.Println(fmt.Sprint(stat, ":", value, "|g"))
	client.Gauge(stat, value, 1.0)
}

func send_counter(client statsd.Statter, v []string, name string, position int64) {
	stat := fmt.Sprint(v[0], ".", name)
	value_at_interval, _ := strconv.ParseInt(v[position], 10, 64)
	value := value_at_interval - initial_response_values[name]
	fmt.Println(fmt.Sprint(stat, ":", value, "|c"))
	client.Inc(stat, value, 1)
}

func get_value(v []string, name string, position int64) int64 {
	value, _ := strconv.ParseInt(v[position], 10, 64)
	return value
}

func get_initial_values(records [][]string) {
	http_responses := []string{"hsrp_1xx", "hsrp_2xx", "hsrp_3xx", "hsrp_4xx", "hsrp_5xx"}
	initial_response_values := map[string]int64{
		"hsrp_1xx": 0,
		"hsrp_2xx": 0,
		"hsrp_3xx": 0,
		"hsrp_4xx": 0,
		"hsrp_5xx": 0,
	}
	for _, v := range records {
		var i int64 = 39
		for _, resp := range http_responses {
			initial_response_values[resp] = get_value(v, resp, i)
			i += 1
		}
	}
}
