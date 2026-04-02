<template>
  <div class="login">
    <h2>Log In</h2>
    
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <input
          id="username"
          class="sketch-input"
          type="text"
          v-model="formData.username"
          placeholder="username"
          :class="{ 'error': errors.username }"
          @blur="validateField('username')"
          required
        />
        <span v-if="errors.username" class="error-message">{{ errors.username }}</span>
      </div>

      <div class="form-group">
        <input
          id="password"
          class="sketch-input"
          type="password"
          v-model="formData.password"
          placeholder="password"
          :class="{ 'error': errors.password }"
          @blur="validateField('password')"
          required
        />
        <span v-if="errors.password" class="error-message">{{ errors.password }}</span>
      </div>

      <button 
        type="submit" 
        class="sketch-button"
        :disabled="isLoading"
      >
        {{ isLoading ? 'Logging in...' : 'Enter' }}
      </button>
      
      <div class="login-footer">
        <router-link to="/signup" class="register-link">Don't have an account?</router-link>
      </div>
    </form>

    <div v-if="errorMessage" class="error-banner">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Form state
const formData = reactive({
  username: '',
  password: ''
})

// UI state
const isLoading = ref(false)
const errorMessage = ref('')
const errors = reactive({
  username: '',
  password: ''
})

// Validation rules
const validateField = (field) => {
  switch (field) {
    case 'username':
      if (!formData.username.trim()) {
        errors.username = 'Username is required'
      } else if (formData.username.length < 3) {
        errors.username = 'Username must be at least 3 characters'
      } else {
        errors.username = ''
      }
      break
    
    case 'password':
      if (!formData.password) {
        errors.password = 'Password is required'
      } else if (formData.password.length < 6) {
        errors.password = 'Password must be at least 6 characters'
      } else {
        errors.password = ''
      }
      break
  }
}

const validateForm = () => {
  validateField('username')
  validateField('password')
  return !errors.username && !errors.password
}

const handleLogin = async () => {
  // Clear previous error
  errorMessage.value = ''
  
  // Validate form
  if (!validateForm()) {
    return
  }

  isLoading.value = true

  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: formData.username.trim(),
        password: formData.password
      })
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || 'Login failed')
    }

    // Store token
    if (data.token) {
      localStorage.setItem('auth_token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
    }

    // Redirect to dashboard after successful login
    router.push('/dashboard')
    
  } catch (error) {
    errorMessage.value = error.message || 'An error occurred during login'
    console.error('Login error:', error)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Your existing styles */
.login {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--color-text);
  font-family: var(--font-1);
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
}

.login h2 {
  margin: 0 0 1rem 0;
  font-size: 2rem;
}

form {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5em;
}

.form-group {
  margin-bottom: 1rem;
  width: 100%;
}

.sketch-input {
  width: 100%;
  padding: 12px 20px;
  border: 3px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  font-size: 1.3rem;
  font-family: var(--font-1);
  background: transparent;
  color: var(--color-text);
  transition: all 0.1s ease;
  box-sizing: border-box;
}

.sketch-input:focus {
  outline: none;
  border-color: var(--color-primary, #4a90e2);
  transform: scale(1.02);
}

.sketch-input.error {
  border-color: #ff4444;
}

.sketch-input::placeholder {
  color: var(--color-text);
  opacity: 50%;
}

.error-message {
  display: block;
  color: #ff4444;
  font-size: 0.75rem;
  margin-top: 0.25rem;
  padding-left: 0.5rem;
}

.sketch-button {
  width: 40%;
  padding: 12px;
  border: 3px solid var(--color-text);
  border-radius: 255px 150px 225px 150px/150px 225px 150px 255px;
  color: var(--color-text);
  font-size: 1.125rem;
  font-family: var(--font-1);
  background: transparent;
  cursor: pointer;
  transition: all 0.1s ease;
  margin-top: 0.5rem;
}

.sketch-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  background: rgba(255, 255, 255, 0.1);
}

.sketch-button:active:not(:disabled) {
  transform: translateY(0);
}

.sketch-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  margin-top: 1rem;
}

.register-link {
  color: var(--color-text);
  text-decoration: underline;
  font-size: 0.875rem;
  transition: opacity 0.1s cubic-bezier();
  cursor: pointer;
}

.register-link:hover {
  text-decoration: underline;
  opacity: 0.8;
}

.error-banner {
  margin-top: 1rem;
  padding: 0.75rem;
  background: rgba(255, 68, 68, 0.1);
  border: 1px solid #ff4444;
  border-radius: 8px;
  color: #ff4444;
  text-align: center;
  font-size: 0.875rem;
  width: 100%;
}

@media (max-width: 768px) {
  .login {
    padding: 1.5rem;
    max-width: 90%;
  }
  
  .sketch-input,
  .sketch-button {
    font-size: 0.875rem;
  }
}
</style>