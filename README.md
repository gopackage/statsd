A "local" UDP server for gathering statistical information from infrastructure and reporting it to stats services like StatHat. Inspired by etsy/statsd.

Test using netcat:

```bash
echo -n "1=hello world" | nc -4u -w0 localhost 9045 
```