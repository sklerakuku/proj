<template>
  <div class="process-view">
    <!-- Верхняя панель инструментов -->
    <div class="process-toolbar">
      <button class="tool-btn" @click="toggleSidebar">
        <span class="icon">☰</span>
      </button>
      <button class="tool-btn" @click="addNewProcess">
        <span class="icon">+</span>
      </button>
      <button class="tool-btn" @click="openSettings">
        <span class="icon">Settings</span>
      </button>
      <button class="tool-btn" @click="showStats">
        <span class="icon">Stats</span>
      </button>
      <div class="process-name">{{ process.name }}</div>
      <div class="process-timer">⏱️ {{ formattedTimer }}</div>
    </div>

    <!-- Граф задач -->
    <div class="graph-container" ref="graphContainer">
      <svg 
        class="graph-svg" 
        :viewBox="`${viewBox.x} ${viewBox.y} ${viewBox.width} ${viewBox.height}`"
        @wheel.prevent="handleZoom"
      >
        <defs>
          <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="10" refY="3.5" orient="auto">
            <polygon points="0 0, 10 3.5, 0 7" fill="var(--color-text)" />
          </marker>
        </defs>
        
        <!-- Линии -->
        <line 
          v-for="edge in edges" 
          :key="edge.id"
          :x1="edge.x1" 
          :y1="edge.y1" 
          :x2="edge.x2" 
          :y2="edge.y2"
          stroke="var(--color-text)"
          stroke-width="2"
          marker-end="url(#arrowhead)"
        />
        
        <!-- Задачи (ноды) -->
        <g 
          v-for="task in tasks" 
          :key="task.id"
          :transform="`translate(${task.x}, ${task.y})`"
          @click="selectTask(task)"
        >
          <foreignObject width="200" height="100">
            <TaskNode 
              :task="task"
              :selected="selectedTask?.id === task.id"
              @complete="completeTask"
            />
          </foreignObject>
        </g>
      </svg>
      
      <!-- Контролы зума -->
      <div class="zoom-controls">
        <button @click="zoomIn">+</button>
        <button @click="zoomOut">-</button>
        <button @click="resetZoom">⟳</button>
      </div>
    </div>

    <!-- Боковое меню задачи -->
    <TaskSidebar 
      v-if="selectedTask"
      :task="selectedTask"
      @close="selectedTask = null"
      @update="updateTask"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TaskNode from './common/TaskNode.vue'
import TaskSidebar from './common/TaskSidebar.vue'

const route = useRoute()
const router = useRouter()
const processId = ref(route.params.id)

// Состояние процесса
const process = ref({
  id: 1,
  name: 'Website Redesign',
  startTime: new Date(Date.now() - 3600000 * 24 * 2) // 2 дня назад
})

// Задачи (ноды) с координатами
const tasks = ref([
  { 
    id: 1, 
    title: 'Analysis', 
    description: 'Gather requirements and analyze market',
    status: 'done',
    assignedTo: 'Anna',
    role: 'analyst',
    x: 50, y: 200,
    attachments: ['requirements.pdf']
  },
  { 
    id: 2, 
    title: 'Design', 
    description: 'Create UI/UX designs',
    status: 'done',
    assignedTo: 'Mike',
    role: 'designer',
    x: 300, y: 200,
    attachments: ['design.fig']
  },
  { 
    id: 3, 
    title: 'Development', 
    description: 'Implement features',
    status: 'in_progress',
    assignedTo: 'John',
    role: 'developer',
    x: 550, y: 200,
    attachments: []
  },
  { 
    id: 4, 
    title: 'Testing', 
    description: 'QA and bug fixes',
    status: 'pending',
    assignedTo: 'Sarah',
    role: 'tester',
    x: 800, y: 200,
    attachments: []
  },
  { 
    id: 5, 
    title: 'Deployment', 
    description: 'Deploy to production',
    status: 'pending',
    assignedTo: 'Ops',
    role: 'devops',
    x: 1050, y: 200,
    attachments: []
  }
])

// Связи между задачами
const edges = computed(() => [
  { id: 1, x1: 250, y1: 250, x2: 300, y2: 250 },
  { id: 2, x1: 500, y1: 250, x2: 550, y2: 250 },
  { id: 3, x1: 750, y1: 250, x2: 800, y2: 250 },
  { id: 4, x1: 1000, y1: 250, x2: 1050, y2: 250 }
])

const selectedTask = ref(null)
const showSidebarFlag = ref(true)

// Zoom и панорамирование
const graphContainer = ref(null)
const viewBox = ref({ x: 0, y: 0, width: 1200, height: 500 })
const scale = ref(1)

// Таймер
const timer = ref(0)
let timerInterval = null

const formattedTimer = computed(() => {
  const hours = Math.floor(timer.value / 3600)
  const minutes = Math.floor((timer.value % 3600) / 60)
  const seconds = timer.value % 60
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})

const toggleSidebar = () => {
  showSidebarFlag.value = !showSidebarFlag.value
  // Можно добавить логику скрытия боковой панели
}

const addNewProcess = () => {
  router.push('/editor')
}

const openSettings = () => {
  alert('Process settings')
}

const showStats = () => {
  alert(`Process stats: ${tasks.value.length} tasks, ${tasks.value.filter(t => t.status === 'done').length} completed`)
}

const selectTask = (task) => {
  selectedTask.value = task
}

const completeTask = (taskId) => {
  const task = tasks.value.find(t => t.id === taskId)
  if (task && task.status !== 'done') {
    task.status = 'done'
    selectedTask.value = null
  }
}

const updateTask = (updatedTask) => {
  const index = tasks.value.findIndex(t => t.id === updatedTask.id)
  if (index !== -1) {
    tasks.value[index] = updatedTask
  }
  selectedTask.value = null
}

// Zoom функции
const zoomIn = () => {
  scale.value *= 1.2
  updateViewBox()
}

const zoomOut = () => {
  scale.value *= 0.8
  updateViewBox()
}

const resetZoom = () => {
  scale.value = 1
  updateViewBox()
}

const handleZoom = (event) => {
  if (event.deltaY < 0) {
    zoomIn()
  } else {
    zoomOut()
  }
}

const updateViewBox = () => {
  viewBox.value.width = 1200 / scale.value
  viewBox.value.height = 500 / scale.value
}

// Таймер
const startTimer = () => {
  if (process.value.startTime) {
    const start = new Date(process.value.startTime).getTime()
    timerInterval = setInterval(() => {
      timer.value = Math.floor((Date.now() - start) / 1000)
    }, 1000)
  }
}

onMounted(() => {
  startTimer()
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})
</script>

<style scoped>
.process-view {
  width: 100%;
  height: calc(100vh - 140px);
  display: flex;
  flex-direction: column;
  background: var(--color-background);
  position: relative;
}

.process-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-bottom: 2px solid var(--color-text);
  background: var(--color-background);
}

.tool-btn {
  padding: 6px 12px;
  background: transparent;
  border: 2px solid var(--color-text);
  border-radius: 255px 150px 225px 150px/150px 225px 150px 255px;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
}

.tool-btn .icon {
  font-size: 1.1rem;
}

.process-name {
  flex: 1;
  font-weight: bold;
  font-family: var(--font-1);
  margin-left: 20px;
}

.process-timer {
  font-family: monospace;
  font-size: 1.2rem;
  background: rgba(0, 0, 0, 0.05);
  padding: 4px 12px;
  border-radius: 20px;
}

.graph-container {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: var(--color-background);
}

.graph-svg {
  width: 100%;
  height: 100%;
  cursor: grab;
}

.graph-svg:active {
  cursor: grabbing;
}

.zoom-controls {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  gap: 8px;
}

.zoom-controls button {
  width: 36px;
  height: 36px;
  border: 2px solid var(--color-text);
  border-radius: 250px 150px 220px 150px/150px 225px 150px 255px;
  background: var(--color-background);
  color: var(--color-text);
  cursor: pointer;
  font-size: 1.2rem;
}

@media (max-width: 768px) {
  .process-toolbar {
    flex-wrap: wrap;
  }
  
  .process-name {
    width: 100%;
    margin-left: 0;
    order: -1;
  }
}

.graph-svg foreignObject {
  overflow: visible;
}

.graph-svg {
  overflow: visible;
}
</style>