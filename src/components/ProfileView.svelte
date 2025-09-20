<script>
  import { onMount } from 'svelte';
  import axios from 'axios';
  
  let profile = null;
  let loading = true;
  let error = null;
  let editing = false;
  
  onMount(async () => {
    await loadProfile();
  });
  
  async function loadProfile() {
    try {
      loading = true;
      const response = await axios.get('/api/v1/profile');
      profile = response.data;
    } catch (err) {
      error = 'Failed to load profile';
      console.error('Error loading profile:', err);
    } finally {
      loading = false;
    }
  }
  
  async function saveProfile() {
    try {
      await axios.put('/api/v1/profile', profile);
      editing = false;
    } catch (err) {
      error = 'Failed to save profile';
      console.error('Error saving profile:', err);
    }
  }
</script>

<div class="profile-container">
  <h1>Your Profile</h1>
  
  {#if loading}
    <div class="loading-container">
      <div class="loading"></div>
      <p>Loading your profile...</p>
    </div>
  {:else if error}
    <div class="error-container">
      <p>{error}</p>
      <button class="btn btn-primary" on:click={loadProfile}>
        Try Again
      </button>
    </div>
  {:else if profile}
    <div class="profile-content">
      <div class="profile-header">
        <div class="avatar-container">
          <img 
            src={profile.avatar_url || '/placeholder-avatar.jpg'} 
            alt="Your profile"
            class="profile-avatar"
          />
          {#if profile.is_verified}
            <div class="verified-badge">✓</div>
          {/if}
        </div>
        
        <div class="profile-info">
          {#if editing}
            <input
              type="text"
              bind:value={profile.display_name}
              class="input profile-name-input"
              placeholder="Display name"
            />
          {:else}
            <h2>{profile.display_name}</h2>
          {/if}
          
          <p class="profile-stats">
            {#if profile.age}{profile.age} years old{/if}
            {#if profile.location_city} • {profile.location_city}{/if}
          </p>
        </div>
        
        <button 
          class="btn btn-primary" 
          on:click={() => editing ? saveProfile() : editing = true}
        >
          {editing ? 'Save' : 'Edit'}
        </button>
      </div>
      
      <div class="profile-details">
        <div class="detail-section">
          <h3>About</h3>
          {#if editing}
            <textarea
              bind:value={profile.bio}
              class="input bio-input"
              placeholder="Tell others about yourself..."
              rows="4"
            ></textarea>
          {:else}
            <p>{profile.bio || 'No bio yet'}</p>
          {/if}
        </div>
        
        <div class="detail-section">
          <h3>Details</h3>
          <div class="details-grid">
            <div class="detail-item">
              <span class="detail-label">Age:</span>
              {#if editing}
                <input
                  type="number"
                  bind:value={profile.age}
                  class="input detail-input"
                  min="18"
                  max="100"
                />
              {:else}
                <span>{profile.age || 'Not specified'}</span>
              {/if}
            </div>
            
            <div class="detail-item">
              <span class="detail-label">Gender:</span>
              {#if editing}
                <select bind:value={profile.gender} class="input detail-input">
                  <option value="">Select...</option>
                  <option value="man">Man</option>
                  <option value="woman">Woman</option>
                  <option value="non-binary">Non-binary</option>
                  <option value="other">Other</option>
                </select>
              {:else}
                <span>{profile.gender || 'Not specified'}</span>
              {/if}
            </div>
          </div>
        </div>
        
        {#if profile.is_premium}
          <div class="premium-badge">
            <span>👑 Premium Member</span>
          </div>
        {:else}
          <div class="upgrade-section">
            <h3>Upgrade to Premium</h3>
            <p>Get unlimited likes, see who liked you, and more!</p>
            <button class="btn btn-primary">
              Upgrade Now
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .profile-container {
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
  .error-container {
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
  
  .profile-content {
    background: white;
    border-radius: var(--border-radius);
    overflow: hidden;
    box-shadow: var(--shadow-lg);
  }
  
  .profile-header {
    padding: 24px;
    display: flex;
    align-items: center;
    gap: 20px;
    border-bottom: 1px solid #e5e7eb;
  }
  
  .avatar-container {
    position: relative;
    flex-shrink: 0;
  }
  
  .profile-avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
  }
  
  .verified-badge {
    position: absolute;
    bottom: 0;
    right: 0;
    background: var(--success-color);
    color: white;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: bold;
    border: 2px solid white;
  }
  
  .profile-info {
    flex: 1;
  }
  
  .profile-info h2 {
    margin: 0 0 4px 0;
    font-size: 24px;
    font-weight: 600;
  }
  
  .profile-name-input {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
  }
  
  .profile-stats {
    margin: 0;
    color: #6b7280;
    font-size: 14px;
  }
  
  .profile-details {
    padding: 24px;
  }
  
  .detail-section {
    margin-bottom: 32px;
  }
  
  .detail-section h3 {
    margin: 0 0 16px 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--dark-color);
  }
  
  .detail-section p {
    margin: 0;
    color: #4b5563;
    line-height: 1.6;
  }
  
  .bio-input {
    resize: vertical;
    min-height: 100px;
  }
  
  .details-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }
  
  .detail-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .detail-label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
  }
  
  .detail-input {
    font-size: 14px;
  }
  
  .premium-badge {
    background: linear-gradient(135deg, #fbbf24, #f59e0b);
    color: white;
    padding: 16px;
    border-radius: var(--border-radius);
    text-align: center;
    font-weight: 600;
  }
  
  .upgrade-section {
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
    padding: 24px;
    border-radius: var(--border-radius);
    text-align: center;
  }
  
  .upgrade-section h3 {
    margin: 0 0 8px 0;
    color: var(--primary-color);
  }
  
  .upgrade-section p {
    margin: 0 0 16px 0;
    color: #6b7280;
  }
  
  @media (max-width: 768px) {
    .profile-header {
      flex-direction: column;
      text-align: center;
      gap: 16px;
    }
    
    .details-grid {
      grid-template-columns: 1fr;
    }
    
    .profile-details {
      padding: 20px;
    }
  }
</style>