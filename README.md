GoBot based Melatonin generator for Gophers

Toy project for GopherCon 2016 - GoBot Hackathon, Wed Jul 13, 2016

Beyond the click-bait, it uses the MacOSX App [f.lux](https://justgetflux.com/) via a Go server (running on that Mac) controlled by a GoBot enabled Go App running on an Edison. The app controls the color temperature of the Mac helping the Gopher naturally generate more Melatonin!! :P. Sorry, if you were looking for the actual chemical, you would have to [buy it yourself](https://g.co/kgs/PJxBtD)

Currently: The button toggles the screen between 2700K (Tungsten) and 6500 (Day light)

TODO:

* Lots of cleanup, removed hardcoded IPs etc.
* Program the rotary so that screen color temperature changes according to the rotary state. Sub todo: Unlike a button, rotary may emit the state/value too many times / too frequently. We need to control / pace the actual action on the Mac.
* Put the light sensor next to a daylight source (say a window) via a really long cable and see if will help control the color based on real day light instead of the current time-based setting available on the Mac.


