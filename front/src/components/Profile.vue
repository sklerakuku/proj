<template>
  <div class="profile-page">
    <h2>Profile Settings</h2>
    
    <div class="profile-card">
      <div class="avatar-section">
        <div class="avatar-large">{{ avatarLetter }}</div>
        <h3>{{ formData.username }}</h3>
        <span class="role-badge">{{ formData.role }}</span>
      </div>
      
      <form @submit.prevent="saveProfile" class="profile-form">
        <div class="form-group">
          <label>Username</label>
          <input 
            v-model="formData.username" 
            class="sketch-input" 
            placeholder="New username"
            :class="{ error: errors.username }"
          />
          <span v-if="errors.username" class="error-msg">{{ errors.username }}</span>
        </div>
        
        <div class="form-group">
          <label>New Password (leave empty to keep current)</label>
          <input 
            v-model="formData.password" 
            type="password" 
            class="sketch-input" 
            placeholder="New password"
          />
        </div>
        
        <div class="form-group">
          <label>Theme</label>
          <div class="theme-selector">
            <button 
              type="button"
              :class="{ active: theme === 'light' }"
              @click="setTheme('light')"
            >☀️ Light</button>
            <button 
              type="button"
              :class="{ active: theme === 'dark' }"
              @click="setTheme('dark')"
            >🌙 Dark</button>
            <button 
              type="button"
              :class="{ active: theme === 'system' }"
              @click="setTheme('system')"
            >💻 System</button>
          </div>
        </div>
        
        <button type="submit" class="save-btn" :disabled="saving">
          {{ saving ? 'Saving...' : 'Save Changes' }}
        </button>
        
        <div v-if="successMsg" class="success-msg">{{ successMsg }}</div>
        <div v-if="errorMsg" class="error-banner">{{ errorMsg }}</div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref({})
const saving = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const theme = ref(localStorage.getItem('theme') || 'system')

const formData = reactive({
  username: '',
  password: '',
  role: ''
})

const errors = reactive({
  username: ''
})

const avatarLetter = computed(() => {
  return formData.username ? formData.username.charAt(0).toUpperCase() : '?'
})

const apiCall = async (url, options = {}) => {
  const token = localStorage.getItem('auth_token')
  if (!token) {
    router.push('/login')
    throw new Error('Not authenticated')
  }

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    ...options.headers
  }

  const response = await fetch(url, { ...options, headers })
  
  if (response.status === 401) {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('user')
    router.push('/login')
    throw new Error('Session expired')
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.error || 'Request failed')
  }

  return response.json()
}

const applyTheme = (t) => {
  document.documentElement.classList.remove('light-theme', 'dark-theme')
  if (t === 'light') {
    document.documentElement.classList.add('light-theme')
  } else if (t === 'dark') {
    document.documentElement.classList.add('dark-theme')
  }
  // system - убираем оба класса, работает auto через prefers-color-scheme
}

const setTheme = (t) => {
  theme.value = t
  localStorage.setItem('theme', t)
  applyTheme(t)
}

const validateForm = () => {
  errors.username = ''
  if (!formData.username.trim()) {
    errors.username = 'Username is required'
    return false
  }
  if (formData.username.length < 3) {
    errors.username = 'Username must be at least 3 characters'
    return false
  }
  return true
}

const saveProfile = async () => {
  successMsg.value = ''
  errorMsg.value = ''
  
  if (!validateForm()) return
  
  saving.value = true
  try {
    const body = { username: formData.username.trim(), role: formData.role }
    if (formData.password.trim()) {
      body.password = formData.password
    }
    
    await apiCall(`/admin/users/${user.value.id}`, {
      method: 'PUT',
      body: JSON.stringify(body)
    })
    
    // Обновляем localStorage
    const updatedUser = { ...user.value, username: formData.username.trim() }
    localStorage.setItem('user', JSON.stringify(updatedUser))
    user.value = updatedUser
    
    successMsg.value = 'Profile updated successfully!'
    formData.password = ''
    
    setTimeout(() => { successMsg.value = '' }, 3000)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  const stored = localStorage.getItem('user')
  if (stored) {
    try {
      user.value = JSON.parse(stored)
      formData.username = user.value.username || ''
      formData.role = user.value.role || ''
    } catch (e) {
      user.value = {}
    }
  }
  applyTheme(theme.value)
})
</script>

<style scoped>
.profile-page {
  width: 100%;
  max-width: 500px;
  padding: 2rem;
}

h2 {
  font-family: var(--font-1);
  margin-bottom: 1.5rem;
  text-align: center;
}

.profile-card {
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sketch);
  padding: 2rem;
  background: var(--color-background);
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 2rem;
}

.avatar-large {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: bold;
  border: 3px solid var(--color-border);
  margin-bottom: 10px;
}

.avatar-section h3 {
  font-family: var(--font-1);
  margin: 5px 0;
}

.role-badge {
  padding: 4px 12px;
  border-radius: var(--radius-md);
  background: var(--color-info-light);
  color: var(--color-info);
  font-size: 0.8rem;
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.sketch-input {
  width: 100%;
  padding: 10px 15px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-1);
  background: transparent;
  color: var(--color-text);
  transition: border-color 0.2s;
}

.sketch-input:focus {
  outline: none;
  border-color: var(--color-primary);
}

.sketch-input.error {
  border-color: var(--color-danger);
}

.error-msg {
  color: var(--color-danger);
  font-size: 0.75rem;
  margin-top: 3px;
}

.theme-selector {
  display: flex;
  gap: 8px;
}

.theme-selector button {
  flex: 1;
  padding: 8px 12px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-btn);
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
  transition: all 0.2s;
}

.theme-selector button:hover {
  background: var(--color-hover);
  transform: translateY(-1px);
}

.theme-selector button.active {
  background: var(--color-text);
  color: var(--color-background);
}

.save-btn {
  padding: 12px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-btn);
  background: var(--color-primary);
  color: white;
  cursor: pointer;
  font-family: var(--font-1);
  font-size: 1rem;
  transition: all 0.2s;
  margin-top: 0.5rem;
}

.save-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.success-msg {
  padding: 10px;
  background: var(--color-success-light);
  color: var(--color-success);
  border-radius: var(--radius-sm);
  text-align: center;
  font-size: 0.9rem;
}

.error-banner {
  padding: 10px;
  background: var(--color-danger-light);
  color: var(--color-danger);
  border-radius: var(--radius-sm);
  text-align: center;
  font-size: 0.9rem;
}
</style>