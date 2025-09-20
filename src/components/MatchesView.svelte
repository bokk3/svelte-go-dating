<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import axios from 'axios';
  
  const dispatch = createEventDispatcher();
  
  let matches = [];
  let loading = true;
  let error = null;
  
  onMount(async () => {
    await loadMatches();
  });
  
  async function loadMatches() {
    try {
      loading = true;
      const response = await axios.get('/api/v1/matches');
      matches = response.data;
    } catch (err) {
      error = 'Failed to load matches';
      console.error('Error loading matches:', err);
    } finally {
      loading = false;
    }
  }
  
  function openChat(match) {
    dispatch('openChat', { match });
  }
  
  function formatTime(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) return 'Just now';
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    return date.toLocaleDateString();
  }
</script>

<div class="matches-container">
  <h1>Your Matches</h1>
  
  {#if loading}
    <div class="loading-container">
      <div class="loading"></div>
      <p>Loading your matches...</p>
    </div>
  {:else if error}
    <div class="error-container">
      <p>{error}</p>
      <button class="btn btn-primary" on:click={loadMatches}>
        Try Again
      </button>
    </div>
  {:else if matches.length === 0}
    <div class="empty-container">
      <div class="empty-icon">💫</div>
      <h2>No matches yet</h2>
      <p>Keep swiping to find your perfect match!</p>
    </div>
  {:else}
    <div class="matches-list">
      {#each matches as match}
        <div class="match-card" on:click={() => openChat(match)}>
          <div class="match-avatar">
            <img 
              src={match.user1?.avatar_url || match.user2?.avatar_url || '/placeholder-avatar.jpg'} 
              alt="Match avatar"
            />
            <div class="online-indicator" class:online={true}></div>
          </div>
          
          <div class="match-info">
            <h3>{match.user1?.display_name || match.user2?.display_name}</h3>
            <p class="last-message">
              {#if match.last_message}
                {match.last_message}
              {:else}
                Say hello! 👋
              {/if}
            </p>
          </div>
          
          <div class="match-meta">
            <span class="match-time">{formatTime(match.matched_at)}</span>
            {#if match.unread_count > 0}
              <div class="unread-badge">{match.unread_count}</div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .matches-container {
    height: 100%;
    overflow-y: auto;
  }
  
  h1 {
    color: white;
    text-align: center;
    margin: 0 0 24px 0;
    font-weight: 700;
  }
  
  .loading-container,
  .error-container,
  .empty-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 60vh;
    text-align: center;
  }
  
  .loading-container p,
  .error-container p {
    color: white;
    margin: 20px 0;
  }
  
  .empty-container {
    color: white;
  }
  
  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
  }
  
  .empty-container h2 {
    margin: 0 0 8px 0;
    font-weight: 600;
  }
  
  .empty-container p {
    margin: 0;
    opacity: 0.8;
  }
  
  .matches-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .match-card {
    background: white;
    border-radius: var(--border-radius);
    padding: 16px;
    display: flex;
    align-items: center;
    gap: 16px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: var(--shadow);
  }
  
  .match-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg);
  }
  
  .match-avatar {
    position: relative;
    flex-shrink: 0;
  }
  
  .match-avatar img {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    object-fit: cover;
  }
  
  .online-indicator {
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: #94a3b8;
    border: 2px solid white;
  }
  
  .online-indicator.online {
    background: var(--success-color);
  }
  
  .match-info {
    flex: 1;
    min-width: 0;
  }
  
  .match-info h3 {
    margin: 0 0 4px 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--dark-color);
  }
  
  .last-message {
    margin: 0;
    color: #6b7280;
    font-size: 14px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  
  .match-meta {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
  }
  
  .match-time {
    font-size: 12px;
    color: #9ca3af;
  }
  
  .unread-badge {
    background: var(--primary-color);
    color: white;
    border-radius: 50%;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: 600;
  }
  
  @media (max-width: 768px) {
    .match-card {
      padding: 12px;
      gap: 12px;
    }
    
    .match-avatar img {
      width: 50px;
      height: 50px;
    }
    
    .match-info h3 {
      font-size: 16px;
    }
  }
</style>