package example

import (
	socketio "github.com/doquangtan/gofiber-socket.io"
	"github.com/gofiber/fiber/v2"
)

func socketIoRoute(app fiber.Router) {
	io := socketio.New()

	io.OnAuthorization(func(params map[string]string) bool {
		// auth, ok := params["Authorization"]
		// if !ok {
		// 	return false
		// }
		return true
	})

	io.OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)

		socket.On("test", func(event *socketio.EventPayload) {
			socket.Emit("test", event.Data...)
		})

		socket.On("disconnecting", func(event *socketio.EventPayload) {
			println(event.SID, "disconnecting")
		})

		socket.On("disconnect", func(event *socketio.EventPayload) {
			println(event.SID, "disconnect")
		})
	})

	io.Of("/admin").OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)
	})

	app.Use("/", io.Middleware)
	app.Route("/socket.io", io.Server)
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Route("/", socketIoRoute)

	app.Listen(":3000")
}
