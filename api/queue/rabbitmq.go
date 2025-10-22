package queue

import (
	"log"
	"os"

	"time"


	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp091.Connection, *amqp091.Channel, error) {
    url := os.Getenv("RABBITMQ_URL")
    if url == "" {
        url = "amqp://guest:guest@rabbitmq:5672/"
    }

    var conn *amqp091.Connection
    var err error
    for i := 0; i < 10; i++ {
        conn, err = amqp091.Dial(url)
        if err == nil {
            break
        }
        log.Println("RabbitMQ не готов, пробуем снова через 3 секунды...")
        time.Sleep(3 * time.Second)
    }
    if err != nil {
        return nil, nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, nil, err
    }

    return conn, ch, nil
}

func DeclareQueue(ch *amqp091.Channel, queueName string) (amqp091.Queue, error) {
    q, err := ch.QueueDeclare(
        queueName,
        true,  
        false, 
        false, 
        false, 
        nil,   
    )
    if err != nil {
        return amqp091.Queue{}, err
    }

    log.Println("Очередь создана:", queueName)
    return q, nil
}



func PublishMessage(ch *amqp091.Channel, queueName string, body []byte) error {

    log.Printf("Публикация сообщения в RabbitMQ: %s", string(body))
    return ch.Publish(
        "",        
        queueName, 
        false,     
        false,     
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}


