A "local" UDP server for gathering statistical information from infrastructure and storing it in a Redis Stream.

Each UDP packet containing a value should be in the format:

```
#1.3=foo
```

And for counts:

```
1=bar
```

Test using netcat:

```bash
echo -n "1=hello world" | nc -4u -w0 localhost 9045 
```

The stats server does basic buffering of data and uses batch updates to reduce load on Redis.

The stats server is configured via Redis and currently supports the following settings:

*
*
