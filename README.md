# Alarm

Simple CLI alarm clock for Linux, MacOS and Windows.

## Usage

The `alarm` CLI supports setting alarms, listing active alarms and removing alarms from the list.

### Setting alarms

Set a one-time alarm at 6:30 in 24h format:

`$ alarm -s 6:30`

12h format is also supported by adding an AM/PM postfix.

Set a daily alarm at 10:00 AM:

`$ alarm -s 10:00AM -r daily`

Other supported recurrences include hourly and weekly.

### List alarms

Listing active alarms can be done with:

`$ alarm -l`

This will return a list of active alarms with associated IDs which can be used to remove alarms.

### Removing alarms

Remove an alarm by first listing and finding it's ID and then calling:

`$ alarm -rm <ID>`

This software was written by Vegard Løkken and is released under the GNU GPLv3 license.
