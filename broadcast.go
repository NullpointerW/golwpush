package GoPush

import "GoPush/pkg"

func broadcaster(broadMsg *pkg.Package) {
	for _, conn := range conns {
		conn.write(broadMsg)
	}
}
