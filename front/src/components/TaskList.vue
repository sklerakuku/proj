<template>
  <div class="dashboard">
    <div class="processes-area">
      <div class="processes-header">
        <h2>Active Processes ({{ activeProcesses.length }})</h2>
        <div class="header-buttons">
          <button class="action-btn" @click="showTemplateModal = true">From Template</button>
          <button class="action-btn plus" @click="showCreateModal = true">+</button>
        </div>
      </div>

      <div v-if="loading" class="status-message">Loading processes...</div>
      <div v-else-if="error" class="status-message error">{{ error }}</div>
      
      <div v-else class="processes-grid">
        <ProcessCard 
          v-for="process in activeProcesses" 
          :key="process.id"
          :process="process"
          @click="openProcess(process.id)"
        />
        <div v-if="activeProcesses.length === 0" class="status-message">
          No active processes. Create one to get started!
        </div>
      </div>

      <!-- Completed -->
      <div v-if="completedProcesses.length > 0" class="completed-section">
        <div class="completed-header" @click="showCompleted = !showCompleted">
          <h3>Completed ({{ completedProcesses.length }}) <span>{{ showCompleted ? '▼' : '▶' }}</span></h3>
        </div>
        <div v-if="showCompleted" class="completed-grid">
          <ProcessCard 
            v-for="process in completedProcesses" 
            :key="process.id"
            :process="process"
            :completed="true"
            @click="openProcess(process.id)"
          />
        </div>
      </div>

      <!-- History (archived) -->
      <div v-if="archivedProcesses.length > 0" class="completed-section">
        <div class="completed-header" @click="showArchived = !showArchived">
          <h3>History ({{ archivedProcesses.length }}) <span>{{ showArchived ? '▼' : '▶' }}</span></h3>
        </div>
        <div v-if="showArchived" class="completed-grid">
          <ProcessCard 
            v-for="process in archivedProcesses" 
            :key="process.id"
            :process="process"
            :completed="true"
            @click="openProcess(process.id)"
          />
        </div>
      </div>
    </div>

    <!-- Модалка создания процесса без шаблона -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <h3>Create New Process</h3>
        <form @submit.prevent="createEmptyProcess">
          <input v-model="newProcessTitle" placeholder="Process title" class="sketch-input" required />
          <textarea v-model="newProcessDescription" placeholder="Description (optional)" class="sketch-input" rows="3"></textarea>
          <div class="modal-buttons">
            <button type="submit" class="btn-save" :disabled="creatingEmpty">Create</button>
            <button type="button" class="btn-cancel" @click="showCreateModal = false">Cancel</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Модалка выбора шаблона -->
    <div v-if="showTemplateModal" class="modal-overlay" @click.self="showTemplateModal = false">
      <div class="modal">
        <h3>Select Template</h3>
        <div v-if="templates.length === 0" class="status-message">No templates available</div>
        <div v-else class="template-list">
          <div v-for="tpl in templates" :key="tpl.id" class="template-item" @click="createFromTemplate(tpl)">
            <strong>{{ tpl.name }}</strong>
            <p>{{ tpl.description }}</p>
            <small>{{ tpl.tasks?.length || 0 }} tasks</small>
          </div>
        </div>
        <button class="btn-cancel" @click="showTemplateModal = false">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ProcessCard from './common/ProcessCard.vue'

const router = useRouter()
const processes = ref([])
const templates = ref([])
const loading = ref(true)
const error = ref(null)
const showCompleted = ref(false)
const showArchived = ref(false)

const showCreateModal = ref(false)
const showTemplateModal = ref(false)
const newProcessTitle = ref('')
const newProcessDescription = ref('')
const creatingEmpty = ref(false)

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

const activeProcesses = computed(() => 
  processes.value.filter(p => p.status !== 'done' && p.status !== 'archived')
)
const completedProcesses = computed(() => 
  processes.value.filter(p => p.status === 'done')
)
const archivedProcesses = computed(() => 
  processes.value.filter(p => p.status === 'archived')
)

const fetchProcesses = async () => {
  loading.value = true
  error.value = null
  try {
    processes.value = await apiCall('/processes') || []
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const fetchTemplates = async () => {
  try {
    templates.value = await apiCall('/templates') || []
  } catch (e) {
    console.error('Failed to load templates', e)
  }
}

const createEmptyProcess = async () => {
  creatingEmpty.value = true
  try {
    const response = await apiCall('/processes/empty', {
      method: 'POST',
      body: JSON.stringify({
        title: newProcessTitle.value
      })
    })
    processes.value.unshift(response) // добавляем в начало списка
    showCreateModal.value = false
    newProcessTitle.value = ''
    newProcessDescription.value = ''
  } catch (err) {
    alert(`Failed to create process: ${err.message}`)
  } finally {
    creatingEmpty.value = false
  }
}

const createFromTemplate = async (tpl) => {
  try {
    await apiCall('/processes', {
      method: 'POST',
      body: JSON.stringify({
        template_id: tpl.id,
        title: `${tpl.name} - ${new Date().toLocaleDateString()}`
      })
    })
    showTemplateModal.value = false
    await fetchProcesses()
  } catch (err) {
    alert(`Failed to create process: ${err.message}`)
  }
}

const openProcess = (id) => {
  router.push(`/processes/${id}`)
}

onMounted(() => {
  fetchProcesses()
  fetchTemplates()
})
</script>

<style scoped>
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
}

.processes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.processes-header h2 { font-family: var(--font-1); }

.header-buttons { display: flex; gap: 10px; }

.action-btn {
  padding: 8px 20px;
  border: 2px solid var(--color-text);
  border-radius: 255px 150px 225px 150px/150px 225px 150px 255px;
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
}
.action-btn.plus { width: 42px; padding: 8px; }
.action-btn:hover { transform: translateY(-1px); }

.processes-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.completed-section {
  margin-top: 2rem;
  border-top: 2px dashed var(--color-text);
  padding-top: 1rem;
}

.completed-header {
  cursor: pointer;
  padding: 10px;
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  font-family: var(--font-1);
}
.completed-header:hover { background: rgba(0, 0, 0, 0.03); }

.completed-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
}

/* Модалки */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--color-background);
  border: 3px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  padding: 2rem;
  min-width: 400px;
  max-width: 90vw;
}

.modal h3 { font-family: var(--font-1); margin-bottom: 1rem; }

.sketch-input {
  width: 100%;
  padding: 10px;
  margin-bottom: 10px;
  border: 2px solid var(--color-text);
  border-radius: 15px;
  font-family: var(--font-1);
  background: transparent;
  color: var(--color-text);
}

.modal-buttons { display: flex; gap: 10px; margin-top: 1rem; }

.btn-save, .btn-cancel {
  padding: 8px 20px;
  border: 2px solid var(--color-text);
  border-radius: 255px 150px 225px 150px/150px 225px 150px 255px;
  cursor: pointer;
  font-family: var(--font-1);
  background: transparent;
  color: var(--color-text);
}
.btn-save { border-color: var(--color-success); color: var(--color-success); }
.btn-cancel { border-color: #999; color: #999; }

.template-list { max-height: 300px; overflow-y: auto; margin-bottom: 1rem; }

.template-item {
  padding: 12px;
  border: 2px solid var(--color-text);
  border-radius: 15px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.1s ease;
}
.template-item:hover { transform: translateY(-2px); }
.template-item p { margin: 4px 0; font-size: 0.85rem; }
.template-item small { opacity: 0.6; }

.status-message {
  text-align: center;
  padding: 2rem;
  font-family: var(--font-1);
  border: 2px dashed var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
}
.status-message.error { border-color: var(--color-danger); color: var(--color-danger); }
</style>