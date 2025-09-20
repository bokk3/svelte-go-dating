import { writable } from 'svelte/store';

// Create WebSocket store
function createWebSocketStore() {
  const { subscribe, set, update } = writable({
    connected: false,
    messages: {},
    onlineUsers: new Set(),
    typing: {}
  });

  let ws = null;

  return {
    subscribe,
    
    connect(token) {
      if (ws) {
        ws.close();
      }
      
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/ws`;
      
      ws = new WebSocket(wsUrl);
      
      ws.onopen = () => {
        console.log('WebSocket connected');
        update(store => ({ ...store, connected: true }));
        
        // Send authentication
        ws.send(JSON.stringify({
          type: 'auth',
          token: token
        }));
      };
      
      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        
        switch (data.type) {
          case 'new_message':
            update(store => {
              const messages = { ...store.messages };
              if (!messages[data.match_id]) {
                messages[data.match_id] = [];
              }
              messages[data.match_id].push(data.data);
              return { ...store, messages };
            });
            break;
            
          case 'user_status':
            update(store => {
              const onlineUsers = new Set(store.onlineUsers);
              if (data.data.status === 'online') {
                onlineUsers.add(data.user_id);
              } else {
                onlineUsers.delete(data.user_id);
              }
              return { ...store, onlineUsers };
            });
            break;
            
          case 'typing':
            update(store => {
              const typing = { ...store.typing };
              typing[data.match_id] = data.user_id;
              
              // Clear typing indicator after 3 seconds
              setTimeout(() => {
                update(s => {
                  const newTyping = { ...s.typing };
                  delete newTyping[data.match_id];
                  return { ...s, typing: newTyping };
                });
              }, 3000);
              
              return { ...store, typing };
            });
            break;
        }
      };
      
      ws.onclose = () => {
        console.log('WebSocket disconnected');
        update(store => ({ ...store, connected: false }));
        
        // Attempt to reconnect after 3 seconds
        setTimeout(() => {
          if (token) {
            this.connect(token);
          }
        }, 3000);
      };
      
      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    },
    
    disconnect() {
      if (ws) {
        ws.close();
        ws = null;
      }
      set({
        connected: false,
        messages: {},
        onlineUsers: new Set(),
        typing: {}
      });
    },
    
    sendMessage(matchId, message) {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: 'send_message',
          match_id: matchId,
          message: message
        }));
      }
    },
    
    sendTyping(matchId) {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: 'typing',
          match_id: matchId
        }));
      }
    }
  };
}

export const wsStore = createWebSocketStore();