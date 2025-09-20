<script>
  import { createEventDispatcher } from 'svelte';
  import { authStore } from '../stores/auth.js';
  
  const dispatch = createEventDispatcher();
  
  let email = '';
  let password = '';
  
  $: loading = $authStore.loading;
  $: error = $authStore.error;
  
  async function handleLogin() {
    if (!email || !password) {
      return;
    }
    
    const success = await authStore.login(email, password);
    if (success) {
      dispatch('success');
    }
  }
  
  function switchToRegister() {
    authStore.clearError();
    dispatch('switchToRegister');
  }
</script>

<div class="login-container">
  <div class="login-card">
    <h1>Welcome Back</h1>
    <p>Sign in to find your perfect match</p>
    
    <form on:submit|preventDefault={handleLogin}>
      <div class="form-group">
        <label for="email" class="label">Email</label>
        <input
          type="email"
          id="email"
          bind:value={email}
          class="input"
          placeholder="Enter your email"
          required
        />
      </div>
      
      <div class="form-group">
        <label for="password" class="label">Password</label>
        <input
          type="password"
          id="password"
          bind:value={password}
          class="input"
          placeholder="Enter your password"
          required
        />
      </div>
      
      {#if error}
        <div class="text-error">{error}</div>
      {/if}
      
      <button type="submit" class="btn btn-primary full-width" disabled={loading}>
        {#if loading}
          <span class="loading"></span>
          Signing in...
        {:else}
          Sign In
        {/if}
      </button>
    </form>
    
    <div class="auth-switch">
      <p>Don't have an account? 
        <button type="button" class="link-btn" on:click={switchToRegister}>
          Sign up
        </button>
      </p>
    </div>
  </div>
</div>

<style>
  .login-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 20px;
  }
  
  .login-card {
    background: white;
    border-radius: var(--border-radius);
    box-shadow: var(--shadow-lg);
    padding: 40px;
    width: 100%;
    max-width: 400px;
  }
  
  h1 {
    margin: 0 0 8px 0;
    font-size: 28px;
    font-weight: 700;
    text-align: center;
    background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  p {
    margin: 0 0 32px 0;
    text-align: center;
    color: #6b7280;
  }
  
  .full-width {
    width: 100%;
    margin-top: 8px;
  }
  
  .auth-switch {
    margin-top: 24px;
    text-align: center;
  }
  
  .auth-switch p {
    margin: 0;
    color: #6b7280;
  }
  
  .link-btn {
    background: none;
    border: none;
    color: var(--primary-color);
    cursor: pointer;
    font-weight: 500;
    text-decoration: underline;
  }
  
  .link-btn:hover {
    color: var(--secondary-color);
  }
</style>