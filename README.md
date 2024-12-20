# gofiber-socket.io

- gofiber-socket.io is library an implementation of [Socket.IO](http://socket.io) in Golang, which is a realtime application framework.
- It using with web-framework [Go Fiber](https://gofiber.io)
- This library support socket.io-client version 3, 4 and only support websocket transport

# Contents

- [Install](#install)
- [Example](#example)

# Install

Install the package with:

```bash
go get github.com/doquangtan/gofiber-socket.io
```

Import it with:

```go
import "github.com/doquangtan/gofiber-socket.io"
```

and use `socketio` as the package name inside the code.

# Documents
## Server
### Constructor
#### socketio.New
```go
import (
	socketio "github.com/doquangtan/gofiber-socket.io"
	"github.com/gofiber/fiber/v2"
)

func socketIoRoute(app fiber.Router) {
	io := socketio.New()

	io.OnConnection(func(socket *socketio.Socket) {
		// ...
	})

	app.Use("/", io.Middleware)
	app.Route("/socket.io", io.Server)
}

func main() {
	app := fiber.New(fiber.Config{})
	app.Route("/", socketIoRoute)
	app.Listen(":3000")
}
```

### Events
#### Event: 'connection'
```go
io.OnConnection(func(socket *socketio.Socket) {
	// ...
})
```

### Methods
#### server.emit(eventName[, ...args])
```go
io.Emit("hello")
```
```go
io.Emit("hello", 1, "2", map[string]interface{}{"3": 4})
```

#### server.of(nsp)
```go
adminNamespace := io.Of("/admin")

adminNamespace.OnConnection(func(socket *socketio.Socket) {
	// ...
})
```

#### server.to(room)
```go
io.To("room-101").Emit("hello", "world")
```

#### server.Sockets()
```go
sockets := io.Sockets()
```

## Socket
### Methods
#### socket.on(eventName, callback)
Register a new handler for the given event.
```go
socket.On("news", func(event *socketio.EventPayload) {
	print(event.Data)
})
```
with several arguments
```go
socket.On("news", func(event *socketio.EventPayload) {
	if len(event.Data) > 0 && event.Data[0] != nil {
		print(event.Data[0])
	}
	if len(event.Data) > 1 && event.Data[1] != nil {
		print(event.Data[1])
	}
	if len(event.Data) > 2 && event.Data[2] != nil {
		print(event.Data[2])
	}
})
```
or with acknowledgement
```go
socket.On("news", func(event *socketio.EventPayload) {
	if event.Callback != nil {
		(*event.Callback)("hello", map[string]interface{}{
			"Test": "ok",
		})
	}
})
```

#### socket.join(room)
Adds the socket to the given room or to the list of rooms.
```go
io.Of("/test").OnConnection(func(socket *socketio.Socket) {
	socket.Join("room 237")

	io.To("room 237").Emit("a new user has joined the room")
})
```

#### socket.leave(room)
Removes the socket from the given room.
```go
io.Of("/test").OnConnection(func(socket *socketio.Socket) {
	socket.Leave("room 237");

	io.To("room 237").Emit("the user has left the room")
})
```
Rooms are left automatically upon disconnection.

#### socket.to(room)
```go
socket.On("room 237", func(event *socketio.EventPayload) {
 	// to one room
	socket.To("room 237").Emit("test", "hello")

	// to multiple rooms
	socket.To("room 237").To("room 123").Emit("test", "hello")
})
```

# Example

Please check more examples into folder in project for details. [Examples](https://github.com/doquangtan/gofiber-socket.io/tree/main/example)

```go
package main

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

	io.Of("/test").OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)
		socket.Join("demo")
		io.To("demo").Emit("test", socket.Id+" join us room...", "server message")

		socket.On("connected", func(event *socketio.EventPayload) {
			socket.Emit("chat message", "Main")
		})
		socket.On("test", func(event *socketio.EventPayload) {
			socket.Emit("test", event.Data...)
		})

		socket.On("join-room", func(event *socketio.EventPayload) {
			if len(event.Data) > 0 && event.Data[0] != nil {
				socket.Join(event.Data[0].(string))
			}
		})

		socket.On("to-room", func(event *socketio.EventPayload) {
			socket.To("demo").To("demo2").Emit("test", "hello")
		})

		socket.On("leave-room", func(event *socketio.EventPayload) {
			socket.Leave("demo")
			socket.Join("demo2")
		})

		socket.On("my-room", func(event *socketio.EventPayload) {
			socket.Emit("my-room", socket.Rooms())
		})

		socket.On("chat message", func(event *socketio.EventPayload) {
			socket.Emit("chat message", event.Data[0])

			if event.Callback != nil {
				(*event.Callback)("hello", map[string]interface{}{
					"Test": "ok",
				})
			}
		})

		socket.On("disconnecting", func(event *socketio.EventPayload) {
			println("disconnecting", socket.Nps, socket.Id)
		})

		socket.On("disconnect", func(event *socketio.EventPayload) {
			println("disconnect", socket.Nps, socket.Id)
		})
	})

	io.OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)
	})

	app.Use("/", io.Middleware)
	app.Route("/socket.io", io.Server)
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Static("/", "./public")

	app.Route("/", socketIoRoute)

	app.Get("/test", func(c *fiber.Ctx) error {
		io := c.Locals("io").(*socketio.Io)

		io.Emit("event", map[string]interface{}{
			"Ok": 1,
		})

		io.Of("/admin").Emit("event", map[string]interface{}{
			"Ok": 1,
		})

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}


```
