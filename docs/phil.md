# Philosophy of Singula

This is an introduction into the approach singula takes in the overall design
of the system. It is inspired by the blockchains of Bitcoin and Ethereum, but
also the Clojure and Go programming languages.

* golang as the language. The system is built on golang and assumes other nodes and clients also run golang. Go has channels as built in primtives and the protocol makes use of that.

* messages and channels. Nodes communicate with messages. Messages can have different types of flows between them, for example Request <> Reply and Subscribe <> Publish.

* modular design. Components should be as pluggable as possible. Protocols should be abstracted away so that higher level dapps can built on them.

* the consensus mechanism and networking layer should be separated. consensus is a higher level concept, building on lower level protocols.

* immutable values