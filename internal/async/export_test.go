package async

// IsClosed checks if a channel is entirely closed or not.
// This function is only exported for testing.
func IsClosed(ch *UnboundedStructChan) bool {
	select {
	case <-ch.close:
		return true
	default:
		return false
	}
}
