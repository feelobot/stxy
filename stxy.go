package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	//	"github.com/cactus/go-statsd-client/statsd"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"log"
	"net/http"
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
			Name:  "haproxy-url, h",
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
			resp, err := http.Get(c.String("h"))
			if err != nil {
				// handle error
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			n := bytes.IndexByte(body, 0)
			s := string(body[:n])
			r := csv.NewReader(strings.NewReader(s))
			records, err := r.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(records)
			//sendStats(c.String("statsd-host"), c.String("prefix"), gauges, counters)
			color.White("-------------------")
			interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
	app.Run(os.Args)
}
