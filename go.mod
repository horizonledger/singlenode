module singula.finance/node

go 1.19

replace singula.finance/netio => ../netio

require (
	github.com/btcsuite/btcd v0.22.0-beta
	singula.finance/netio v0.0.0-00010101000000-000000000000
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp v0.0.0-20220826205824-bd9bcdd0b820 // indirect
)
