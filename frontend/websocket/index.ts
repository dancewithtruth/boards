class WebSocketConnection {
  private socket: WebSocket;

  constructor(url: string) {
    this.socket = new WebSocket(url);

    this.socket.onopen = () => {
      console.log('WebSocket connection established');
    };

    this.socket.onmessage = (event) => {
      console.log('Received message:', event.data);
    };

    this.socket.onclose = () => {
      console.log('WebSocket connection closed');
    };
  }

  send(message: string) {
    this.socket.send(message);
  }

  close() {
    this.socket.close();
  }
}

export default WebSocketConnection;
