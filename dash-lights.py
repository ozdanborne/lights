"""lights_witch

Usage:
  lights_witch.py daemon
  lights_witch.py sniff

Description:
  daemon:   Run the light daemon.
  sniff:    Run the light daemon without light actions.
"""

import json
from docopt import docopt

from scapy.all import sniff, ARP
from phue import Bridge

import random

KNOWN_MAC_FILE = 'lights.json'

bridge = Bridge('10.0.0.2')

def get_known_macs(filename):
	f = open(filename, 'r')
	data = f.read()
	return json.loads(data)

known_macs = get_known_macs(KNOWN_MAC_FILE)

def get_light_ids():
    return [id for id in bridge.get_light()]

def mysniff(toggle_lights):
	def arp_display(pkt):
	    try:
		if pkt[ARP].op == 1: # who-has (request)
		    if True:
			mac = pkt[ARP].hwsrc
			try:
			    room = known_macs[mac]
			except KeyError:
				if not toggle_lights:
					print "Unknown Arp: " + mac
			else:
			    print "toggling room: %s" % room
                            if toggle_lights:
                                    set_lights_rainbow(get_light_ids())
		    else:
			print pkt[ARP].psrc
	    except IndexError:
		pass
        return arp_display

def set_lights_rainbow(ids):
    for id in ids:
        bridge.set_light(int(id), {'on': True,
                                   'effect': 'colorloop',
                                   'bri': 254,
                                   'hue': random.randint(0,65535),
                                   'sat': random.randint(128,254)})

if __name__ == '__main__':
	arguments = docopt(__doc__)
	if arguments['sniff']:
                print "running sniffer"
		print sniff(prn=mysniff(toggle_lights=False), filter="arp", store=0, count=0)
        if arguments['daemon']:
                print ("running daemon mode")
                print sniff(prn=mysniff(toggle_lights=True), filter="arp", store=0, count=0)

