<script>
  import { onMount } from 'svelte';
  import axios from 'axios';
  
  let profiles = [];
  let currentIndex = 0;
  let loading = true;
  let error = null;
  let swipeAnimation = '';
  
  $: currentProfile = profiles[currentIndex];
  
  onMount(async () => {
    await loadProfiles();
  });
  
  async function loadProfiles() {
    try {
      loading = true;
      const response = await axios.get('/api/v1/potential-matches');
      profiles = response.data;
      
      if (profiles.length === 0) {
        error = 'No more profiles to show. Check back later!';
      }
    } catch (err) {
      error = 'Failed to load profiles';
      console.error('Error loading profiles:', err);
    } finally {
      loading = false;
    }
  }
  
  async function swipe(liked) {
    if (!currentProfile) return;
    
    // Animate card
    swipeAnimation = liked ? 'swipe-right' : 'swipe-left';
    
    try {
      const response = await axios.post('/api/v1/swipe', {
        target_user_id: currentProfile.user_id,
        liked: liked
      });
      
      if (response.data.matched) {
        // Show match animation/modal
        showMatchNotification(currentProfile);
      }
      
      // Move to next profile after animation
      setTimeout(() => {
        currentIndex++;
        swipeAnimation = '';
        
        // Load more profiles if running low
        if (currentIndex >= profiles.length - 2) {
          loadProfiles();
        }
      }, 300);
      
    } catch (err) {
      console.error('Error swiping:', err);
      swipeAnimation = '';
    }
  }
  
  function showMatchNotification(profile) {
    // TODO: Implement match notification modal
    alert(`It's a match with ${profile.display_name}! 🎉`);
  }
  
  // Touch handling for mobile swipe gestures
  let startX = 0;
  let startY = 0;
  let deltaX = 0;
  
  function handleTouchStart(e) {
    startX = e.touches[0].clientX;
    startY = e.touches[0].clientY;
  }
  
  function handleTouchMove(e) {
    if (!startX || !startY) return;
    
    deltaX = e.touches[0].clientX - startX;
    const deltaY = e.touches[0].clientY - startY;
    
    // Only handle horizontal swipes
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
      e.preventDefault();
      
      // Apply transform to card
      const card = e.currentTarget;
      const rotation = deltaX * 0.1;
      card.style.transform = `translateX(${deltaX}px) rotate(${rotation}deg)`;
      card.style.opacity = 1 - Math.abs(deltaX) / 300;
    }
  }
  
  function handleTouchEnd(e) {
    const card = e.currentTarget;
    
    // Reset card style
    card.style.transform = '';
    card.style.opacity = '';
    
    // Determine swipe direction
    if (Math.abs(deltaX) > 100) {
      swipe(deltaX > 0);
    }
    
    // Reset values
    startX = 0;
    startY = 0;
    deltaX = 0;
  }
</script>

<div class="swipe-container">
  {#if loading}
    <div class="loading-container">
      <div class="loading"></div>
      <p>Finding amazing people for you...</p>
    </div>
  {:else if error}
    <div class="error-container">
      <p>{error}</p>
      <button class="btn btn-primary" on:click={loadProfiles}>
        Try Again
      </button>
    </div>
  {:else if currentProfile}
    <div class="card-stack">
      <!-- Current card -->
      <div 
        class="profile-card {swipeAnimation}"
        on:touchstart={handleTouchStart}
        on:touchmove={handleTouchMove}
        on:touchend={handleTouchEnd}
      >
        <div class="card-image">
          <img 
            src={currentProfile.avatar_url || '/placeholder-avatar.jpg'} 
            alt={currentProfile.display_name}
            loading="lazy"
          />
          {#if currentProfile.is_verified}
            <div class="verified-badge">✓</div>
          {/if}
        </div>
        
        <div class="card-info">
          <div class="card-header">
            <h2>{currentProfile.display_name}</h2>
            {#if currentProfile.age}
              <span class="age">{currentProfile.age}</span>
            {/if}
          </div>
          
          {#if currentProfile.bio}
            <p class="bio">{currentProfile.bio}</p>
          {/if}
          
          {#if currentProfile.location_city}
            <div class="location">
              📍 {currentProfile.location_city}
              {#if currentProfile.location_country}, {currentProfile.location_country}{/if}
            </div>
          {/if}
        </div>
      </div>
      
      <!-- Preview of next card -->
      {#if profiles[currentIndex + 1]}
        <div class="profile-card preview">
          <div class="card-image">
            <img 
              src={profiles[currentIndex + 1].avatar_url || '/placeholder-avatar.jpg'} 
              alt={profiles[currentIndex + 1].display_name}
              loading="lazy"
            />
          </div>
        </div>
      {/if}
    </div>
    
    <div class="actions">
      <button 
        class="btn btn-circle btn-danger" 
        on:click={() => swipe(false)}
        title="Pass"
      >
        ✕
      </button>
      
      <button 
        class="btn btn-circle btn-success" 
        on:click={() => swipe(true)}
        title="Like"
      >
        ♥
      </button>
    </div>
  {:else}
    <div class="empty-container">
      <h2>No more profiles</h2>
      <p>You've seen everyone in your area! Check back later for new people.</p>
      <button class="btn btn-primary" on:click={loadProfiles}>
        Refresh
      </button>
    </div>
  {/if}
</div>

<style>
  .swipe-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100%;
    max-width: 400px;
    margin: 0 auto;
    position: relative;
  }
  
  .loading-container,
  .error-container,
  .empty-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 70vh;
    text-align: center;
  }
  
  .loading-container p,
  .error-container p,
  .empty-container p {
    color: white;
    margin: 20px 0;
  }
  
  .empty-container h2 {
    color: white;
    margin: 0 0 10px 0;
  }
  
  .card-stack {
    position: relative;
    width: 100%;
    height: 70vh;
    margin-bottom: 20px;
  }
  
  .profile-card {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: white;
    border-radius: var(--border-radius);
    box-shadow: var(--shadow-lg);
    overflow: hidden;
    cursor: grab;
    transition: transform 0.3s ease;
    user-select: none;
  }
  
  .profile-card.preview {
    transform: scale(0.95) translateY(10px);
    z-index: -1;
    opacity: 0.8;
  }
  
  .profile-card:active {
    cursor: grabbing;
  }
  
  .profile-card.swipe-left {
    transform: translateX(-100vw) rotate(-30deg);
    opacity: 0;
  }
  
  .profile-card.swipe-right {
    transform: translateX(100vw) rotate(30deg);
    opacity: 0;
  }
  
  .card-image {
    position: relative;
    height: 70%;
    overflow: hidden;
  }
  
  .card-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .verified-badge {
    position: absolute;
    top: 16px;
    right: 16px;
    background: var(--success-color);
    color: white;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: bold;
  }
  
  .card-info {
    padding: 20px;
    height: 30%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }
  
  .card-header h2 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
  }
  
  .age {
    font-size: 20px;
    color: #6b7280;
    font-weight: 400;
  }
  
  .bio {
    margin: 0 0 12px 0;
    color: #4b5563;
    font-size: 16px;
    line-height: 1.4;
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  
  .location {
    font-size: 14px;
    color: #6b7280;
    display: flex;
    align-items: center;
    gap: 4px;
  }
  
  .actions {
    display: flex;
    justify-content: space-around;
    width: 100%;
    max-width: 200px;
    gap: 40px;
  }
  
  .actions button {
    transform: scale(1);
    transition: transform 0.1s ease;
  }
  
  .actions button:active {
    transform: scale(0.95);
  }
  
  @media (max-width: 768px) {
    .swipe-container {
      padding: 0 10px;
    }
    
    .card-stack {
      height: 65vh;
    }
    
    .card-info {
      padding: 16px;
    }
    
    .actions {
      gap: 60px;
    }
  }
</style>