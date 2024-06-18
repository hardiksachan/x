# xmessage

`xmessage` is a Go package that provides an asynchronous messaging system. It is part of the larger X project, which provides a collection of libraries for various functionalities.

## Features

- **Message Struct**: `xmessage` introduces a `Message` struct that carries information about the message, including a unique identifier, type, and payload.

- **Publishing Struct**: `xmessage` provides a `Publishing` struct that represents a publishing of an event. It includes the topic of the event and the message of the event.

- **Delivery Interface**: `xmessage` introduces a `Delivery` interface that represents a delivery of an event. It provides methods to acknowledge (`Ack`) or not acknowledge (`Nack`) the delivery.

- **Inbox and Outbox Processing**: `xmessage` provides out-of-the-box support for _transactional outbox_ and _idempotent consumer_. The `inbox` package provides a `Repository` interface for saving and retrieving messages, and the `outbox` package provides a function to create a new message.

---

Happy asy - messag - nchronous - ing
