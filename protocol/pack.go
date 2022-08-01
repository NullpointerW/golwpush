package protocol

func Pack(msg string) (b []byte) {
	return []byte(msg + string(EndFlag))
}
