<template>
  <div class="dashboard">
    <div class="main-content">

      <!-- Основной контент (80%, со скроллом) -->
      <div class="processes-area">
        <div class="processes-header">
          <h2>Active Tasks ({{ activeProcesses.length }})</h2>
          <div class="header-buttons">
            <button class="action-btn" @click="createFromTemplate">From Template</button>
            <button class="action-btn plus" @click="createEmpty">+</button>
          </div>
        </div>

        <!-- Плитки процессов -->
        <div class="processes-grid">
          <ProcessCard 
            v-for="process in filteredProcesses" 
            :key="process.id"
            :process="process"
            @click="openProcess(process.id)"
          />
        </div>

        <!-- Завершённые процессы (аккордеон) -->
        <div class="completed-section">
          <div class="completed-header" @click="showCompleted = !showCompleted">
            <h3>History ({{ completedProcesses.length }}) <span>{{ showCompleted ? '▼' : '▶' }}</span></h3>
            
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
      </div>
      
      <!-- Боковая панель (фиксированная, 20%) -->
      <div class="sidebar-fixed">
        <SearchFilter 
          @search="handleSearch" 
          @filter-change="handleFilterChange"
          @date-change="handleDateChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ProcessCard from './common/ProcessCard.vue'
import SearchFilter from './common/Sidebar.vue'

const router = useRouter()
const showCompleted = ref(false)
const searchQuery = ref('')
const filters = ref({ priority: [] })
const selectedDate = ref(null)

// Мок данные процессов
const processes = ref([
  {
    id: 1,
    name: 'Website Redesign',
    status: 'active',
    priority: 'high',
    currentStage: 'Development',
    completedTasks: 2,
    totalTasks: 5,
    waitingFor: 'Design review',
    timer: '2d 4h',
    stages: ['Analysis', 'Design', 'Development', 'Testing', 'Deploy'],
    activeStage: 2
  },
  {
    id: 2,
    name: 'Mobile App Launch',
    status: 'active',
    priority: 'medium',
    currentStage: 'Testing',
    completedTasks: 3,
    totalTasks: 4,
    waitingFor: 'QA team',
    timer: '1d 2h',
    stages: ['Planning', 'Development', 'Testing', 'Deploy'],
    activeStage: 2
  },
  {
    id: 3,
    name: 'Database Migration',
    status: 'active',
    priority: 'high',
    currentStage: 'Backup',
    completedTasks: 1,
    totalTasks: 3,
    waitingFor: 'DBA approval',
    timer: '5h',
    stages: ['Backup', 'Migration', 'Verification'],
    activeStage: 0
  },
  {
    id: 4,
    name: 'API Documentation',
    status: 'completed',
    priority: 'low',
    currentStage: 'Done',
    completedTasks: 3,
    totalTasks: 3,
    waitingFor: null,
    timer: '0h'
  }
])

const activeProcesses = computed(() => 
  processes.value.filter(p => p.status === 'active')
)

const completedProcesses = computed(() => 
  processes.value.filter(p => p.status === 'completed')
)

const filteredProcesses = computed(() => {
  let result = [...activeProcesses.value]
  
  if (searchQuery.value) {
    result = result.filter(p => 
      p.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  if (filters.value.priority && filters.value.priority.length > 0) {
    result = result.filter(p => filters.value.priority.includes(p.priority))
  }
  
  return result
})

const createFromTemplate = () => {
  router.push('/editor')
}

const createEmpty = () => {
  alert('Create new process')
}

const openProcess = (id) => {
  router.push(`/process/${id}`)
}

const handleSearch = (query) => {
  searchQuery.value = query
}

const handleFilterChange = (filter) => {
  filters.value = filter
}

const handleDateChange = (date) => {
  selectedDate.value = date
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: var(--color-background);
}

.main-content {
  display: flex;
  padding: 2rem 5% 0 5%;
  gap: 2rem;
  align-items: flex-start;
}

.sidebar-fixed {
  flex: 1;
  position: sticky;
  top: 5%;
  align-self: flex-start;
}

.processes-area {
  width: 70em;
  flex: 4;
  max-height: calc(90vh - 150px);
  overflow-y: auto;
  padding-right: 8px;
}

.processes-area::-webkit-scrollbar {
  width: 8px;
}

.processes-area::-webkit-scrollbar-track {
  background: var(--color-background);
  border-radius: 10px;
}

.processes-area::-webkit-scrollbar-thumb {
  background: var(--color-text);
  border-radius: 10px;
  opacity: 0.5;
}

.processes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.processes-header h2 {
  font-family: var(--font-1);
  margin: 0;
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

.action-btn:hover {
  transform: translateY(-1px);
}

.action-btn.plus {
  width: 42px;
  padding: 8px;
}

.processes-grid {
  display: flex;
  flex-direction: column;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.completed-section {
  margin-top: 2rem;
  border-top: 2px dashed var(--color-text);
  padding-top: 1rem;
}

.completed-header {
  display: flex;
  justify-content: space-between;
  cursor: pointer;
  padding: 10px;
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  transition: all 0.1s ease;
}

.completed-header:hover {
  background: rgba(0, 0, 0, 0.03);
}

.completed-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }
  
  .sidebar-fixed {
    position: static;
    width: 100%;
  }
  
  .processes-grid {
    grid-template-columns: 1fr;
  }
}
</style>