# logstats - A Simple Logfile Analyzer

This project implements a very basic parser for nginx access logs, which counts
the number of requests and sent bytes and things like that.

## Why?

I'm using this to monitor some of my web projects. I have a cronjob running this
program once every minute, letting it count how many requests hit my project and
what HTTP status codes I generated. The result is then consumed by a script that
pumps the data into InfluxDB.

It's useful for projects that have no telemetry built into them, but are proxied
behind a reverse proxy.

## How?

1. Clone the repository
2. `make` in `cmd/logstats` to compile the binary. Get the required dependencies
   as needed. Sorry, no vendoring yet.
3. Run it.

To run, copy the `config.dist.yaml`, adjust to your needs and fire away:

    $ ./logstats myconfig.yaml access.log

Make sure to pipe the log output (stderr) away to only process the resulting
JSON on stdout.

## Caveats

If you're using a custom nginx log format, this will not work for you without
adjusting the code base.

The code is probably not the most efficient. Make sure to use the `read` parameter
in the config file accordingly. Reading a 500MB log file is somewhere around
20 seconds, but is relatively memory efficient (using less than 10 MB).

## License

WTFPL
