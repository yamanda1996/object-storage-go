package rabbitmq

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"object-storage-go/api-server/model"
	"strconv"
	"strings"
)

var RabbitMqDialUrl string

func GetRabbitMqDialUrl() string {

	if len(RabbitMqDialUrl) == 0 {
		var builder strings.Builder
		builder.WriteString("amqp://")
		builder.WriteString(model.Config.RabbitMqConfig.RabbitMqUser)
		builder.WriteString(":" + model.Config.RabbitMqConfig.RabbitMqPwd)
		builder.WriteString("@" + model.Config.RabbitMqConfig.RabbitMqAddress)
		builder.WriteString(":" + strconv.Itoa(model.Config.RabbitMqConfig.RabbitMqPort))
		RabbitMqDialUrl = builder.String()
	}
	return RabbitMqDialUrl
}

type RabbitMQ struct {
	channel 		*amqp.Channel
	Name 			string
	exchange 		string
}

func New(url string) *RabbitMQ {

	connection, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("connect to mq failed, dial url [%s]", url)
		return nil
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Error("create channel failed")
		return nil
	}

	queue, err := channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("declare queue failed")
		return nil
	}

	mq := new(RabbitMQ)
	mq.channel = channel
	mq.Name = queue.Name

	return mq
}

func (q *RabbitMQ) Bind(exchange string)  {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	q.exchange = exchange
}

func (q *RabbitMQ) Send(queue string, body interface{})  {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if err != nil {
		panic(err)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{})  {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if err != nil {
		panic(err)
	}
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	ch, err := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	return ch
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
}
