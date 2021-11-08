package chat

// Channel containing messages from
type Client struct {
	IncomingMessages chan []byte
}

type BroadcastGroup struct {
	messages chan []byte
	entering chan Client
	leaving  chan Client
	clients  map[Client]bool
}

func NewBroadCastGroup() *BroadcastGroup {
	return &BroadcastGroup{
		messages: make(chan []byte),
		entering: make(chan Client),
		leaving:  make(chan Client),
		clients:  make(map[Client]bool),
	}
}

func (bg *BroadcastGroup) Get(m []byte) {
	bg.messages <- m
}

func (bg *BroadcastGroup) Enter(c Client) {
	bg.entering <- c
}

func (bg *BroadcastGroup) Leave(c Client) {
	bg.leaving <- c
}

func (bg *BroadcastGroup) HandleBroadcasts() {
	for {
		select {
		// A new incoming message arrived.
		// Dispatch the message to all the clients
		case msg := <-bg.messages:
			for client := range bg.clients {
				client.IncomingMessages <- msg
			}

		// A client is entering the group
		case client := <-bg.entering:
			bg.clients[client] = true

		// A client is leaving the group
		// Close all incoming messages and delete from group
		case client := <-bg.leaving:
			close(client.IncomingMessages)
			delete(bg.clients, client)
		}
	}
}
