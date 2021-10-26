package server

type Option func(s *Server)

func OptionStaticPath(p string) Option {
	return func(s *Server) {
		s.static = p
	}
}

func OptionUserName(u string) Option {
	return func(s *Server) {
		s.username = u
	}
}

func OptionPassword(p string) Option {
	return func(s *Server) {
		s.password = p
	}
}

func OptionPort(p int) Option {
	return func(s *Server) {
		s.port = p
	}
}
