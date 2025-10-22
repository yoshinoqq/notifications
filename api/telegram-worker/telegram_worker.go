package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/tovma/ruslanzadacha/api/models"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
    rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
    queueName = "notifications"
    botToken  = "7981315627:AAHRycToW7UpHq8Ptfmv9iUYLfRoTBruIVQ"
    chatID    =  1398641935
)



func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func sendToTelegram(bot *tgbotapi.BotAPI, message string) {
    msg := tgbotapi.NewMessage(chatID, message)
    _, err := bot.Send(msg)
    if err != nil {
        log.Println("Ошибка отправки в Telegram:", err)
    }
    log.Printf("Отправка сообщения в Telegram: %s", message)
}

func main() {
    bot, err := tgbotapi.NewBotAPI(botToken)
    failOnError(err, "Не удалось создать Telegram бота")

    for {
        conn, err := amqp.Dial(rabbitURL)
        if err != nil {
            log.Println("Не удалось подключиться к RabbitMQ, повтор через 5 секунд...")
            time.Sleep(5 * time.Second)
            continue
        }

        ch, err := conn.Channel()
        if err != nil {
            log.Println("Не удалось открыть канал:", err)
            conn.Close()
            time.Sleep(5 * time.Second)
            continue
        }

        queue, err := ch.QueueDeclare(
            queueName,
            true,
            false,
            false,
            false,
            nil,
        )
        if err != nil {
            log.Println("Не удалось создать очередь:", err)
            ch.Close()
            conn.Close()
            time.Sleep(5 * time.Second)
            continue
        }

        msgs, err := ch.Consume(
            queue.Name,
            "",
            true,
            false,
            false,
            false,
            nil,
        )
        if err != nil {
            log.Println("Не удалось зарегистрировать потребителя:", err)
            ch.Close()
            conn.Close()
            time.Sleep(5 * time.Second)
            continue
        }

        log.Println(" [*] Слушаем очередь:", queue.Name)


        disconnected := make(chan bool)
        go func() {
            for d := range msgs {
                var notif models.Notification
                if err := json.Unmarshal(d.Body, &notif); err != nil {
                    log.Println("Ошибка парсинга JSON:", err)
                    continue
                }
                log.Printf("Получено сообщение: %s", notif.Message)
                sendToTelegram(bot, notif.Message)
            }
            disconnected <- true
        }()

  
        <-disconnected

        log.Println("Переподключение к RabbitMQ через 5 секунд...")
        ch.Close()
        conn.Close()
        time.Sleep(5 * time.Second)
    }
}
