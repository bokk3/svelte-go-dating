<script>
  import { createEventDispatcher } from 'svelte';
  import { authStore } from '../stores/auth.js';
  
  const dispatch = createEventDispatcher();
  
  let email = '';
  let password = '';
  let confirmPassword = '';
  let displayName = '';
  let age = '';
  let gender = '';
  let gdprConsent = false;
  
  $: loading = $authStore.loading;
  $: error = $authStore.error;
  $: passwordsMatch = password === confirmPassword;
  $: formValid = email && password && confirmPassword && displayName && age && gender && gdprConsent && passwordsMatch;
  
  async function handleRegister() {
    if (!formValid) {
      return;
    }
    
    const userData = {
      email,
      password,
      display_name: displayName,
      age: parseInt(age),
      gender,
      gdpr_consent: gdprConsent
    };
    
    const success = await authStore.register(userData);
    if (success) {
      dispatch('success');
    }
  }
  
  function switchToLogin() {
    authStore.clearError();
    dispatch('switchToLogin');
  }
</script>

<div class="register-container">
  <div class="register-card">
    <h1>Join Today</h1>
    <p>Create your account and start meeting amazing people</p>
    
    <form on:submit|preventDefault={handleRegister}>
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
        <label for="displayName" class="label">Display Name</label>
        <input
          type="text"
          id="displayName"
          bind:value={displayName}
          class="input"
          placeholder="How should we call you?"
          required
        />
      </div>
      
      <div class="form-row">
        <div class="form-group">
          <label for="age" class="label">Age</label>
          <input
            type="number"
            id="age"
            bind:value={age}
            class="input"
            placeholder="Age"
            min="18"
            max="100"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="gender" class="label">Gender</label>
          <select id="gender" bind:value={gender} class="input" required>
            <option value="">Select...</option>
            <option value="man">Man</option>
            <option value="woman">Woman</option>
            <option value="non-binary">Non-binary</option>
            <option value="other">Other</option>
          </select>
        </div>
      </div>
      
      <div class="form-group">
        <label for="password" class="label">Password</label>
        <input
          type="password"
          id="password"
          bind:value={password}
          class="input"
          placeholder="Create a password"
          required
        />
      </div>
      
      <div class="form-group">
        <label for="confirmPassword" class="label">Confirm Password</label>
        <input
          type="password"
          id="confirmPassword"
          bind:value={confirmPassword}
          class="input"
          placeholder="Confirm your password"
          class:error={confirmPassword && !passwordsMatch}
          required
        />
        {#if confirmPassword && !passwordsMatch}
          <div class="text-error">Passwords don't match</div>
        {/if}
      </div>
      
      <div class="form-group">
        <label class="checkbox-label">
          <input
            type="checkbox"
            bind:checked={gdprConsent}
            required
          />
          <span class="checkmark"></span>
          I agree to the <a href="/privacy" target="_blank">Privacy Policy</a> and <a href="/terms" target="_blank">Terms of Service</a>
        </label>
      </div>
      
      {#if error}
        <div class="text-error">{error}</div>
      {/if}
      
      <button type="submit" class="btn btn-primary full-width" disabled={loading || !formValid}>
        {#if loading}
          <span class="loading"></span>
          Creating account...
        {:else}
          Create Account
        {/if}
      </button>
    </form>
    
    <div class="auth-switch">
      <p>Already have an account? 
        <button type="button" class="link-btn" on:click={switchToLogin}>
          Sign in
        </button>
      </p>
    </div>
  </div>
</div>

<style>
  .register-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 20px;
  }
  
  .register-card {
    background: white;
    border-radius: var(--border-radius);
    box-shadow: var(--shadow-lg);
    padding: 40px;
    width: 100%;
    max-width: 450px;
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
  
  .form-row {
    display: flex;
    gap: 16px;
  }
  
  .form-row .form-group {
    flex: 1;
  }
  
  .full-width {
    width: 100%;
    margin-top: 8px;
  }
  
  .input.error {
    border-color: var(--danger-color);
  }
  
  .checkbox-label {
    display: flex;
    align-items: flex-start;
    cursor: pointer;
    font-size: 14px;
    line-height: 1.4;
    position: relative;
    padding-left: 30px;
  }
  
  .checkbox-label input[type="checkbox"] {
    position: absolute;
    opacity: 0;
    cursor: pointer;
  }
  
  .checkmark {
    position: absolute;
    top: 2px;
    left: 0;
    height: 20px;
    width: 20px;
    background-color: #f3f4f6;
    border: 2px solid #d1d5db;
    border-radius: 4px;
    transition: all 0.2s ease;
  }
  
  .checkbox-label input:checked ~ .checkmark {
    background-color: var(--primary-color);
    border-color: var(--primary-color);
  }
  
  .checkmark:after {
    content: "";
    position: absolute;
    display: none;
    left: 6px;
    top: 2px;
    width: 6px;
    height: 10px;
    border: solid white;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
  }
  
  .checkbox-label input:checked ~ .checkmark:after {
    display: block;
  }
  
  .checkbox-label a {
    color: var(--primary-color);
    text-decoration: none;
  }
  
  .checkbox-label a:hover {
    text-decoration: underline;
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
  
  @media (max-width: 768px) {
    .form-row {
      flex-direction: column;
      gap: 0;
    }
    
    .register-card {
      padding: 30px 24px;
    }
  }
</style>