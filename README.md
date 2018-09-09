# wsping

## Installation

```bash
go get -u github.com/ariefrahmansyah/wsping
```

## Usage

```bash
$ wsping -h
usage: wsping [<flags>] <target>

Ping your WebSocket endpoint with ease!

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
  -m, --write-message="ping"  Message to write to the server.
  -i, --ping-duration=1s      Period at which to send ping message to the server. Zero value means no ping at all.
  -v, --verbose               Make the operation more talkative.

Args:
  <target>  Endpoint target.
```
