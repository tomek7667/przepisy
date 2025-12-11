package przepisy

func (s *Server) SetupRoutes() {
	s.PostLogin()
	s.AddUsersRoutes()
}
