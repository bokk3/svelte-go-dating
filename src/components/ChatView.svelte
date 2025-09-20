<script>
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  // Placeholder component for chat functionality
  function goBack() {
    dispatch('back');
  }
</script>

<div class="chat-container">
  <div class="chat-header">
    <button class="back-btn" on:click={goBack}>
      ← Back
    </button>
    <h2>Chat</h2>
  </div>
  
  <div class="chat-content">
    <p>Chat functionality coming soon...</p>
  </div>
</div>

<style>
  .chat-container {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  
  .chat-header {
    background: white;
    padding: 16px;
    border-radius: var(--border-radius);
    margin-bottom: 16px;
    display: flex;
    align-items: center;
    gap: 16px;
  }
  
  .back-btn {
    background: none;
    border: none;
    font-size: 16px;
    cursor: pointer;
    color: var(--primary-color);
  }
  
  .chat-header h2 {
    margin: 0;
    color: var(--dark-color);
  }
  
  .chat-content {
    flex: 1;
    background: white;
    border-radius: var(--border-radius);
    padding: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>