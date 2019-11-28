module github.com/orbs-network/orbs-network-javascript-plugin

go 1.12

replace github.com/orbs-network/orbs-network-go => ../orbs-network-go

replace github.com/ry/v8worker2 => ./vendor/github.com/ry/v8worker2

require (
	github.com/VividCortex/ewma v1.1.1
	github.com/beevik/ntp v0.2.0
	github.com/c9s/goprocinfo v0.0.0-20190309065803-0b2ad9ac246b
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd
	github.com/ethereum/go-ethereum v1.9.6
	github.com/google/go-cmp v0.3.1
	github.com/orbs-network/go-mock v1.1.0
	github.com/orbs-network/govnr v0.2.0
	github.com/orbs-network/lean-helix-go v0.2.4
	github.com/orbs-network/membuffers v0.4.0
	github.com/orbs-network/orbs-client-sdk-go v0.12.0
	github.com/orbs-network/orbs-contract-sdk v1.4.0
	github.com/orbs-network/orbs-network-go v0.0.0
	github.com/orbs-network/orbs-spec v0.0.0-20191114152037-24b26e24030e
	github.com/orbs-network/scribe v0.2.3
	github.com/pkg/errors v0.8.1
	github.com/ry/v8worker2 v0.0.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
)
