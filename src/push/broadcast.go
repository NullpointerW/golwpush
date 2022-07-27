package push

func broadcaster(msg string) {
	for _, conn := range conns {
		conn.write(msg)
	}
}
