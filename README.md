# Chat System Backend

This project is a chat system backend implemented in Go, utilizing a microservices architecture with MQTT, REST API, and gRPC for communication. The system supports real-time messaging and is designed for scalability and modularity.

## Setup Instructions

### Prerequisites

- **Go**: Make sure you have Go installed. You can download it from [golang.org](https://golang.org/dl/).
- **Protocol Buffers Compiler (`protoc`)**: Required for generating Go code from `.proto` files. Installation instructions are available [here](https://grpc.io/docs/protoc-installation/).

### Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/chat-system-backend.git
    cd chat-system-backend
    ```

2. **Generate gRPC and Protocol Buffer code**:
    Use the `Makefile` to generate necessary code from `.proto` files.
    ```bash
    make gen
    ```

3. **Run the server**:
    Start the chat system backend server.
    ```bash
    make run
    ```

## Project Structure

- **proto/**: Contains `.proto` files defining the gRPC services and messages.
- **gen/**: Generated Go code from the `.proto` files.
- **cmd/**: Contains the main entry point for the server.

## Contributing

Feel free to open issues or submit pull requests if you find any bugs or have suggestions for improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
