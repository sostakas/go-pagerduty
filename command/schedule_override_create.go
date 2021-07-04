package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type ScheduleOverrideCreate struct {
	Meta
}

func ScheduleOverrideCreateCommand() (cli.Command, error) {
	return &ScheduleOverrideCreate{}, nil
}

func (c *ScheduleOverrideCreate) Help() string {
	helpText := `
	pd schedule-override create <SERVICE-ID> <FILE> Create a schedule override from json file
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleOverrideCreate) Synopsis() string {
	return "Create an override for a specific user"
}

func (c *ScheduleOverrideCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("schedule-override create")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	client.SetDebugFlag(pagerduty.DebugCaptureLastResponse)
	var o pagerduty.Override
	if len(flags.Args()) != 2 {
		log.Error("Please specify input json file")
		return -1
	}
	log.Info("service id is:", flags.Arg(0))
	log.Info("Input file is:", flags.Arg(1))
	f, err := os.Open(flags.Arg(1))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&o); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", o)
	o1, err := client.CreateOverride(flags.Arg(0), o)
	if err != nil {
		// log.Error(err)

		fmt.Printf("\n\nerr: %s\n\n", err)

		resp, ok := client.LastAPIResponse()
		if ok {
			fmt.Println("resp:")
			fmt.Printf("resp.Status: %s\n", resp.Status)
			fmt.Printf("resp.Header: %#v\n", resp.Header)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("resp.Body:\n%s\n", string(body))
		}
		return -1
	}
	log.Println("New override id:", o1.ID)
	return 0
}
