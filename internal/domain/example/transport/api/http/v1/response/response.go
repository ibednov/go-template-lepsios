package response

import "github.com/ibednov/go-template-lepsios/internal/domain/example/entity"

type Hello struct {
	Text string `json:"text"`
}

func FromMessage(m entity.Message) Hello {
	return Hello{Text: m.Text}
}
