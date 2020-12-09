package main

import (
	"fmt"
	"math/rand"
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
	config, err := config.Load()
	if err != nil {
		return err
	}

	bridge, err := huego.Discover()
	if err != nil {
		return fmt.Errorf("bridge discovery failed: %v", err)
	}

	if config.User == "" {
		config.User, err = bridge.CreateUser("lightcli")
		if err != nil {
			if hueErr, ok := err.(*huego.APIError); ok && hueErr.Type == 101 {
				return fmt.Errorf("not authenticated! press the hue link button then rerun")
			}
			return fmt.Errorf("failed to create user: %s", err.Error())
		}
		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to save auth to disk: %s", err.Error())
		}
	}

	b := bridge.Login(config.User)
	if config.Group == 0 {
		groups, err := b.GetGroups()
		if err != nil {
			return fmt.Errorf("failed to get groups: %v", err)
		}
		config.Group = groups[0].ID

		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to save auth to disk: %s", err.Error())
		}
	}

	switch command {
	case "on":
		b.SetGroupState(config.Group, huego.State{
			On:  true,
			Bri: uint8(255),
			Hue: 1,
			Sat: 1,
		})
	case "off":
		b.SetGroupState(config.Group, huego.State{
			On: false,
		})
	case "bright":
		b.SetGroupState(config.Group, huego.State{
			On:  true,
			Bri: 128,
		})
	case "dim":
		b.SetGroupState(config.Group, huego.State{
			On:  true,
			Bri: 12,
		})
	case "random":
		lights, err := b.GetLights()
		if err != nil {
			return err
		}
		for _, light := range lights {
			b.SetLightState(light.ID, huego.State{
				On:  true,
				Hue: uint16(rand.Uint32()),
			})
		}
	case "rainbow":
		b.SetGroupState(config.Group, huego.State{
			On:     true,
			Effect: "colorloop",
		})
	default:
		return fmt.Errorf("unrecognized command: %s", command)
	}

	return nil
}
