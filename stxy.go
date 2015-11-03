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
		interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
		for {
			client, err := statsd.NewClient(c.String("s"), c.String("p"))
			// handle any errors
			if err != nil {
				log.Fatal(err)
			}
			// make sure to clean up
			defer client.Close()
			initial_stats := get_stats(c.String("haproxy-url"))
			previous := map[string]int64{}
			for _, v := range initial_stats {
				if v[1] == "BACKEND" {
					previous[fmt.Sprint("1xx_", v[0])] = get_value(v, "hrsp_1xx", 39)
					previous[fmt.Sprint("2xx_", v[0])] = get_value(v, "hrsp_2xx", 40)
					previous[fmt.Sprint("3xx_", v[0])] = get_value(v, "hrsp_3xx", 41)
					previous[fmt.Sprint("4xx_", v[0])] = get_value(v, "hrsp_4xx", 42)
					previous[fmt.Sprint("5xx_", v[0])] = get_value(v, "hrsp_5xx", 43)
				}
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
			records := get_stats(c.String("haproxy-url"))
			for _, record := range records {
				if record[1] == "BACKEND" {
					go send_gauge(client, record, "scur", 4)
					go send_gauge(client, record, "smax", 5)
					go send_gauge(client, record, "ereq", 12)
					go send_gauge(client, record, "econ", 13)
					go send_gauge(client, record, "rate", 33)
					go send_gauge(client, record, "bin", 8)
					go send_gauge(client, record, "bout", 9)
					go send_counter(previous[fmt.Sprint("1xx_", record[0])], client, record, "hrsp_1xx", 39)
					go send_counter(previous[fmt.Sprint("2xx_", record[0])], client, record, "hrsp_2xx", 40)
					go send_counter(previous[fmt.Sprint("3xx_", record[0])], client, record, "hrsp_3xx", 41)
					go send_counter(previous[fmt.Sprint("4xx_", record[0])], client, record, "hrsp_4xx", 42)
					go send_counter(previous[fmt.Sprint("5xx_", record[0])], client, record, "hrsp_5xx", 43)
					go send_gauge(client, record, "qtime", 58)
					go send_gauge(client, record, "ctime", 59)
					go send_gauge(client, record, "rtime", 60)
					go send_gauge(client, record, "ttime", 61)
				}
			}
			color.White("-------------------")
			time.Sleep(time.Duration(interval) * time.Millisecond)
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

func send_counter(previous int64, client statsd.Statter, v []string, name string, position int64) {
	stat := fmt.Sprint(v[0], ".", name)
	value_at_interval, _ := strconv.ParseInt(v[position], 10, 64)
	value := value_at_interval - previous
	fmt.Println(fmt.Sprint(stat, ":", value, "|c"))
	client.Inc(stat, value, 1)
}

func get_value(v []string, name string, position int64) int64 {
	value, _ := strconv.ParseInt(v[position], 10, 64)
	return value
}

func get_stats(url string) [][]string {
	stats, _ := sh.Command("curl", url).Output()
	r := csv.NewReader(strings.NewReader(string(stats)))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records
}
