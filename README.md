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
  -i, --ping-duration=1s      Period at which to send ping message to the server. Zero value means no ping at all, only dialing the server.
  -v, --verbose               Make the operation more talkative.

Args:
  <target>  Endpoint target.

$ wsping localhost:8080/hi
level=info ts=2018-09-09T18:02:22.059863439Z target=ws://localhost:8080/hi msg="Message received" message=Hi message_type=1
level=info ts=2018-09-09T18:02:23.061861151Z target=ws://localhost:8080/hi msg="Message received" message=Hi message_type=1
level=info ts=2018-09-09T18:02:24.06182165Z target=ws://localhost:8080/hi msg="Message received" message=Hi message_type=1
level=info ts=2018-09-09T18:02:25.061694477Z target=ws://localhost:8080/hi msg="Message received" message=Hi message_type=1
level=info ts=2018-09-09T18:02:26.061771817Z target=ws://localhost:8080/hi msg="Message received" message=Hi message_type=1
```
