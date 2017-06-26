package types

type Subscriber interface {
	Subscribe(string) chan interface{}
	Unsubscribe(chan interface{})
}

type Publisher interface {
	PublishWithTags(interface{}, map[string]interface{})
}

func PublishEventNewBlock(pubsub Publisher, block *Block) {
	pubsub.PublishWithTags(TMEventData{EventDataNewBlock{block}}, map[string]interface{}{"tm.events.type": EventStringNewBlock()})
}
