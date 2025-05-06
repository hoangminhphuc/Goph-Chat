package pubsub

import (
	"sync"

	"github.com/hoangminhphuc/goph-chat/common/logger"
)

type LocalPubSub struct {
	name         string
	messageQueue chan *Message // Messages flow from Publish → messageQueue → subscribers.
	mapTopic     map[Topic][]chan *Message // Slice of subscribers for each topic.
	lock         *sync.RWMutex
	logger       logger.ZapLogger
}

func NewLocalPubSub(name string) *LocalPubSub {
	return &LocalPubSub{
		name:         name,
		messageQueue: make(chan *Message, 10000),
		mapTopic:     make(map[Topic][]chan *Message),
		lock :         new(sync.RWMutex),
		logger:       logger.NewZapLogger(),
	}
}

func (l *LocalPubSub) Publish(data *Message) {
	go func() {
		l.messageQueue <- data // Send data to message queue
	} ()
}

func (l *LocalPubSub) Subscribe(topic Topic) (ch <-chan *Message, unsubscribe func()) {
	c := make(chan *Message)

	l.lock.Lock()
	// Even if topic not exists, it will return empty slice
	l.mapTopic[topic] = append(l.mapTopic[topic], c)

	l.lock.Unlock()

	l.logger.Log.Info("Registered topic ", topic, " successfully")

	return c, func() {
		if chans, ok := l.mapTopic[topic]; ok {
			for i, ch := range chans {
				if ch == c {
					chans = append(chans[:i], chans[i+1:]...)

					l.lock.Lock()
					l.mapTopic[topic] = chans
					l.lock.Unlock()

					close(c)
					break
				}
			}
		}
	}
}

func (l *LocalPubSub) Name() string {
	return l.name
}
func (l *LocalPubSub)	InitFlags() {
	// No flags needed for now
}

func (l *LocalPubSub)	Run() error {
	go func() {
		for {
			data, ok := <-l.messageQueue // Get data from message queue
			if !ok {
				l.logger.Log.Info("Message queue closed, stopping pubsub dispatcher")
				return
			}
			if data == nil {
				l.logger.Log.Warn("Nil message received, skipping")
				continue
			}
			
			l.lock.RLock()

			// Publish data to all subscribers of the topic
			if chans, ok := l.mapTopic[data.GetTopic()]; ok {
				for _, ch := range chans {
					go func(c chan *Message) {
						c <- data
					}(ch)
				}
			}

			l.lock.RUnlock()
		}
	}()
	return nil
}
func (l *LocalPubSub) Stop() <-chan error {
	errc := make(chan error, 1)
	go func() {
		close(l.messageQueue)

		l.lock.Lock()
		for topic, chans := range l.mapTopic {
			for _, ch := range chans {
				close(ch)
			}
			delete(l.mapTopic, topic)
		}
		l.lock.Unlock()

		l.logger.Log.Info("Local PubSub stopped successfully")
		errc <- nil
	}()
	return errc
}

func (l *LocalPubSub) Get() interface{} {
	return l
}