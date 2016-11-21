xinput
======

XInput (Xbox 360/Xbox One compatible controller) library for [Go][] on Windows.

Based on: https://github.com/tajtiattila/xinput

Modified version that combines analog input into digital and gives access to modify deadzones.

[Go]: http://golang.org

Requirements
------------

- Windows 8.1: xinput1_4.dll
- Windows 7: xinput1_3.dll
- Windows Vista: xinput9_1_0.dll

Tested on Windows 10 x64.

Installation
------------

`go get github.com/Nudua/xinput`

Usage
-----

`IsLoaded` checks if XInput was loaded successfully.

`GetState` retrieves the state thumbsticks, d-pad, buttons, triggers and digital versions of the analog input.

`GetSimpleState` retrieves the state thumbsticks, d-pad, buttons and triggers.

`SetState` sets vibration state.

Deadzones and thresholds for analog inputs
---------
`LEFT_THUMB_DEADZONE` default '7849' (-32767 to 32767)

`RIGHT_THUMB_DEADZONE` default '8689' (-32767 to 32767)

`TRIGGER_TRESHOLD` default '50' (0 to 255)

Example
-------
TODO