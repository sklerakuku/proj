<template>
  <div class="signin">
    <h2>Sign In</h2>
    
    <form @submit.prevent="handleSignup">
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

      <div class="form-group">
        <input
          id="rep-password"
          class="sketch-input"
          type="password"
          v-model="formData.confirmPassword"
          placeholder="repeat password"
          :class="{ 'error': errors.confirmPassword }"
          @blur="validateField('confirmPassword')"
          required
        />
        <span v-if="errors.confirmPassword" class="error-message">{{ errors.confirmPassword }}</span>
      </div>

      <!-- Compact CAPTCHA - click on canvas to refresh -->
      <div class="captcha-line">
        <canvas 
          ref="captchaCanvas" 
          class="sketch-captcha"
          width="150" 
          height="50"
          @click="refreshCaptcha"
          title="Click to refresh"
        ></canvas>
        <input
          class="sketch-input captcha-input"
          type="text"
          v-model="formData.captchaCode"
          placeholder="captcha code"
          :class="{ 'error': errors.captchaCode }"
          @blur="validateField('captchaCode')"
          maxlength="6"
          required
        />
      </div>
      <span v-if="errors.captchaCode" class="error-message">{{ errors.captchaCode }}</span>

      <button 
        type="submit" 
        class="sketch-button"
        :disabled="isLoading"
      >
        {{ isLoading ? 'Creating account...' : 'Create account' }}
      </button>
      
      <div class="signin-footer">
        <router-link to="/login" class="register-link">Already have an account?</router-link>
      </div>
    </form>

    <div v-if="errorMessage" class="error-banner">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Form state
const formData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  captchaCode: ''
})

// UI state
const isLoading = ref(false)
const errorMessage = ref('')
const captchaCanvas = ref(null)
const currentCaptcha = ref('')

const errors = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  captchaCode: ''
})

// Generate random CAPTCHA
const generateCaptcha = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ0123456789'
  let captcha = ''
  for (let i = 0; i < 6; i++) {
    captcha += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  currentCaptcha.value = captcha
  
  // Draw on canvas
  if (captchaCanvas.value) {
    const ctx = captchaCanvas.value.getContext('2d')
    const canvas = captchaCanvas.value
    
    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    
    // Background
    ctx.fillStyle = '#f5f5f5'
    ctx.fillRect(0, 0, canvas.width, canvas.height)
    
    // Add noise dots
    for (let i = 0; i < 100; i++) {
      ctx.fillStyle = `rgba(0,0,0,${Math.random() * 0.3})`
      ctx.fillRect(Math.random() * canvas.width, Math.random() * canvas.height, 1, 1)
    }
    
    // Draw CAPTCHA text with distortion
    for (let i = 0; i < captcha.length; i++) {
      ctx.save()
      ctx.translate(20 + i * 18, 30)
      ctx.rotate((Math.random() - 0.5) * 0.4)
      ctx.font = `bold ${22 + Math.random() * 4}px monospace`
      ctx.fillStyle = `#${Math.floor(Math.random() * 16777215).toString(16)}`
      ctx.fillText(captcha[i], 0, 0)
      ctx.restore()
    }
    
    // Add random lines
    for (let i = 0; i < 4; i++) {
      ctx.beginPath()
      ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height)
      ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height)
      ctx.strokeStyle = `rgba(0,0,0,${Math.random() * 0.5})`
      ctx.lineWidth = 1.5
      ctx.stroke()
    }
  }
}

const refreshCaptcha = () => {
  generateCaptcha()
  formData.captchaCode = ''
  errors.captchaCode = ''
}

// Validation rules
const validateField = (field) => {
  switch (field) {
    case 'username':
      if (!formData.username.trim()) {
        errors.username = 'Username is required'
      } else if (formData.username.length < 3) {
        errors.username = 'Username must be at least 3 characters'
      } else if (formData.username.length > 20) {
        errors.username = 'Username must be less than 20 characters'
      } else if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
        errors.username = 'Username can only contain letters, numbers, and underscores'
      } else {
        errors.username = ''
      }
      break
    
    case 'password':
      if (!formData.password) {
        errors.password = 'Password is required'
      } else if (formData.password.length < 6) {
        errors.password = 'Password must be at least 6 characters'
      } else if (!/(?=.*[A-Z])(?=.*[0-9])/.test(formData.password)) {
        errors.password = 'Password must contain at least one uppercase letter and one number'
      } else {
        errors.password = ''
      }
      break
    
    case 'confirmPassword':
      if (!formData.confirmPassword) {
        errors.confirmPassword = 'Please confirm your password'
      } else if (formData.confirmPassword !== formData.password) {
        errors.confirmPassword = 'Passwords do not match'
      } else {
        errors.confirmPassword = ''
      }
      break
    
    case 'captchaCode':
      if (!formData.captchaCode.trim()) {
        errors.captchaCode = 'CAPTCHA code is required'
      } else if (formData.captchaCode.toUpperCase() !== currentCaptcha.value) {
        errors.captchaCode = 'Invalid CAPTCHA code'
      } else {
        errors.captchaCode = ''
      }
      break
  }
}

const validateForm = () => {
  validateField('username')
  validateField('password')
  validateField('confirmPassword')
  validateField('captchaCode')
  
  return !errors.username && 
         !errors.password && 
         !errors.confirmPassword && 
         !errors.captchaCode
}

// Handle signup submission
const handleSignup = async () => {
  // Clear previous error
  errorMessage.value = ''
  
  // Validate form
  if (!validateForm()) {
    return
  }

  isLoading.value = true

  try {
    // Prepare user data
    const userData = {
      username: formData.username.trim(),
      password: formData.password,
      email: `${formData.username.trim()}@example.com`
    }

    // Make API request to create account
    const response = await fetch('/api/auth/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData)
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || 'Account creation failed')
    }

    // Store token if received
    if (data.token) {
      localStorage.setItem('auth_token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
    }

    console.log('Account created successfully:', data)
    
    // Redirect to login
    setTimeout(() => {
      router.push('/login')
    }, 1500)
    
  } catch (error) {
    errorMessage.value = error.message || 'An error occurred during account creation'
    console.error('Signup error:', error)
    
    // Refresh CAPTCHA on error
    refreshCaptcha()
  } finally {
    isLoading.value = false
  }
}

// Initialize CAPTCHA on component mount
onMounted(() => {
  generateCaptcha()
})
</script>

<style scoped>
.signin {
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

.signin h2 {
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

/* CAPTCHA Styles */
.captcha-line {
  display: flex;
  gap: 12px;
  align-items: center;
  width: 100%;
  margin-bottom: 0.25rem;
}

.sketch-captcha {
  border: 3px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  cursor: pointer;
  background: #f5f5f5;
  flex-shrink: 0;
  transition: all 0.1s ease;
}

.sketch-captcha:hover {
  transform: scale(1.02);
  border-color: var(--color-primary, #4a90e2);
}

.captcha-input {
  flex: 1;
  min-width: 0;
  padding: 12px 15px;
  font-size: 1.1rem;
}

.error-message {
  display: block;
  color: #ff4444;
  font-size: 0.75rem;
  margin-top: 0.25rem;
  padding-left: 0.5rem;
}

.sketch-button {
  width: 100%;
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

.signin-footer {
  text-align: center;
  margin-top: 1rem;
}

.register-link {
  color: var(--color-text);
  text-decoration: underline;
  font-size: 0.875rem;
  transition: opacity 0.1s cubic-bezier();
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
  .signin {
    padding: 1.5rem;
    max-width: 90%;
  }
  
  .sketch-input,
  .sketch-button {
    font-size: 0.875rem;
  }
  
  .captcha-line {
    gap: 8px;
  }
  
  .sketch-captcha {
    width: 120px;
    height: 42px;
  }
  
  .captcha-input {
    font-size: 0.875rem;
  }
}
</style>