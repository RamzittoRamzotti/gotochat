package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ActiveClients = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gotochat_active_clients_total",
		Help: "Number of currently connected WebSocket clients",
	})

	ActiveRooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gotochat_active_rooms_total",
		Help: "Number of currently active chat rooms",
	})

	MessagesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gotochat_messages_total",
		Help: "Total number of messages sent",
	})
)

func Register() {
	prometheus.MustRegister(ActiveClients, ActiveRooms, MessagesTotal)
}
