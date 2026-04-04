# OnIt

small status app for running in a Getty-like interface (800x480 screens, ideally)

starts a small API with three endpoints:

- /alarm?time=[Linux time UTC]&title=[string]&description=[string]&alert=[true]
- /display?text=[string]
- /clear: clears alarm and display

presents:

- time locally and in Hong Kong (where most of my colleagues work)
- an alarm time if set, with a description (and will flash when reached)
- anything else sent via the display API
