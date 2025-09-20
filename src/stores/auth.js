import { writable } from 'svelte/store';
import axios from 'axios';

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
    },
    
    logout() {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      delete axios.defaults.headers.common['Authorization'];
      
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