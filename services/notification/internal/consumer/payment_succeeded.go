package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"time"

	"github.com/IBM/sarama"
	"github.com/zyncc/ecommerce-microservice/services/notification/pkg/types/topics"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

func (h *NotificationEventHandler) paymentSucceededEvent(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event topics.PaymentSucceededEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.log.Error("invalid payment event, dropping", zap.Error(err), zap.ByteString("payload", msg.Value))
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	order, err := h.orderClient.FindOrderByOrderID(ctx, event.OrderID)
	if err != nil {
		cancel()
		return err
	}
	defer cancel()

	data := welcomeEmailData{
		Name:  order.FirstName,
		Email: order.Email,
	}

	htmlString, err := renderWelcomeEmail(data)
	if err != nil {
		return err
	}

	if err := h.sendEmail(data.Email, htmlString); err != nil {
		return err
	}

	return nil
}

type welcomeEmailData struct {
	Name  string
	Email string
}

func renderWelcomeEmail(data welcomeEmailData) (string, error) {
	tmpl, err := template.ParseFiles("templates/welcome_email.html")
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}

func (h *NotificationEventHandler) sendEmail(to, html string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", h.env.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Accound Created")
	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		h.env.FromEmailSMTP,
		h.env.SMTPPort,
		h.env.FromEmail,
		h.env.SMTPPassword,
	)

	return d.DialAndSend(m)
}
