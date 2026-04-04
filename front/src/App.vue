<template>
  <div>
    <header v-if="isAuthenticated">
      <h1>Tasks Workflow Manager</h1>
      <div class="user" @click="toggleMenu">
        <h3 class="name">{{ username }}</h3>
        <div class="user-picture">{{ avatarLetter }}</div>
        <div v-if="showMenu" class="dropdown-menu">
          <button @click="goToProfile">Profile</button>
          <button @click="logout">Logout</button>
        </div>
      </div>
    </header>
    
    <main>
      <router-view></router-view>
    </main>

    <footer v-if="isAuthenticated">
      <p> &copy; {{ new Date().getFullYear() }} =^-^=</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isAuthenticated = ref(false)
const username = ref('')
const showMenu = ref(false)

const avatarLetter = computed(() => {
  return username.value ? username.value.charAt(0).toUpperCase() : '?'
})

const toggleMenu = () => {
  showMenu.value = !showMenu.value
}

const goToProfile = () => {
  showMenu.value = false
  router.push('/profile')
}

const logout = () => {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('user')
  isAuthenticated.value = false
  router.push('/login')
}

// Закрыть меню при клике вне
const handleClickOutside = (event) => {
  if (!event.target.closest('.user')) {
    showMenu.value = false
  }
}

onMounted(() => {
  const token = localStorage.getItem('auth_token')
  const user = localStorage.getItem('user')
  
  if (token && user) {
    isAuthenticated.value = true
    try {
      const userData = JSON.parse(user)
      username.value = userData.username || 'User'
    } catch (e) {
      username.value = 'User'
    }
  }
  
  document.addEventListener('click', handleClickOutside)
})
</script>

<style scoped>
header {
  position: sticky;
  top: 0;
  padding: 1rem 15%;
  background: var(--color-background);
  color: var(--color-text);
  border-bottom: 2px solid var(--color-text);
  transform: rotate(-0.5deg);
  font-family: var(--font-1);
  font-size: 1.5rem;
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  backdrop-filter: blur(10px);
  z-index: 100;
  display: flex;
  justify-content: space-between;
}

header h1 {
  margin: 0;
  font-size: inherit;
}

.user {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  position: relative;
}

.name {
  font-size: 1rem;
  margin: 0;
}

.user-picture {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #4a90e2;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.2rem;
  border: 2px solid var(--color-text);
}

.dropdown-menu {
  position: absolute;
  top: 50px;
  right: 0;
  background: var(--color-background);
  border: 2px solid var(--color-text);
  border-radius: 15px;
  padding: 8px 0;
  min-width: 120px;
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
}

.dropdown-menu button:hover {
  background: rgba(0, 0, 0, 0.05);
}

.nav-link {
  color: var(--color-text);
  text-decoration: none;
  font-size: 0.9rem;
  padding: 4px 8px;
  border-radius: 15px;
  transition: all 0.1s ease;
}

.nav-link:hover {
  background: rgba(0, 0, 0, 0.05);
  transform: translateY(-1px);
}

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
  margin: 0;
  color: var(--color-text);
  opacity: 0.7;
  font-size: 0.75rem;
  transform: rotate(0.5deg);
}

@media (max-width: 768px) {
  header {
    padding: 1rem 5%;
    font-size: 1.2rem;
  }
  
  main {
    padding: 1rem;
  }
}
</style>