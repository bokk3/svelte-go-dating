<script>
  import { onMount } from 'svelte';
  import { authStore } from './stores/auth.js';
  import { wsStore } from './stores/websocket.js';
  
  // Components
  import Login from './components/Login.svelte';
  import Register from './components/Register.svelte';
  import SwipeView from './components/SwipeView.svelte';
  import MatchesView from './components/MatchesView.svelte';
  import ChatView from './components/ChatView.svelte';
  import ProfileView from './components/ProfileView.svelte';
  
  let currentView = 'swipe';
  let showLogin = true;
  
  $: isAuthenticated = $authStore.token !== null;
  
  onMount(() => {
    // Check for stored token
    const token = localStorage.getItem('access_token');
    if (token) {
      authStore.setToken(token);
      wsStore.connect(token);
    }
  });
  
  function handleLogin() {
    showLogin = false;
    currentView = 'swipe';
    wsStore.connect($authStore.token);
  }
  
  function handleLogout() {
    authStore.logout();
    wsStore.disconnect();
    showLogin = true;
  }
  
  function switchView(view) {
    currentView = view;
  }
</script>

<main>
  {#if !isAuthenticated}
    <div class="auth-container">
      {#if showLogin}
        <Login on:success={handleLogin} on:switchToRegister={() => showLogin = false} />
      {:else}
        <Register on:success={handleLogin} on:switchToLogin={() => showLogin = true} />
      {/if}
    </div>
  {:else}
    <div class="app-container">
      <!-- Navigation -->
      <nav class="bottom-nav">
        <button 
          class="nav-btn" 
          class:active={currentView === 'swipe'}
          on:click={() => switchView('swipe')}
        >
          🔥
        </button>
        <button 
          class="nav-btn" 
          class:active={currentView === 'matches'}
          on:click={() => switchView('matches')}
        >
          💬
        </button>
        <button 
          class="nav-btn" 
          class:active={currentView === 'profile'}
          on:click={() => switchView('profile')}
        >
          👤
        </button>
        <button class="nav-btn" on:click={handleLogout}>
          🚪
        </button>
      </nav>
      
      <!-- Main Content -->
      <div class="main-content">
        {#if currentView === 'swipe'}
          <SwipeView />
        {:else if currentView === 'matches'}
          <MatchesView on:openChat={(e) => switchView('chat')} />
        {:else if currentView === 'chat'}
          <ChatView on:back={() => switchView('matches')} />
        {:else if currentView === 'profile'}
          <ProfileView />
        {/if}
      </div>
    </div>
  {/if}
</main>

<style>
  .auth-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 20px;
  }
  
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }
  
  .main-content {
    flex: 1;
    overflow: hidden;
    padding: 20px;
    padding-bottom: 80px; /* Space for bottom nav */
  }
  
  .bottom-nav {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: white;
    display: flex;
    justify-content: space-around;
    padding: 10px 0;
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
    z-index: 100;
  }
  
  .nav-btn {
    background: none;
    border: none;
    font-size: 24px;
    padding: 12px;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.2s ease;
    width: 50px;
    height: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .nav-btn:hover {
    background: rgba(102, 126, 234, 0.1);
  }
  
  .nav-btn.active {
    background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
    color: white;
  }
  
  @media (max-width: 768px) {
    .main-content {
      padding: 10px;
      padding-bottom: 70px;
    }
    
    .nav-btn {
      font-size: 20px;
      width: 45px;
      height: 45px;
    }
  }
</style>