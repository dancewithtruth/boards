import { useEffect, useRef, useState, useCallback } from 'react';

export const useWebSocket = (url: string) => {
  const [messages, setMessages] = useState<any[]>([]);
  const [error, setError] = useState<Event | null>(null);
  const [readyState, setReadyState] = useState<WebSocket['readyState']>();
  const ws = useRef<WebSocket>();

  const send = useCallback((data: string) => {
    if (ws.current) {
      ws.current.send(data);
    }
  }, []);

  useEffect(() => {
    ws.current = new WebSocket(url);

    ws.current.onopen = () => {
      console.log(`Connected: ${url}`);
      setReadyState(ws.current?.readyState);
    };

    ws.current.onmessage = (event) => {
      const receivedData = event.data;
      const messages = receivedData.split('\n');
      const parsedMessages: any[] = [];

      messages.forEach((message: any) => {
        try {
          const parsedMessage = JSON.parse(message);
          console.log(parsedMessage);
          // Process each individual message
          parsedMessages.push(parsedMessage);
        } catch (error) {
          console.log('Error parsing message:', error);
        }
      });
      setMessages(parsedMessages);
    };

    ws.current.onerror = (event) => {
      setError(event);
    };

    ws.current.onclose = (event) => {
      setReadyState(ws.current?.readyState);
    };

    return () => {
      ws.current?.close();
    };
  }, [url]);

  return { send, messages, error, readyState };
};
