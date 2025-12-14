package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In

	Lc       fx.Lifecycle
	Log      *zap.Logger
	Commands []commands.Command `group:"commands"`
}

type Server struct {
	listeners net.Listener
	log       *zap.Logger
	handlers  map[string]commands.Command
}

func (s *Server) Serve(listener net.Listener) {
	s.listeners = listener

	for {
		connection, err := s.listeners.Accept()
		if err != nil {
			s.log.Error("Error accepting connection", zap.Error(err))
			return
		}

		s.log.Debug("Accepting new client", zap.String("address", connection.RemoteAddr().String()))
		go s.handleClient(connection)
	}
}

func (s *Server) handleClient(connection net.Conn) {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			s.log.Error("Couldn't close the connection")
		}
	}(connection)

	session := client.NewSession(connection)

	for {
		cmd, err := session.ParseCommand()

		if err == io.EOF {
			s.log.Info("Client disconnected", zap.String("client", connection.RemoteAddr().String()))
			return
		} else if err != nil {
			s.log.Error("Can't parse command", zap.Error(err))
			return
		}

		handler, ok := s.handlers[strings.ToUpper(cmd[0])]
		if !ok {
			s.log.Info("Unknown command", zap.String("command", cmd[0]), zap.String("client", connection.RemoteAddr().String()))
			err = session.SendString(fmt.Sprintf("Unknown command: %s", cmd[0]))
			if err != nil {
				s.log.Error("Couldn't write to client", zap.Error(err))
				return
			}
		} else {
			err = handler.Handle(session, cmd)
			if err != nil {
				s.log.Error("Failed to handle command", zap.String("command", handler.Mnemonic()), zap.Error(err))
			}
		}
	}
}

func (s *Server) Shutdown() error {
	s.log.Debug("Stopping TCP server...")
	if err := s.listeners.Close(); err != nil {
		s.log.Error("Failed to shutdown TCP binding", zap.Error(err))
		return err
	}

	return nil
}

func NewServer(p ServerParams) *Server {
	addr := "0.0.0.0:6379"

	p.Log.Debug(fmt.Sprintf("Received %d commands", len(p.Commands)))

	handlers := make(map[string]commands.Command)

	for _, cmd := range p.Commands {
		handlers[cmd.Mnemonic()] = cmd
	}

	srv := &Server{
		log:      p.Log,
		handlers: handlers,
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			srv.log.Debug("Binding TCP server", zap.String("addr", addr))
			listener, err := net.Listen("tcp", addr)

			if err != nil {
				srv.log.Error("Failed to bind TCP server", zap.Error(err))
				return err
			}

			go srv.Serve(listener)
			srv.log.Info("Server up & running")
			return nil
		},

		OnStop: func(ctx context.Context) error {
			return srv.Shutdown()
		},
	})

	return srv
}
