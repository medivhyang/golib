package http

type Midware = func(HandlerFunc) HandlerFunc
