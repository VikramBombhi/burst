package io

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/vikrambombhi/burst/messages"
)

const STATUS_CLOSED = 0
const STATUS_OPEN = 1

type IO interface {
	readMessages()
	writeMessages()
	GetStatus() int
}

// TODO: use builder pattern to create io interface for both web connections and files.
type ioBuilder interface {
	setConn(conn *websocket.Conn)
	setReadChannel(fromIO chan<- messages.Message)
	buildIO() *IO
}

type WebIOBuilder struct {
	conn        *websocket.Conn
	readChannel chan<- messages.Message
}

func (w *WebIOBuilder) SetConn(conn *websocket.Conn) {
	w.conn = conn
}

func (w *WebIOBuilder) SetReadChannel(fromIO chan<- messages.Message) {
	w.readChannel = fromIO
}

func (w *WebIOBuilder) BuildIO() (IO, chan<- messages.Message, error) {
	if w.conn == nil || w.readChannel == nil {
		return nil, nil, errors.New("connection and/or read channel needs to be set")
	}
	web, toClient := createWebIO(w.conn, w.readChannel)
	return web, toClient, nil
}

type FileIOBuilder struct {
	filename    string
	readChannel chan<- messages.Message
}

func (f *FileIOBuilder) SetFilename(filename string) {
	f.filename = filename
}

func (f *FileIOBuilder) SetReadChannel(fromIO chan<- messages.Message) {
	f.readChannel = fromIO
}

func (f *FileIOBuilder) BuildIO() (IO, chan<- messages.Message, error) {
	if f.filename == "" || f.readChannel == nil {
		return nil, nil, errors.New("filename and/or read channel needs to be set")
	}
	file, toClient := createFileIO(f.filename, f.readChannel)
	return file, toClient, nil
}
