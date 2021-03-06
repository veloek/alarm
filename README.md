# Alarm

Simple CLI alarm clock for Linux, MacOS and Windows.

An alarm daemon is started and runs in the background as soon as you set, list
or remove alarms. This daemon runs a scheduler that triggers configured alarms.

Alarm database is stored in a file called `.alarms` in you home directory.

## Build

Prerequisites:
* If Debian/Ubuntu, run `apt install libasound2-dev`.
* If FreeBSD, run `pkg install openal-soft`.

Run `go get -u` to install dependencies, then run `go build` or `go install` to
create or install executable.

## Usage

The `alarm` CLI supports setting alarms, listing active alarms and removing
alarms from the list.

### Setting alarms

Set a one-time alarm at 6:30 in 24h format:

`$ alarm -s 6:30`

12h format is also supported by adding an AM/PM postfix.

Set a daily alarm at 10:00 AM:

`$ alarm -s 10:00AM -r daily`

Recurrence can also be hourly.

### List alarms

Listing active alarms can be done with:

`$ alarm -l`

This will return a list of active alarms with associated IDs which can be used
to remove alarms.

### Removing alarms

Remove an alarm by first listing and finding it's ID and then calling:

`$ alarm -rm <ID>`

This software was written by Vegard Løkken and is released under the GNU GPLv3
license.
