import { writable } from 'svelte/store';
import axios from 'axios';

// Import WebSocket store (avoid circular dependency by importing lazily)
let wsStore;

// Create auth store
function createAuthStore() {
  const { subscribe, set, update } = writable({
    token: null,
    user: null,
    loading: false,
    error: null
  });

  return {
    subscribe,
    
    async login(email, password) {
      update(store => ({ ...store, loading: true, error: null }));
      
      try {
        const response = await axios.post('/api/v1/login', {
          email,
          password
        });
        
        const { user, tokens } = response.data;
        
        // Store tokens
        localStorage.setItem('access_token', tokens.access_token);
        localStorage.setItem('refresh_token', tokens.refresh_token);
        
        // Update store
        update(store => ({
          ...store,
          token: tokens.access_token,
          user,
          loading: false
        }));
        
        // Set up axios default headers
        axios.defaults.headers.common['Authorization'] = `Bearer ${tokens.access_token}`;
        
        // Connect WebSocket
        if (!wsStore) {
          const { wsStore: ws } = await import('./websocket.js');
          wsStore = ws;
        }
        wsStore.connect(tokens.access_token);
        
        return true;
      } catch (error) {
        const errorMessage = error.response?.data?.error || 'Login failed';
        update(store => ({
          ...store,
          loading: false,
          error: errorMessage
        }));
        return false;
      }
    },
    
    async register(userData) {
      update(store => ({ ...store, loading: true, error: null }));
      
      try {
        const response = await axios.post('/api/v1/register', userData);
        
        const { user, tokens } = response.data;
        
        // Store tokens
        localStorage.setItem('access_token', tokens.access_token);
        localStorage.setItem('refresh_token', tokens.refresh_token);
        
        // Update store
        update(store => ({
          ...store,
          token: tokens.access_token,
          user,
          loading: false
        }));
        
        // Set up axios default headers
        axios.defaults.headers.common['Authorization'] = `Bearer ${tokens.access_token}`;
        
        // Connect WebSocket
        if (!wsStore) {
          const { wsStore: ws } = await import('./websocket.js');
          wsStore = ws;
        }
        wsStore.connect(tokens.access_token);
        
        return true;
      } catch (error) {
        const errorMessage = error.response?.data?.error || 'Registration failed';
        update(store => ({
          ...store,
          loading: false,
          error: errorMessage
        }));
        return false;
      }
    },
    
    setToken(token) {
      update(store => ({ ...store, token }));
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      
      // Verify token and get user info
      this.verifyToken();
    },
    
    async verifyToken() {
      try {
        const response = await axios.get('/api/v1/me');
        update(store => ({ ...store, user: response.data.user }));
        return true;
      } catch (error) {
        console.error('Token verification failed:', error);
        this.logout();
        return false;
      }
    },
    
    logout() {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      delete axios.defaults.headers.common['Authorization'];
      
      // Disconnect WebSocket
      if (wsStore) {
        wsStore.disconnect();
      }
      
      set({
        token: null,
        user: null,
        loading: false,
        error: null
      });
    },
    
    clearError() {
      update(store => ({ ...store, error: null }));
    }
  };
}

export const authStore = createAuthStore();