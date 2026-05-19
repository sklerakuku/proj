
<template>
  <div>
    <header>
      <h1 @click="goHome" class="logo">{{ t('app.title') }}</h1>
      
      <div class="language-switcher">
        <button @click="setLocale('ru')" :class="{ active: locale === 'ru' }">RU</button>
        <button @click="setLocale('en')" :class="{ active: locale === 'en' }">EN</button>
      </div>
      
      <div v-if="isAuthenticated" class="user" @click="toggleMenu">
        <h3 class="name">{{ username }}</h3>
        <div class="user-picture">{{ avatarLetter }}</div>
        <div v-if="showMenu" class="dropdown-menu" @click.stop>
          <button @click="goToProfile">{{ t('nav.profile') }}</button>
          <button v-if="isAdmin" @click="goToAdmin">{{ t('nav.admin') }}</button>
          <button @click="logout">{{ t('nav.logout') }}</button>
        </div>
      </div>
    </header>
    
    <main>
      <router-view />
    </main>

    <footer>
      <p>&copy; {{ new Date().getFullYear() }} {{ t('footer.copyright') }}</p>
    </footer>
  </div>
</template>


<script setup>
import { useI18n } from './i18n'
const { t, locale, setLocale } = useI18n()
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const isAuthenticated = ref(false)
const username = ref('')
const showMenu = ref(false)
const isAdmin = ref(false)

const avatarLetter = computed(() => {
  return username.value ? username.value.charAt(0).toUpperCase() : '?'
})

const updateAuthState = () => {
  const token = localStorage.getItem('auth_token')
  const user = localStorage.getItem('user')
  
  if (token && user) {
    isAuthenticated.value = true
    try {
      const userData = JSON.parse(user)
      username.value = userData.username || 'User'
      isAdmin.value = userData.role === 'admin'
    } catch (e) {
      username.value = 'User'
    }
  } else {
    isAuthenticated.value = false
    username.value = ''
    isAdmin.value = false
  }
}

watch(() => route.path, () => {
  updateAuthState()
})

const toggleMenu = () => {
  showMenu.value = !showMenu.value
}

const goHome = () => {
  if (isAuthenticated.value) {
    router.push('/dashboard')
  }
}

const goToProfile = () => {
  showMenu.value = false
  router.push('/profile')
}

const goToAdmin = () => {
  showMenu.value = false
  router.push('/admin')
}

const logout = () => {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('user')
  isAuthenticated.value = false
  isAdmin.value = false
  showMenu.value = false
  router.push('/login')
}

const handleClickOutside = (event) => {
  if (!event.target.closest('.user')) {
    showMenu.value = false
  }
}

onMounted(() => {
  updateAuthState()
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('storage', updateAuthState)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('storage', updateAuthState)
})
</script>

<style scoped>
header {
  position: sticky;
  top: 0;
  padding: 1rem 5%;
  background: var(--color-background);
  color: var(--color-text);
  border-bottom: 2px solid var(--color-border);
  font-family: var(--font-1);
  font-size: 1.5rem;
  border-radius: var(--radius-sketch);
  backdrop-filter: blur(10px);
  z-index: 100;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo { margin: 0; font-size: inherit; cursor: pointer; }
.logo:hover { opacity: 0.8; }

.user { display: flex; align-items: center; gap: 12px; cursor: pointer; position: relative; }

.name { font-size: 1rem; margin: 0; }

.user-picture {
  width: 40px; height: 40px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.2rem;
  border: 2px solid var(--color-border);
}

.dropdown-menu {
  position: absolute;
  top: 50px; right: 0;
  background: var(--color-background);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 8px 0;
  min-width: 160px;
  z-index: 200;
}

.dropdown-menu button {
  display: block;
  width: 100%;
  padding: 8px 16px;
  background: none;
  border: none;
  text-align: left;
  font-family: var(--font-1);
  cursor: pointer;
  color: var(--color-text);
  transition: background 0.15s;
}

.dropdown-menu button:hover { background: var(--color-hover); }

main {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 140px);
  overflow-x: hidden;
}

footer {
  font-family: var(--font-1);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: 0.75rem;
  padding: 1rem;
}
</style>