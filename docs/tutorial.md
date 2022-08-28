# Tutorial 

Singula consists of two separate layers. The netio base layer and the consensus logic on top of this.
the node reads and writes from TPC/IP and all interaction with network is "virtualized" with go channels.

Instead of reading and writing to the network directly, higher level code needs to use the channels (Ntchan). This allow to enforce good protocol use and re-use of these protocols through the node software. Instead of thinking of code in terms of separate networks, we think of it as channels and connections of these channels. This allows for drastic simplification in the logic, because we have always a separate of concerns between communication and application.

For example lets say we want to contruct a protocol which reads on network and echoes back whatever it reads. We can write code which simulates this behaviour without using any network code. We have 
a reader and writer channel and the writer channel will write whatever has been read on the reader channel.

```
func ReadProcessorEcho(ntchan Ntchan) {

	for {		
		msgString := <-ntchan.Reader_queue
		
		if len(msgString) > 0 {
			reply := "echo: " + msgString
			ntchan.Writer_queue <- reply
		}
	}

}
```

This code will achieve this echo protocol. Note that this is independent of any network logic and 
based purely on golang channels. A separate process needs to ensure that reader and writer queue are mapped to the underlying TCP/IP network.

The network always operates in terms of discrete messages which start with "{" and end with "}". If the program maintains multiple connections there are multiple channels (note: how 
multiple connections are proritized is an open question).

To accomplish this the setup requires at least 4 go routines

```
//reads from the actual "physical" network
go ReadLoop(ntchan)
//process of reads in X_in chans
go ReadProcessor(ntchan)
//processor of X_out chans
go WriteProcessor(ntchan)
//write to network whatever is in writer queue
go WriteLoop(ntchan, 300*time.Millisecond)
```

Messages that flow between nodes have the following structure. Currently everything is built on JSON messages, but this is envisioned as pluggable between different message encodings and a pluggable parser.

```
type MessageJSON struct {
	//type of message i.e. the communications protocol
	MessageType string `json:"messagetype"`
	//Specific message command
	Command string `json:"command"`
	//any data, can be empty. gets interpreted downstream to other structs
	Data *json.RawMessage `json:"data,omitempty"`	
}
```

The message type will show where messages have to be routed to, for example a request will be handled in a dedicated request in channel. It is important that all data stored can be read and written from the message format. Currently any data associated with a message will be stored in the data field with json.RawMessage type. This can be encoded and decoded on the fly.

(temporary note: Lets say a node sends a request and we want to track whether this request has been handled we need to compare replies in and check against previous request out, this would be a more sophisticated use. Another example would be we send mulitple requests to several peers and cancel a request once one has already fullfilled one of the requests)




