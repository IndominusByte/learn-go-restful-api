package endpoint_http

import (
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/config"
	"github.com/creent-production/cdk-go/mail"
	"github.com/creent-production/cdk-go/queue"
	"github.com/creent-production/cdk-go/response"
	"github.com/go-chi/chi/v5"
)

func AddEmail(r *chi.Mux, cfg *config.Config, m *mail.Mail) {
	r.Route("/email", func(r chi.Router) {
		r.Post("/send", func(rw http.ResponseWriter, r *http.Request) {
			q := queue.NewQueue(func(val interface{}) {
				m.SendEmail(
					[]string{"/app/static/icon-categories/default.jpg"},
					[]string{"nyomanpradipta120@gmail.com", "pradipta.nyoman@tokopedia.com", "mangokky@gmail.com", "paulusbsimanjuntak@gmail.com"},
					"Activated User",
					"dont-reply@example.com",
					"templates/email/EmailConfirm.html",
					struct{ Name string }{Name: "oman"},
				)
			}, 20)
			q.Push("send")

			response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
				"_app": "Email notification has sended.",
			})
		})
	})
}
