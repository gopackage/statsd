A "local" UDP server for gathering statistical information. Unlike other [statsd](https://github.com/statsd/statsd) style servers, this goal for this statsd implementation is to create an all-in-one stats solution. There is no need for a "backend" server and statistics are viewed using graphs displayed by the included web server.

Statistical data should be submitted to the server using standard [statsd types](https://github.com/statsd/statsd/blob/master/docs/metric_types.md).

To report a `1` in the `count` stat with a sampling rate of .1:

```
count:1|c|@.1
```

Sampling rate is optional and assumed to be `1`.

Timing is reported in a similar way (stat is foo, time is 320, sample rate of .1):

```
foo:320|ms|@0.1
```

Gauges store current value of an item (stat here is `bar`, value is 12)

```
bar:12|g
```

You can also increment gauges by putting a `+` or `-` before the gauge value.

Finally, sets count unique occurrences of a statistic (stat in the example is `hello`, value is `42`):

```
hello:42|s
```

Test using netcat:

```bash
echo -n "howdy:1|c" | nc -4u -w0 localhost 9045 
```

Run `statsd --help` to see all available options.

## Developer

Build using `go build _app/statsd/statsd.go`

