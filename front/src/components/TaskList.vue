<template>
  <div class="dashboard">
    <div class="processes-area">
      <div class="processes-header">
        <h2>Active Processes ({{ processes.length }})</h2>
        <div class="header-buttons">
          <!-- Кнопка для создания нового процесса пока выключена, т.к. нужен template_id -->
          <button class="action-btn" @click="createNewProcess" :disabled="creating">
            {{ creating ? 'Creating...' : 'New Process (Test)' }}
          </button>
        </div>
      </div>

      <div v-if="loading" class="status-message">Loading processes...</div>
      <div v-else-if="error" class="status-message error">{{ error }}</div>
      
      <div v-else class="processes-grid">
        <ProcessCard 
          v-for="process in processes" 
          :key="process.id"
          :process="process"
          @click="openProcess(process.id)"
        />
        <div v-if="processes.length === 0" class="status-message">
          No processes yet. Create one to get started!
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ProcessCard from './common/ProcessCard.vue'

const router = useRouter()
const processes = ref([])
const loading = ref(true)
const error = ref(null)
const creating = ref(false)

// Функция для выполнения API-запросов с авторизацией
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
    throw new Error(errorData.error || `Request failed with status ${response.status}`)
  }

  return response.json()
}

const fetchProcesses = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await apiCall('/processes')
    processes.value = data || []
  } catch (err) {
    error.value = err.message
    console.error(err)
  } finally {
    loading.value = false
  }
}

// Временная функция для создания процесса с ID шаблона 0.
// В реальном приложении здесь должен быть выбор из существующих шаблонов.
const createNewProcess = async () => {
  creating.value = true
  try {
    const newProcess = await apiCall('/processes', {
      method: 'POST',
      body: JSON.stringify({
        template_id: 1, // Здесь должен быть реальный ID существующего шаблона
        title: `New Process ${new Date().toLocaleTimeString()}`
      })
    })
    processes.value.push(newProcess)
  } catch (err) {
    alert(`Failed to create process: ${err.message}`)
    console.error(err)
  } finally {
    creating.value = false
  }
}

const openProcess = (id) => {
  router.push(`/processes/${id}`)
}

onMounted(() => {
  fetchProcesses()
})
</script>

<style scoped>
/* Оставляем ваши оригинальные стили */
.dashboard {
  width: 100%;
  min-height: 100vh;
  background: var(--color-background);
  display: flex;
  justify-content: center;
}

.processes-area {
  width: 100%;
  max-width: 900px;
  padding: 2rem;
  max-height: calc(100vh - 100px);
  overflow-y: auto;
}

.processes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.processes-header h2 {
  font-family: var(--font-1);
}

.header-buttons {
  display: flex;
  gap: 10px;
}

.action-btn {
  padding: 8px 20px;
  border: 2px solid var(--color-text);
  border-radius: 255px 150px 225px 150px/150px 225px 150px 255px;
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
  transition: all 0.1s ease;
}

.action-btn:hover:not(:disabled) {
  transform: translateY(-1px);
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.processes-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.status-message {
  text-align: center;
  padding: 2rem;
  font-family: var(--font-1);
  border: 2px dashed var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
}

.status-message.error {
  border-color: #ff4444;
  color: #ff4444;
}

@media (max-width: 768px) {
  .processes-area {
    padding: 1rem;
  }
}
</style>