package main

import (
	"fmt"
	"os"

	"github.com/ozdanborne/lights/pkg/config"

	"github.com/amimof/huego"
)

const hueIP = "192.168.1.161"

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`usage: lights <command>

commands:
  on
  off
  dim
  rainbow`)
		os.Exit(1)
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(command string) error {
	config := config.Load()

	bridge, err := huego.Discover()
	if err != nil {
		return fmt.Errorf("bridge discovery failed: %v", err)
	}

	// bridge.DeleteUser("lightci")
	if config.User == "" {
		config.User, err = bridge.CreateUser("lightcli")
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to persist user: %v", err)
		}
	}

	b := bridge.Login(config.User)

	groups, err := b.GetGroups()
	if err != nil {
		return fmt.Errorf("failed to get groups: %v", err)
	}

	if len(groups) != 1 {
		return fmt.Errorf("sorry, only 1 group supported currently")
	}

	group := groups[0]

	switch command {
	case "on":
		b.SetGroupState(group.ID, huego.State{
			On:  true,
			Bri: uint8(255),
			Hue: 1,
			Sat: 1,
		})
	case "off":
		b.SetGroupState(group.ID, huego.State{
			On: false,
		})
	case "dim":
		b.SetGroupState(group.ID, huego.State{
			On:  true,
			Bri: 12,
		})
	case "rainbow":
		b.SetGroupState(group.ID, huego.State{
			On:     true,
			Effect: "colorloop",
		})
	default:
		return fmt.Errorf("unrecognized command: %s", command)
	}

	return nil
}
