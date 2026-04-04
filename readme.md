# OnIt

small status app for running in a getty-like interface (800x480 screens, ideally)

starts a small api with three endpoints:

- /alarm?time=[linuxtimeutc]&title=[string]&description=[string]&alert=[true]
- /display?text=[string]
- /clear: clears alarm and display

presents:

- time locally and in hong kong (where most of my colleagues work)
- an alarm time if set, with a description (and will flash when reached)
- anything else sent via the display api
