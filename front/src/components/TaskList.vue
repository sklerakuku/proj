<template>
  <div class="dashboard">
    <div class="main-content">
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
    
      <!-- Sidebar -->
      <Sidebar 
        @search="handleSearch" 
        @filter-change="handleFilterChange"
        @date-change="handleDateChange"
      />
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
import Sidebar from './common/Sidebar.vue'

const router = useRouter()
const processes = ref([])
const templates = ref([])
const loading = ref(true)
const error = ref(null)
const showCompleted = ref(false)
const showArchived = ref(false)

const searchQuery = ref('')
const priorityFilter = ref([])
const dateFilter = ref(null)

const showCreateModal = ref(false)
const showTemplateModal = ref(false)
const newProcessTitle = ref('')
const newProcessDescription = ref('')
const creatingEmpty = ref(false)

const apiCall = async (url, options = {}) => {
  const token = localStorage.getItem('auth_token')
  if (!token) { router.push('/login'); throw new Error('Not authenticated') }

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

const filterProcesses = (list) => {
  let result = [...list]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(p => p.title?.toLowerCase().includes(q))
  }
  if (dateFilter.value) {
    const filterDate = new Date(dateFilter.value).toDateString()
    result = result.filter(p => {
      if (!p.started_at) return false
      return new Date(p.started_at).toDateString() === filterDate
    })
  }
  return result
}

// Единые computed с фильтрацией
const activeProcesses = computed(() => 
  filterProcesses(processes.value.filter(p => p.status !== 'done' && p.status !== 'archived'))
)
const completedProcesses = computed(() => 
  filterProcesses(processes.value.filter(p => p.status === 'done'))
)
const archivedProcesses = computed(() => 
  filterProcesses(processes.value.filter(p => p.status === 'archived'))
)

const handleSearch = (query) => { searchQuery.value = query }
const handleFilterChange = (filter) => { priorityFilter.value = filter.priority || [] }
const handleDateChange = (date) => { dateFilter.value = date }

const fetchProcesses = async () => {
  loading.value = true
  error.value = null
  try { processes.value = await apiCall('/processes') || [] }
  catch (err) { error.value = err.message }
  finally { loading.value = false }
}

const fetchTemplates = async () => {
  try { templates.value = await apiCall('/templates') || [] }
  catch (e) { console.error('Failed to load templates', e) }
}

const createEmptyProcess = async () => {
  creatingEmpty.value = true
  try {
    const response = await apiCall('/processes/empty', {
      method: 'POST',
      body: JSON.stringify({ title: newProcessTitle.value })
    })
    processes.value.unshift(response)
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

const openProcess = (id) => { router.push(`/processes/${id}`) }

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
}

.main-content {
  display: flex;
  padding: 2rem 5%;
  gap: 2rem;
  align-items: flex-start;
  max-width: 1400px;
  margin: 0 auto;
}

.processes-area {
  flex: 3;
  max-height: calc(100vh - 150px);
  overflow-y: auto;
  padding-right: 8px;
}

.processes-area::-webkit-scrollbar {
  width: 8px;
}
.processes-area::-webkit-scrollbar-track {
  background: transparent;
}
.processes-area::-webkit-scrollbar-thumb {
  background: var(--color-muted);
  border-radius: 4px;
}

.processes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  position: sticky;
  top: 0;
  background: var(--color-background);
  z-index: 5;
  padding: 10px 0;
}
.processes-header h2 { font-family: var(--font-1); }

.header-buttons { display: flex; gap: 10px; }

.action-btn {
  padding: 8px 20px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-btn);
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
  border-top: 2px dashed var(--color-muted);
  padding-top: 1rem;
}

.completed-header {
  cursor: pointer;
  padding: 10px;
  border-radius: var(--radius-sketch);
  font-family: var(--font-1);
}
.completed-header:hover { background: var(--color-hover); }

.completed-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
}

.status-message {
  text-align: center;
  padding: 2rem;
  font-family: var(--font-1);
  border: 2px dashed var(--color-border);
  border-radius: var(--radius-sketch);
}
.status-message.error { border-color: var(--color-danger); color: var(--color-danger); }

@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
    padding: 1rem;
  }
  .processes-area {
    max-height: none;
    overflow: visible;
  }
}

/* Модалка */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: var(--color-overlay);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: var(--color-background);
  border: 3px solid var(--color-border);
  border-radius: var(--radius-sketch);
  padding: 2rem;
  min-width: 400px;
}
.modal h3 { font-family: var(--font-1); margin-bottom: 1rem; }

.modal-buttons { display: flex; gap: 10px; margin-top: 1rem; }

.template-list {
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 16px;
}

.template-item {
  padding: 12px;
  border: 1px solid var(--color-border, #ddd);
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
}

.template-item:hover {
  background: var(--color-hover, #f5f5f5);
}

.btn-save, .btn-cancel {
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  border: none;
}

.btn-save {
  background: var(--color-primary, #4caf50);
  color: white;
}

.btn-cancel {
  background: var(--color-muted, #ccc);
}
</style>