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

	/*
		heartbeat := time.NewTicker(heartbeatInterval)
		defer heartbeat.Stop()
	*/

	for {
		select {
		case <-modifications:
			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					return err
				}

				message := &Message{scanner.Bytes()}
				if err := messages.Send(message); err != nil {
					return err
				}
			}

			/*
				case <-heartbeat.C:
					if _, err := w.Write(heartbeatMessage); err != nil {
						return
					}
					flusher.Flush()
			*/

		}
	}
	return nil
}
