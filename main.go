package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "Ping your WebSocket endpoint with ease!")
	a.HelpFlag.Short('h')

	var (
		target       string
		writeMessage string
		pingInterval time.Duration
		verbose      bool
	)

	a.Arg("target", "Endpoint target.").
		Required().StringVar(&target)

	a.Flag("write-message", "Message to write to the server.").Short('m').
		Default("ping").StringVar(&writeMessage)

	a.Flag("ping-duration", "Period at which to send ping message to the server. Zero value means no ping at all, only dialing the server.").Short('i').
		Default("1s").DurationVar(&pingInterval)

	a.Flag("verbose", "Make the operation more talkative.").Short('v').
		Default("false").BoolVar(&verbose)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	levelAllowed := level.AllowInfo()
	if verbose {
		levelAllowed = level.AllowDebug()
	}
	l = level.NewFilter(l, levelAllowed)
	l = log.With(l, "ts", log.DefaultTimestampUTC)

	wsScheme := "ws://"
	if strings.HasPrefix(target, wsScheme) {
		wsScheme = ""
	}

	targetURL, err := url.Parse(wsScheme + target)
	if err != nil {
		level.Error(l).Log("msg", "Fail to parse url", "err", err)
		return
	}

	l = log.With(l, "target", targetURL)
	level.Debug(l).Log("msg", "Target URL parsed")

	wsConn, _, err := websocket.DefaultDialer.Dial(targetURL.String(), nil)
	if err != nil {
		level.Error(l).Log("msg", "Fail to dial WebSocket", "err", err)
		return
	}
	defer wsConn.Close()

	level.Debug(l).Log("msg", "WebSocket connected")

	done := make(chan struct{})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		defer close(done)

		for {
			messageType, message, err := wsConn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					level.Info(l).Log("msg", "WebSocket abnormally disconnected", "err", err)
				}
				return
			}

			level.Info(l).Log("msg", "Message received", "message", message, "message_type", messageType)
		}
	}()

	if pingInterval > 0 {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ticker.C:
					err := wsConn.WriteMessage(websocket.TextMessage, []byte(writeMessage))
					if err != nil {
						return
					}
				}
			}
		}()
	}

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				level.Error(l).Log("err", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			return
		}
	}
}
