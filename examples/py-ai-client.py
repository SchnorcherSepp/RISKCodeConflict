#!/usr/bin/env python

"""
This file contains an example in Python for an AI controlled client.
Use this example to program your own AI in Python.
"""

import json
import socket
import time
from threading import Lock

# CONFIG
TCP_IP = '127.0.0.1'
TCP_PORT = 1234

# TCP connection
conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
conn.connect((TCP_IP, TCP_PORT))
conn_make_file = conn.makefile()

# thread lock
lock = Lock()


# ------ Helper ------------------------------------------------------------------------------------------------------ #


# command send a single command and return the response
def command(cmd):
    lock.acquire()  # <---- LOCK

    # remove protocol break
    cmd = cmd.replace('\n', '')
    cmd = cmd.replace('\r', '')

    # send command
    conn.send(bytes(cmd, 'utf8') + b'\n')
    print("SEND:", cmd)  # DEBUG !!!

    # read response
    resp = conn_make_file.readline()
    resp = resp.replace('\n', '')
    resp = resp.replace('\r', '')
    print("RESP:", resp)  # DEBUG !!!

    lock.release()  # <---- UNLOCK

    # return
    return resp


# ----------- COMMANDS ------------------------------------------------------------------------------------------------#


# player send a 'add new player' command.
def add_player(name, color_r, color_g, color_b):
    return command("PLAYER|%s|%d|%d|%d" % (name, color_r, color_g, color_b))


# status returns a json with all world data.
def world_status():
    return command("STATUS")


# attack_or_move send a move/attack command.
def attack_or_move(attacker, defender, strength):
    return command("MOVE|%s|%s|%d" % (attacker, defender, strength))


# reinforcement send a reinforcement command.
def reinforcement(country, strength):
    return attack_or_move(country, country, strength)


# end_turn send a 'end your turn' command.
def end_turn():
    return command("END")


# --------- MY AI ---------------------------------------------------------------------------------------------------- #


if __name__ == '__main__':

    # commands:
    #   add_player(name, color_r, color_g, color_b) -> error
    #   world_status() -> json
    #   attack_or_move(attacker, defender, strength) -> error
    #   reinforcement(country, strength) -> error
    #   end_turn() -> error

    # add player
    playerName = "My Mega AI 9000"
    err = add_player(name=playerName, color_r=252, color_g=3, color_b=236)
    if err != "OK":
        exit(err)

    # Main AI loop
    while True:
        time.sleep(0.30)  # Prevent server denial of service (DoS) by pacing requests.

        # Get the current state of the game world from the server.
        json_str = world_status()
        world = json.loads(json_str)

        # Check if it's the specified player's turn.
        if not world['Freeze'] and len(world['PlayerQueue']) > 1 and world['PlayerQueue'][0]['Name'] == playerName:
            print("MY TURN")
            #--------------------------------------------------------------------------------------

            # TODO: implement your AI here

            ##############################
            # useful variables and lists #
            ##############################

            # all countries
            all_countries = list(world['Countries'].values())
            my_countries = []
            my_recruiting_countries = []
            my_fortress_countries = []
            my_border_countries = []
            my_reinforcement = world['PlayerQueue'][0]['Reinforcement']

            # fill the lists
            for c in all_countries:
                if c['Occupier']['Player'] == playerName:
                    my_countries.append(c)
                    if c['RecruitingRegion']:
                        my_recruiting_countries.append(c)
                    if c['FortressRegion']:
                        my_fortress_countries.append(c)
                    if c['BorderRegion']:
                        my_border_countries.append(c)

            #################################
            # various brilliant ai commands #
            #################################

            # EXAMPLE: recruiting
            for c in my_recruiting_countries:
                err = reinforcement(c['Name'], 1)
                if err != "OK":
                    print(err)

            # EXAMPLE: move or attack
            for c in my_recruiting_countries:
                attacker = c['Name']
                defender = c['Neighbors'][0]
                err = attack_or_move(attacker, defender, 1)
                if err != "OK":
                    print(err)

            # EXAMPLE: End the turn and wait briefly before continuing.
            time.sleep(0.4)  # Sleep for 400 milliseconds
            err = end_turn()
            if err != "OK":
                print(err)
