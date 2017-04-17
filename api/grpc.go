package api

type GrpcApi struct {
}

func (s *GrpcApi) OpenStream(arguments *OpenStreamArguments, messages StreamApi_OpenStreamServer) error {
	stream, err := router.GetStream(arguments.Id)
	if err != nil {
		return err
	}

	scanner, err := stream.NewScanner(arguments.Offset)
	if err != nil {
		return err
	}
	defer scanner.Close()

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		message := &Message{scanner.Bytes()}
		if err := messages.Send(message); err != nil {
			return err
		}
	}

	watcher := stream.NewWatcher()
	defer watcher.Close()
	modifications := watcher.Watch()

	for range modifications {
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}

			message := &Message{scanner.Bytes()}
			if err := messages.Send(message); err != nil {
				return err
			}
		}
	}
	return nil
}
