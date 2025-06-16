# TWAMP

A minimal [TWAMP Light](https://datatracker.ietf.org/doc/html/rfc5357) implementation in Go for active round-trip time (RTT) measurement over UDP. It follows the TWAMP-Test packet format from RFC 5357 §4.1.2 and uses NTP-style 64-bit timestamps.

## Features

- Sends and receives TWAMP-Test packets over UDP
- NTP-based timestamps for accurate RTT
- Summary statistics including packet loss, min/max/avg RTT, and jitter
- Simple `Sender` and `Reflector` interfaces
- Optional Prometheus metrics instrumentation
- Configurable padding for packet size control
- Unit tests and fuzz testing

## Example

### Reflector (Listener)

```go
go func() {
    listener, _ := twamp.NewListener(":9000", 512, nil, nil)
    _ = listener.Run()
}()
```

### Sender (Client)

```go
sender, _ := twamp.NewSender(":0", "127.0.0.1:9000", time.Second, nil, nil)
res, _ := sender.SendProbe()
fmt.Printf("RTT = %v (seq %d)\n", res.RTT, res.Seq)
```

## Packet Flow Diagram (TWAMP Light)

```text
+---------+                                 +------------+
| Sender  | -- TWAMP-Test packet (UDP) -->  | Reflector  |
|         |                                 |            |
|         | <-- echoed with timestamp ----- |            |
+---------+                                 +------------+
```

## Interfaces

```go
type Sender interface {
	SendProbe() (*ProbeResult, error)
	SendProbeWithPadding(paddingLen int) (*ProbeResult, error)
	Summary() *ProbeSummary
	Close() error
}

type Reflector interface {
    HandleProbe(msg []byte, from net.Addr) ([]byte, error)
}

type Listener interface {
    Run() error
    Close() error
    LocalAddr() net.Addr
}
```

## TWAMP-Test Packet Format

This library uses the minimal TWAMP-Test format (RFC 5357 §4.1.2):

| Field                | Offset | Size (bytes) |
| -------------------- | ------ | ------------ |
| Sequence Number      | 0      | 4            |
| Timestamp (seconds)  | 4      | 4            |
| Timestamp (fraction) | 8      | 4            |
| Error Estimate       | 12     | 2            |
| MBZ (reserved)       | 14     | 2            |
| Reflector TS (sec)   | 16     | 4            |
| Reflector TS (frac)  | 20     | 4            |
| Padding              | 24     | variable     |

Timestamps follow NTP 64-bit format: 32 bits for seconds, 32 bits for fractional seconds since Jan 1, 1900 UTC.

## Output Example

```
RTT = 212.400µs (seq 4)
```

## Limitations (Not Yet Implemented)

This library implements the essentials of TWAMP Light. Some features are not yet included:

- DSCP / TOS field configuration
- Don't Fragment (DF) bit
- CLI tools (`twamp-client`, `twamp-reflector`)
