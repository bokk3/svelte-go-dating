<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import { authStore } from '../stores/auth.js';
  import { wsStore } from '../stores/websocket.js';
  import axios from 'axios';
  
  export let matchId;
  
  const dispatch = createEventDispatcher();
  
  let messages = [];
  let newMessage = '';
  let messageInput;
  let messagesContainer;
  let match = null;
  let otherUser = null;
  let loading = true;
  let error = null;
  
  $: currentUser = $authStore.user;
  $: wsMessages = $wsStore.messages[matchId] || [];
  $: isConnected = $wsStore.connected;
  $: isTyping = $wsStore.typing[matchId];
  
  // Combine loaded messages with real-time messages
  $: allMessages = [...messages, ...wsMessages].sort((a, b) => 
    new Date(a.created_at) - new Date(b.created_at)
  );
  
  onMount(async () => {
    await loadMatchData();
    await loadMessages();
    scrollToBottom();
    
    // Connect WebSocket if not connected
    if (!isConnected && $authStore.token) {
      wsStore.connect($authStore.token);
    }
  });
  
  onDestroy(() => {
    // Optional: disconnect WebSocket when leaving chat
    // wsStore.disconnect();
  });
  
  async function loadMatchData() {
    try {
      const response = await axios.get(`/api/v1/matches/${matchId}`);
      match = response.data.match;
      otherUser = response.data.other_user;
    } catch (err) {
      error = 'Failed to load match data';
      console.error('Error loading match:', err);
    }
  }
  
  async function loadMessages() {
    try {
      const response = await axios.get(`/api/v1/matches/${matchId}/messages`);
      messages = response.data.messages || [];
      loading = false;
      setTimeout(scrollToBottom, 100);
    } catch (err) {
      error = 'Failed to load messages';
      loading = false;
      console.error('Error loading messages:', err);
    }
  }
  
  function sendMessage() {
    if (!newMessage.trim() || !isConnected) return;
    
    wsStore.sendMessage(matchId, newMessage.trim());
    newMessage = '';
    setTimeout(scrollToBottom, 100);
  }
  
  function handleKeyPress(event) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      sendMessage();
    } else {
      // Send typing indicator
      wsStore.sendTyping(matchId);
    }
  }
  
  function scrollToBottom() {
    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }
  
  function goBack() {
    dispatch('back');
  }
  
  function formatTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = (now - date) / (1000 * 60 * 60);
    
    if (diffInHours < 24) {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } else {
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
  }
</script>

<div class="chat-container">
  <div class="chat-header">
    <button class="back-btn" on:click={goBack}>
      ← Back
    </button>
    {#if otherUser}
      <div class="user-info">
        <h2>{otherUser.display_name}</h2>
        {#if isConnected}
          <span class="connection-status online">● Online</span>
        {:else}
          <span class="connection-status offline">● Offline</span>
        {/if}
      </div>
    {:else}
      <h2>Loading...</h2>
    {/if}
  </div>
  
  {#if loading}
    <div class="loading">Loading messages...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else}
    <div class="messages-container" bind:this={messagesContainer}>
      {#if allMessages.length === 0}
        <div class="no-messages">
          <p>Start your conversation!</p>
          <p>Say hello to {otherUser?.display_name || 'your match'}</p>
        </div>
      {:else}
        {#each allMessages as message}
          <div class="message {message.sender_id === currentUser?.id ? 'sent' : 'received'}">
            <div class="message-content">
              {message.message}
            </div>
            <div class="message-time">
              {formatTime(message.created_at)}
            </div>
          </div>
        {/each}
      {/if}
      
      {#if isTyping && isTyping !== currentUser?.id}
        <div class="typing-indicator">
          <div class="typing-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <span class="typing-text">{otherUser?.display_name} is typing...</span>
        </div>
      {/if}
    </div>
    
    <div class="message-input-container">
      <div class="input-wrapper">
        <textarea
          bind:this={messageInput}
          bind:value={newMessage}
          on:keypress={handleKeyPress}
          placeholder="Type a message..."
          rows="1"
          disabled={!isConnected}
        ></textarea>
        <button 
          class="send-btn" 
          on:click={sendMessage}
          disabled={!newMessage.trim() || !isConnected}
        >
          Send
        </button>
      </div>
      {#if !isConnected}
        <div class="connection-warning">
          Connecting to chat...
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .chat-container {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--bg-color, #f5f5f5);
  }
  
  .chat-header {
    background: white;
    padding: 16px;
    border-bottom: 1px solid #e0e0e0;
    display: flex;
    align-items: center;
    gap: 16px;
    flex-shrink: 0;
  }
  
  .back-btn {
    background: none;
    border: none;
    font-size: 16px;
    cursor: pointer;
    color: var(--primary-color, #ff4458);
    padding: 4px;
  }
  
  .user-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  
  .chat-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }
  
  .connection-status {
    font-size: 12px;
    font-weight: 500;
  }
  
  .connection-status.online {
    color: #4caf50;
  }
  
  .connection-status.offline {
    color: #999;
  }
  
  .loading, .error {
    display: flex;
    justify-content: center;
    align-items: center;
    flex: 1;
    color: #666;
  }
  
  .error {
    color: #f44336;
  }
  
  .messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .no-messages {
    text-align: center;
    color: #666;
    margin-top: 50%;
    transform: translateY(-50%);
  }
  
  .no-messages p {
    margin: 8px 0;
  }
  
  .message {
    display: flex;
    flex-direction: column;
    max-width: 70%;
    word-wrap: break-word;
  }
  
  .message.sent {
    align-self: flex-end;
    align-items: flex-end;
  }
  
  .message.received {
    align-self: flex-start;
    align-items: flex-start;
  }
  
  .message-content {
    background: white;
    padding: 12px 16px;
    border-radius: 18px;
    margin-bottom: 4px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }
  
  .message.sent .message-content {
    background: var(--primary-color, #ff4458);
    color: white;
  }
  
  .message-time {
    font-size: 11px;
    color: #999;
    margin: 0 8px;
  }
  
  .typing-indicator {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #666;
    font-size: 14px;
    margin-left: 16px;
  }
  
  .typing-dots {
    display: flex;
    gap: 2px;
  }
  
  .typing-dots span {
    width: 6px;
    height: 6px;
    background: #999;
    border-radius: 50%;
    animation: typing 1.4s infinite;
  }
  
  .typing-dots span:nth-child(2) {
    animation-delay: 0.2s;
  }
  
  .typing-dots span:nth-child(3) {
    animation-delay: 0.4s;
  }
  
  @keyframes typing {
    0%, 60%, 100% {
      transform: translateY(0);
    }
    30% {
      transform: translateY(-10px);
    }
  }
  
  .message-input-container {
    background: white;
    padding: 16px;
    border-top: 1px solid #e0e0e0;
    flex-shrink: 0;
  }
  
  .input-wrapper {
    display: flex;
    gap: 12px;
    align-items: flex-end;
  }
  
  .input-wrapper textarea {
    flex: 1;
    border: 1px solid #ddd;
    border-radius: 20px;
    padding: 12px 16px;
    resize: none;
    max-height: 100px;
    font-family: inherit;
    font-size: 14px;
    outline: none;
  }
  
  .input-wrapper textarea:focus {
    border-color: var(--primary-color, #ff4458);
  }
  
  .input-wrapper textarea:disabled {
    background: #f5f5f5;
    color: #999;
  }
  
  .send-btn {
    background: var(--primary-color, #ff4458);
    color: white;
    border: none;
    border-radius: 20px;
    padding: 12px 20px;
    cursor: pointer;
    font-weight: 500;
    transition: opacity 0.2s;
  }
  
  .send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .send-btn:not(:disabled):hover {
    background: var(--primary-color-dark, #e63946);
  }
  
  .connection-warning {
    text-align: center;
    font-size: 12px;
    color: #f44336;
    margin-top: 8px;
  }
  
  /* Mobile responsiveness */
  @media (max-width: 768px) {
    .message {
      max-width: 85%;
    }
    
    .chat-header {
      padding: 12px;
    }
    
    .messages-container {
      padding: 12px;
    }
    
    .message-input-container {
      padding: 12px;
    }
  }
</style>