# Hegemonie

Hegemonie is an online management, strategy and diplomacy RPG. The current
repository is a reboot of what [Hegemonie](http://www.hegemonie.be) was
between 1999 and 2003. It is under heavily inactive construction.

A Web interface...
* manages the authentication of the players
* displays the status of the country managed by each player
* proposes actions to evolve in the world.

Meanwhile, the game engine is managed by a standalone daemon that makes
the world evolve at periodical ticks: long term actions progress a bit toward
their completion, the movements are executed, attacks started, resources
produced, etc etc.

Technical facts:
1. Written in 100% in Golang: for the sake of Simplicity and Portability. The
   code mostly depends on [Go-Macaron](https://go-macaron.com). At the moment
   no special attention has been paid to the performance of the whole thing:
   this will happend after the release of a very first MVP.
2. No database required: the system has all its components in RAM while it is
   alive, it periodically persist its state and restore it at the startup.
   The status is written in [JSON](https://json.org) to ease the daily
   administration.
3. Notifications are emitted upon special events in the game, they are managed
   by a [Redis](https://redis.io) service.

A game instance for a small community is lightweight enough to run on a small
ARM-based board.

Architecture:
* ``haproxy`` as an SSL termination frontend.
* **hege-front** serves HTTP pages for the human beings
* **hege-world** manage the game's world through a simple HTTP/JSON API
* **hege-ticker** triggers the rounds in the game's world.
* ``redis-server``

The configuration is very simple:
* A central configuration file for **hege-front**, **hege-world** and
  **hege-ticker**
* A standard configuration file for **redis**
* The ``crontab`` configuration to trigger the ticker.

