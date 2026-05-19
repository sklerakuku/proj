<template>
  <div class="process-view">
    <div class="process-toolbar">
      <button class="tool-btn" @click="goBack">← Back</button>
      <div class="process-name">{{ process.title }}</div>
      <div class="process-status" :class="'status-' + process.status">{{ process.status }}</div>
      <button class="tool-btn" @click="showAddTask = true">+ Task</button>
      <button class="tool-btn" @click="showAllAttachments = true">📎 All Attachments</button>
      <button 
        v-if="process.status === 'done'" 
        class="tool-btn" 
        @click="archiveProcess"
      >Archive</button>
      <button class="tool-btn" v-if="selectedTaskForLink" @click="cancelLink">✕</button>
    </div>

    <div v-if="loading" class="status-message">Loading...</div>
    <div v-else-if="error" class="status-message error">{{ error }}</div>
    
    <div v-else class="flow-container">
      <div v-if="selectedTaskForLink" class="link-hint">
        Click on another task to create dependency from "{{ selectedTaskForLink.title }}"
      </div>
      <VueFlow
        v-model="elements"
        :node-types="nodeTypes"
        :default-viewport="{ zoom: 1 }"
        :min-zoom="0.2"
        :max-zoom="4"
        @node-click="onNodeClick"
        @connect="onConnect"
      >
        <Background pattern-color="#aaa" :gap="16" />
        <Controls />
        
        <template #node-taskNode="{ data }">
          <TaskNode 
            :task="data.task"
            @complete="(id) => completeTask(id)"
            @click="selectTask(data.task)"
            @link="startLinking(data.task)"
          />
        </template>
      </VueFlow>
    </div>

    <TaskSidebar 
      v-if="selectedTask"
      :task="selectedTask"
      @close="selectedTask = null"
      @update="updateTask"
      @complete="(id) => completeTask(id)"
    />

<!-- Модалка со всеми вложениями -->
<div v-if="showAllAttachments" class="modal-overlay" @click.self="showAllAttachments = false">
  <div class="modal modal-large">
    <h3>All Attachments</h3>
    <div class="attachments-list">
      <div v-for="task in process.tasks" :key="task.id" class="task-attachments-group">
        <h4>{{ task.title }}</h4>
        <div v-if="task.attachments?.length" class="attachment-items">
          <div v-for="(att, idx) in task.attachments" :key="idx" class="attachment-item">
            <a :href="att.file_path" download>📄 {{ att.file_path.split('/').pop() }}</a>
            <span class="file-size">({{ att.file_size_kb }} KB)</span>
          </div>
        </div>
        <div v-else class="no-attachments">No attachments</div>
      </div>
    </div>
    <button class="btn-cancel" @click="showAllAttachments = false">Close</button>
  </div>
</div>

    <!-- Модалка добавления задачи -->
<div v-if="showAddTask" class="modal-overlay" @click.self="showAddTask = false">
  <div class="modal">
    <h3>Add New Task</h3>
    <form @submit.prevent="addTask">
      <input v-model="newTask.title" placeholder="Task title" class="sketch-input" required />
      <select v-model="newTask.role" class="sketch-input">
        <option value="worker">Worker</option>
        <option value="manager">Manager</option>
        <option value="admin">Admin</option>
      </select>
      <select v-model="newTask.dependsOn" class="sketch-input">
        <option value="">No dependency (standalone)</option>
        <option v-for="t in process.tasks" :key="t.id" :value="t.id">
          After: {{ t.title }}
        </option>
      </select>
      <div class="modal-buttons">
        <button type="submit" class="btn-save">Add</button>
        <button type="button" class="btn-cancel" @click="showAddTask = false">Cancel</button>
      </div>
    </form>
  </div>
</div>
</div>
</template>

<script setup>
import { ref, markRaw, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import TaskNode from './common/TaskNode.vue'
import TaskSidebar from './common/TaskSidebar.vue'

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
const selectedTaskForLink = ref(null)
const showAddTask = ref(false)
const newTask = ref({ title: '', role: 'worker', dependsOn: '' })

const nodeTypes = { taskNode: markRaw(TaskNode) }
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
    throw new Error(errorData.error || `Request failed`)
  }

  return response.json()
}

const calculateLayout = (tasks) => {
  if (!tasks || tasks.length === 0) return { nodes: [], edges: [] }
  
  const nodes = []
  const edges = []
  const levels = {}
  const taskMap = {}
  
  tasks.forEach(task => { taskMap[task.id] = task })
  
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
  
  const tasksByLevel = {}
  tasks.forEach(task => {
    const level = levels[task.id] || 0
    if (!tasksByLevel[level]) tasksByLevel[level] = []
    tasksByLevel[level].push(task)
  })
  
  const levelGap = 250
  const nodeGap = 150
  
  Object.keys(tasksByLevel).sort().forEach(level => {
    const tasksInLevel = tasksByLevel[level]
    const totalHeight = (tasksInLevel.length - 1) * nodeGap
    const startY = -totalHeight / 2
    
    tasksInLevel.forEach((task, index) => {
      nodes.push({
        id: task.id.toString(),
        type: 'taskNode',
        position: { x: level * levelGap, y: startY + index * nodeGap },
        data: { task }
      })
    })
  })
  
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
            stroke: task.status === 'done' ? 'var(--status-done)' : 'var(--color-muted)',
            strokeWidth: 2
          }
        })
      })
    }
  })
  
  return { nodes, edges }
}

const updateProcessStatus = async () => {
  const tasks = process.value.tasks || []
  if (tasks.length === 0) return
  
  const allDone = tasks.every(t => t.status === 'done')
  if (allDone && process.value.status !== 'done') {
    process.value.status = 'done'
    const { nodes, edges } = calculateLayout(tasks)
    elements.value = [...nodes, ...edges]
  }
}

const archiveProcess = async () => {
  if (!confirm('Archive this process?')) return
  try {
    await apiCall(`/processes/archive/${processId.value}`, { method: 'PATCH' })
    process.value.status = 'archived'
  } catch (err) {
    alert(`Failed to archive: ${err.message}`)
  }
}

const fetchProcess = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await apiCall(`/processes/${processId.value}`)
    process.value = data
    const { nodes, edges } = calculateLayout(data.tasks || [])
    elements.value = [...nodes, ...edges]
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const addTask = async () => {
  try {
    const taskId = Date.now()
    const dependsOn = newTask.value.dependsOn ? [parseInt(newTask.value.dependsOn)] : []
    
    const newTaskObj = {
      id: taskId,
      process_id: parseInt(processId.value),
      title: newTask.value.title,
      for_role: newTask.value.role,
      status: 'pending',
      depends_on: dependsOn
    }
    
    if (!process.value.tasks) {
      process.value.tasks = []
    }
    process.value.tasks.push(newTaskObj)
    
    const { nodes, edges } = calculateLayout(process.value.tasks)
    elements.value = [...nodes, ...edges]
    showAddTask.value = false
    newTask.value = { title: '', role: 'worker', dependsOn: '' }
  } catch (err) {
    alert(`Failed to add task: ${err.message}`)
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
    await updateProcessStatus()
    
    const { nodes, edges } = calculateLayout(process.value.tasks)
    elements.value = [...nodes, ...edges]
    
    if (selectedTask.value?.id === taskId) {
      selectedTask.value = { ...task }
    }
  } catch (err) {
    alert(`Failed to complete task: ${err.message}`)
  }
}

const startLinking = (task) => {
  selectedTaskForLink.value = task
}

const cancelLink = () => {
  selectedTaskForLink.value = null
}

const onNodeClick = ({ node }) => {
  if (selectedTaskForLink.value && node.data?.task) {
    // Создаем связь
    const sourceTask = selectedTaskForLink.value
    const targetTask = node.data.task
    
    if (sourceTask.id === targetTask.id) return
    
    if (!targetTask.depends_on) targetTask.depends_on = []
    if (!targetTask.depends_on.includes(sourceTask.id)) {
      targetTask.depends_on.push(sourceTask.id)
      const { nodes, edges } = calculateLayout(process.value.tasks)
      elements.value = [...nodes, ...edges]
    }
    
    selectedTaskForLink.value = null
    return
  }
  
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
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
}

.process-toolbar {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 2px solid var(--color-border);
  background: var(--color-background);
  z-index: 10;
  flex-wrap: wrap;
  gap: 10px;
}

.tool-btn {
  padding: 6px 12px;
  background: transparent;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-btn);
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
}
.tool-btn:hover { transform: translateY(-1px); }

.process-name { flex: 1; font-weight: bold; font-family: var(--font-1); margin-left: 20px; }

.process-status {
  font-family: var(--font-1);
  padding: 4px 12px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
}
.status-draft { color: var(--status-draft); }
.status-in_progress { color: var(--status-progress); }
.status-done { color: var(--status-done); }

.flow-container { flex: 1; width: 100%; position: relative; }

.link-hint {
  position: absolute;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--color-info);
  color: white;
  padding: 8px 20px;
  border-radius: var(--radius-md);
  font-family: var(--font-1);
  z-index: 20;
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
.sketch-input {
  width: 100%;
  padding: 10px;
  margin-bottom: 10px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-1);
  background: transparent;
  color: var(--color-text);
}
.modal-buttons { display: flex; gap: 10px; margin-top: 1rem; }
.btn-save, .btn-cancel {
  padding: 8px 20px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-btn);
  cursor: pointer;
  font-family: var(--font-1);
  background: transparent;
  color: var(--color-text);
}
.btn-save { border-color: var(--color-success); color: var(--color-success); }
.btn-cancel { border-color: var(--color-muted); color: var(--color-muted); }
</style>