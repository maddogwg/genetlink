module github.com/mdlayher/genetlink

go 1.25.0

require (
	github.com/google/go-cmp v0.7.0
	github.com/mdlayher/netlink v1.9.0
	golang.org/x/net v0.57.0
	golang.org/x/sys v0.47.0
)

require (
	github.com/mdlayher/socket v0.6.1 // indirect
	golang.org/x/sync v0.20.0 // indirect
)

replace github.com/mdlayher/netlink => github.com/maddogwg/netlink v0.0.0-20260801185548-7f04072a805d
