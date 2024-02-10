package server

//// RunServer missing godoc.
//func RunGRPCServer(app *app.App) error {
//	listen, err := net.Listen("tcp", app.Config.Host)
//	if err != nil {
//		log.Printf("Error listen: %v\n", err)
//	}
//	s := grpc.NewServer()
//	pb.RegisterAddURLServiceServer(s, &UsersServer{})
//
//	fmt.Println("Сервер gRPC начал работу")
//
//	if err := s.Serve(listen); err != grpc.ErrServerStopped {
//		log.Printf("Error starting the server: %v\n", err)
//	}
//	return nil
//}
