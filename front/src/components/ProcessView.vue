<template>
  <div class="process-view">
    <div class="process-toolbar">
      <button class="tool-btn" @click="goBack">← Back</button>
      <div class="process-name">{{ process.title }}</div>
      <div class="process-status">Status: {{ process.status }}</div>
    </div>

    <div v-if="loading" class="status-message">Loading...</div>
    <div v-else-if="error" class="status-message error">{{ error }}</div>
    
    <div v-else class="flow-container">
      <VueFlow
        v-model="elements"
        :node-types="nodeTypes"
        :default-viewport="{ zoom: 1 }"
        :min-zoom="0.2"
        :max-zoom="4"
        @node-click="onNodeClick"
      >
        <Background pattern-color="#aaa" :gap="16" />
        <Controls />
        
        <template #node-taskNode="{ data }">
          <TaskNode 
            :task="data.task"
            @complete="completeTask"
            @click="selectTask(data.task)"
          />
        </template>
      </VueFlow>
    </div>

    <TaskSidebar 
      v-if="selectedTask"
      :task="selectedTask"
      @close="selectedTask = null"
      @update="updateTask"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { VueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import TaskNode from './common/TaskNode.vue'
import TaskSidebar from './common/TaskSidebar.vue'

// Стили Vue Flow
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'

const route = useRoute()
const router = useRouter()
const processId = ref(route.params.id)

const process = ref({ tasks: [] })
const loading = ref(true)
const error = ref(null)
const selectedTask = ref(null)

// Типы узлов
const nodeTypes = {
  taskNode: markRaw(TaskNode)
}

// Элементы для VueFlow
const elements = ref([])

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

// Функция расчета позиций для узлов графа
const calculateLayout = (tasks) => {
  if (!tasks || tasks.length === 0) return { nodes: [], edges: [] }
  
  const nodes = []
  const edges = []
  
  // Простая схема: задачи располагаются слева направо
  // Зависимые задачи правее
  const levels = {}
  const taskMap = {}
  
  // Создаем карту задач
  tasks.forEach(task => {
    taskMap[task.id] = task
  })
  
  // Определяем уровень для каждой задачи
  const getLevel = (taskId, visited = new Set()) => {
    if (visited.has(taskId)) return 0
    if (levels[taskId] !== undefined) return levels[taskId]
    
    visited.add(taskId)
    const task = taskMap[taskId]
    if (!task || !task.depends_on || task.depends_on.length === 0) {
      levels[taskId] = 0
      return 0
    }
    
    const maxDepLevel = Math.max(...task.depends_on.map(depId => getLevel(depId, visited)))
    levels[taskId] = maxDepLevel + 1
    return levels[taskId]
  }
  
  tasks.forEach(task => getLevel(task.id))
  
  // Группируем задачи по уровням
  const tasksByLevel = {}
  tasks.forEach(task => {
    const level = levels[task.id] || 0
    if (!tasksByLevel[level]) tasksByLevel[level] = []
    tasksByLevel[level].push(task)
  })
  
  // Создаем узлы с позициями
  const levelGap = 250  // расстояние между уровнями
  const nodeGap = 150   // расстояние между узлами на одном уровне
  
  Object.keys(tasksByLevel).sort().forEach(level => {
    const tasksInLevel = tasksByLevel[level]
    const totalHeight = (tasksInLevel.length - 1) * nodeGap
    const startY = -totalHeight / 2
    
    tasksInLevel.forEach((task, index) => {
      nodes.push({
        id: task.id.toString(),
        type: 'taskNode',
        position: { 
          x: level * levelGap, 
          y: startY + index * nodeGap 
        },
        data: { task }
      })
    })
  })
  
  // Создаем связи
  tasks.forEach(task => {
    if (task.depends_on) {
      task.depends_on.forEach(depId => {
        edges.push({
          id: `${depId}-${task.id}`,
          source: depId.toString(),
          target: task.id.toString(),
          type: 'smoothstep',
          animated: task.status === 'in_progress',
          style: {
            stroke: task.status === 'done' ? '#4caf50' : '#999',
            strokeWidth: 2
          }
        })
      })
    }
  })
  
  return { nodes, edges }
}

const fetchProcess = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await apiCall(`/processes/${processId.value}`)
    process.value = data
    
    // Рассчитываем раскладку
    const { nodes, edges } = calculateLayout(data.tasks || [])
    elements.value = [...nodes, ...edges]
    
  } catch (err) {
    error.value = err.message
    console.error(err)
  } finally {
    loading.value = false
  }
}

const completeTask = async (taskId) => {
  try {
    const task = process.value.tasks.find(t => t.id === taskId)
    if (!task || task.status === 'done') return

    await apiCall(`/tasks/${taskId}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status: 'done' })
    })

    task.status = 'done'
    
    // Пересчитываем граф для обновления анимаций
    const { nodes, edges } = calculateLayout(process.value.tasks)
    elements.value = [...nodes, ...edges]
    
    if (selectedTask.value?.id === taskId) {
      selectedTask.value = { ...task }
    }
  } catch (err) {
    alert(`Failed to complete task: ${err.message}`)
    console.error(err)
  }
}

const onNodeClick = ({ node }) => {
  if (node.data?.task) {
    selectedTask.value = { ...node.data.task }
  }
}

const selectTask = (task) => {
  selectedTask.value = { ...task }
}

const updateTask = (updatedTask) => {
  const index = process.value.tasks.findIndex(t => t.id === updatedTask.id)
  if (index !== -1) {
    process.value.tasks[index] = { ...process.value.tasks[index], ...updatedTask }
    
    // Пересчитываем граф
    const { nodes, edges } = calculateLayout(process.value.tasks)
    elements.value = [...nodes, ...edges]
    
    selectedTask.value = null
  }
}

const goBack = () => {
  router.push('/dashboard')
}

onMounted(() => {
  fetchProcess()
})
</script>

<style scoped>
.process-view {
  width: 100%;
  height: calc(100vh - 140px);
  display: flex;
  flex-direction: column;
}

.process-toolbar {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 2px solid var(--color-text);
  background: var(--color-background);
  z-index: 10;
  flex-wrap: wrap;
  gap: 10px;
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

.process-name {
  flex: 1;
  font-weight: bold;
  font-family: var(--font-1);
  margin-left: 20px;
}

.process-status {
  font-family: var(--font-1);
  padding: 4px 12px;
  border: 2px solid var(--color-text);
  border-radius: 15px;
}

.flow-container {
  flex: 1;
  width: 100%;
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

:deep(.vue-flow__node) {
  background: transparent;
  border: none;
  padding: 0;
}

:deep(.vue-flow__edge-path) {
  stroke-width: 2;
}

:deep(.vue-flow__controls) {
  border: 2px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  overflow: hidden;
}

:deep(.vue-flow__controls-button) {
  background: var(--color-background);
  border-bottom: 1px solid var(--color-text);
  color: var(--color-text);
}

:deep(.vue-flow__background) {
  background: var(--color-background);
}
</style>