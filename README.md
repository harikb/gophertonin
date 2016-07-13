GoBot based Melatonin generator for Gophers

Toy project for GopherCon 2016 - GoBot Hackathon, Wed Jul 13, 2016

Beyond the click-bait, it uses the MacOSX App f.lux via a Go server (running on that Mac) controlled by a GoBot enabled Go App running on an Edison. The app controls the color temperature of the Mac helping the Gopher naturally generate more Melatonin!! :P.

Currently: The button toggles the screen between 2700K (Tungsten) and 6500 (Day light)

TODO:

* Program the rotary so that screen color temperature changes according to the rotary state. Sub todo: Unlike a button, rotary may emit the state/value too many times / too frequently. We need to control / pace the actual action on the Mac.
* Put the light sensor next to a daylight source (say a window) via a really long cable and see if will help control the color based on real day light instead of the current time-based setting available on the Mac.


