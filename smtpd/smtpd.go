// Copyright 2024 The OpenMail Authors
// SPDX-License-Identifier: Apache-2.0

package smtpd

import (
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"time"

	"github.com/linuxuser586/openmail/telemetry"
)

var ErrUnknownCommand = errors.New("unknown command")

// Network listener enum
type Network string

const (
	Tcp        Network = "tcp"
	Tcp4       Network = "tcp4"
	Tcp6       Network = "tcp6"
	Unix       Network = "unix"
	Unixpacket Network = "unixpacket"
)

// Command enum
type Command string

const (
	Unknown       Command = ""
	ExtendedHello Command = "EHLO"
	Hello         Command = "HELO"
	Mail          Command = "MAIL"
	Recipient     Command = "RCPT"
	Data          Command = "DATA"
	Reset         Command = "RSET"
	NoOperation   Command = "NOOP"
	Quit          Command = "QUIT"
	Verify        Command = "VRFY"
)

func (c Command) String() string {
	return string(c)
}

// ToCommand converts a string to a command enum
func ToCommand(s string) (Command, error) {
	// TODO finish implementation
	switch strings.ToUpper(s) {
	case ExtendedHello.String():
		return ExtendedHello, nil
	case Quit.String():
		return Quit, nil
	default:
		return Unknown, ErrUnknownCommand
	}
}

// Config contains read only configurations for smtpd
type Config interface {
	// Network type
	Network() Network
	// Host to listen on
	Host() string
	// Port to listen on
	Port() int
	// InitialTimeout in seconds
	InitialTimeout() int
	// MailCmdTimeout in seconds
	MailCmdTimeout() int
	// RecipientCmdTimeout in seconds
	RecipientCmdTimeout() int
	// DataInitTimeout in seconds
	DataInitTimeout() int
	// DataBlockTimeout in seconds
	DataBlockTimeout() int
	// DataTerminationTimeout in seconds
	DataTerminationTimeout() int
	// HostName for this instance
	HostName() string
}

// Repository manages persistence
type Repository interface {
}

// TextConnection works with [net/textproto]
type TextConnection struct {
	conn net.Conn
}

func (t *TextConnection) Read(p []byte) (n int, err error) {
	return t.conn.Read(p)
}

func (t *TextConnection) Write(p []byte) (n int, err error) {
	return t.conn.Write(p)
}

func (t *TextConnection) Close() error {
	return t.conn.Close()
}

type Smtpd struct {
	config     Config
	telemetry  telemetry.Telemetry
	repository Repository
}

func New(config Config, telemetry telemetry.Telemetry, repository Repository) Smtpd {
	return Smtpd{config: config, telemetry: telemetry, repository: repository}
}

func (s *Smtpd) ListenAndServe() error {
	log := s.telemetry.Logger()
	listen, err := net.Listen(string(s.config.Network()), fmt.Sprintf("%s:%d", s.config.Host(), s.config.Port()))
	if err != nil {
		return err
	}
	defer listen.Close()
	log.Infof("smtpd started on %s", listen.Addr().String())

	return s.serve(listen)
}

func (s *Smtpd) serve(listen net.Listener) error {
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		config := s.config
		telemetry := s.telemetry
		repository := s.repository
		go handleConn(config, telemetry, repository, conn)
	}
}

func handleConn(
	config Config,
	telemetry telemetry.Telemetry,
	repository Repository,
	conn net.Conn) {
	textConn := TextConnection{conn: conn}
	h := handler{
		config:     config,
		telemetry:  telemetry,
		repository: repository,
		textConn:   textConn,
		protoConn:  textproto.NewConn(&textConn),
	}
	h.handleConn()
}

type handler struct {
	config     Config
	telemetry  telemetry.Telemetry
	repository Repository
	textConn   TextConnection
	protoConn  *textproto.Conn
}

func (h *handler) handleConn() {
	tp := textproto.NewConn(&h.textConn)
	conn := h.textConn.conn
	log := h.telemetry.Logger()
	log.Infof("request client address %s", conn.RemoteAddr().String())
	tp.PrintfLine(fmt.Sprintf("220 %s ESMTP OpenMail", h.config.HostName()))
	conn.SetDeadline(time.Now().Add(time.Duration(h.config.InitialTimeout()) * time.Second))
	for {
		line, err := tp.ReadLine()
		if err != nil {
			if h.handleReadTimeout(err, tp) {
				return
			}
			log.Error(err.Error())
			// TODO respond to client with appropriate error message
			continue
		}

		log.Debug(line)

		cmd, err := ToCommand(line)
		if err != nil {
			tp.PrintfLine("500 5.5.2 Error: command not recognized")
			continue
		}
		if err = h.handleCommand(cmd); err != nil {
			// TODO use custom error
			tp.PrintfLine(err.Error())
		}
	}
}

func (h handler) handleReadTimeout(err error, tp *textproto.Conn) bool {
	if netErr := err.(*net.OpError); netErr != nil {
		if netErr.Timeout() {
			conn := h.textConn.conn
			// reset the deadline so that the message can be sent to the client
			conn.SetDeadline(time.Now().Add(5 * time.Second))
			tp.PrintfLine("421 connection timeout")
			if err := h.textConn.Close(); err != nil {
				h.telemetry.Logger().Error(err.Error())
			}
			return true
		}
	}
	return false
}

func (h handler) handleCommand(cmd Command) error {
	log := h.telemetry.Logger()
	switch cmd {
	case Quit:
		if err := h.textConn.Close(); err != nil {
			log.Error(err.Error())
		}
		return nil
	default:
		log.Infof("TODO: %s not implemented", cmd.String())
	}

	return nil
}
