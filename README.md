# Hegemonie

Hegemonie is an online management, strategy and diplomacy RPG. The current
repository is a reboot of what [Hegemonie](http://www.hegemonie.be) was
between 1999 and 2003. It is under heavily inactive construction.

A Web interface manages the authentication of the players, displays the
status of the country managed by each player and proposes actions (as HTML forms)
to mpact the world.

Meanwhile, the game engine is managed by a standalone daemon that makes
the world evolve with external triggers: long term actions progress a bit
toward their completion, the movements are executed, attacks started, resources
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
3. Notifications will be emitted upon special events in the game.
   No technical solution has been chosen yet.
   It is likely to be split into a collect by either [Redis](https://redis.io),
   [Kafka](https://kafka.apache.org) or [Beanstalkd](https://beanstalkd.github.io),
   and then forwarded to any IM (instant messenging) application like
   [Discord](https://discord.io/), [Slack](https://slack.com),
   [RocketChat](https://rocket.chat), [Riot](https://riot.im) or whatever.

A game instance for a small community is lightweight enough to run on a small
ARM-based board.

## Architecture

* **hege-front** serves HTTP pages for the human beings
* **hege-world** manage the game's world through a simple HTTP/JSON API
* **hege-ticker** triggers the rounds in the game's world.
* ``haproxy`` is OPTIONAL but recommanded as an SSL termination frontend.
* ``cron`` is OPTINAL but recommanded to trigger the ticker.

### Scalability

This is not the topic yet.

However there are a few opportunities:
* The front service is stateless, you might deploy many of them.
* The world service is stateful and it manages all the game entities.

