#!/usr/bin/env python3
"""lights

Usage:
  lights daemon
  lights sniff
  lights terminal
  lights rainbow
  lights off

Description:
  daemon:   Run the light daemon.
  sniff:    Run the light daemon without light actions.
"""

import json
import sys
import socket
import random

from docopt import docopt

from scapy.layers.l2 import ARP
from scapy.sendrecv import sniff
from phue import Bridge, PhueRegistrationException


KNOWN_MAC_FILE = 'lights.json'

GROUP_ID = 1
BRIGHT_SCENE_ID = 'qgqUDIXZ56KFs5L'

try:
    bridge = Bridge('philips-hue')
except socket.error:
    print ("Error: Couldn't connect to bridge. Do you have the right IP?")
    sys.exit(1)
except PhueRegistrationException:
    print ("Error: Not registered withe bridge. Hit the button!")
    sys.exit(1)
except Exception as e:
    print ("Error connecting to bridge: %s" % e)
    sys.exit(1)


def get_known_macs(filename):
    f = open(filename, 'r')
    data = f.read()
    return json.loads(data)


def get_light_ids():
    return [light_id for light_id in bridge.get_light()]

def mysniff(toggle_lights):
    known_macs = get_known_macs(KNOWN_MAC_FILE)
    def arp_display(pkt):
        try:
            if pkt[ARP].op == 1:  # who-has (request)
                mac = pkt[ARP].hwsrc
                try:
                    room = known_macs[mac]
                except KeyError:
                    if not toggle_lights:
                        print ("Unknown Arp: " + mac)
                else:
                    print ("toggling room: %s" % room)
                    if room == 'color-loop':
                        set_lights_rainbow(get_light_ids())
                    elif room == 'on-off':
                        # TODO: toggle room
                        if bridge.get_group(GROUP_ID)['state']['any_on']:
                            bridge.set_group(GROUP_ID, 'on', False)
                        else:
                            bridge.set_group(GROUP_ID, 'on', value=True)
                    else:
                        print ("unknown room: %s" % room)
        except IndexError as e:
            print (e)
    return arp_display


def set_lights_rainbow(ids):
    for id in ids:
        bridge.set_light(int(id), {'on': True,
                                   'effect': 'colorloop',
                                   'bri': 254,
                                   'hue': random.randint(0, 65535),
                                   'sat': random.randint(128, 254)})

def turn_lights_off(ids):
    for id in ids:
        bridge.set_light(int(id), {'on': False})

def set_lights_on(ids):
    for id in ids:
        bridge.set_light(int(id), {'on': True, 'bri': 254})

def set_lights_off(ids):
    for id in ids:
        bridge.set_light(int(id), {'on': False})

if __name__ == '__main__':
    args = docopt(__doc__)
    if args['sniff']:
        print ("running sniffer")
        print (sniff(prn=mysniff(toggle_lights=False), filter="arp", store=0, count=0))
    if args['daemon']:
        print (("running daemon mode"))
        print (sniff(prn=mysniff(toggle_lights=True), filter="arp", store=0, count=0))
    if args['terminal']:
        import pdb; pdb.set_trace()
    if args['rainbow']:
        set_lights_rainbow(get_light_ids())
    if args['off']:
        turn_lights_off(get_light_ids())
    else:
        print ('done')
