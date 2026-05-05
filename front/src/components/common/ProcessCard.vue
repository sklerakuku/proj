<template>
  <div class="process-card" :class="{ 'card-completed': completed }" @click="$emit('click')">
    <!-- Строка 1: Название и таймер текущего этапа -->
    <div class="card-row-1">
      <h3>{{ process.title || 'Untitled Process' }}</h3>
      <span v-if="currentTask" class="stage-timer" :class="{ overdue: isOverdue }">
        ⏱ {{ stageDuration }}
      </span>
    </div>
    
    <!-- Строка 2: Счётчик этапов и ожидание -->
    <div class="card-row-2">
      <span class="task-counter">{{ completedTasks }}/{{ totalTasks }} tasks</span>
      <span v-if="waitingFor" class="waiting-badge">
        waiting {{ waitingFor }}
      </span>
    </div>
    
    <!-- Строка 3: Хлебные крошки этапов -->
    <div class="card-row-3">
      <div class="breadcrumbs">
        <template v-for="(task, idx) in tasksList" :key="task.id">
          <span v-if="idx > 0" class="crumb-arrow">></span>
          <span class="crumb" :class="crumbClass(task)">{{ task.title }}</span>
        </template>
        <span v-if="tasksList.length === 0" class="crumb">No tasks yet</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  process: { type: Object, required: true },
  completed: { type: Boolean, default: false }
})

defineEmits(['click'])

const tasksList = computed(() => props.process.tasks || [])
const totalTasks = computed(() => tasksList.value.length)
const completedTasks = computed(() => tasksList.value.filter(t => t.status === 'done').length)

const currentTask = computed(() => {
  return tasksList.value.find(t => t.status === 'in_progress') || null
})

const stageDuration = computed(() => {
  if (!currentTask.value?.started_at) return '--'
  const start = new Date(currentTask.value.started_at)
  const now = new Date()
  const diff = now - start
  const hours = Math.floor(diff / (1000 * 60 * 60))
  if (hours < 1) return '< 1h'
  if (hours < 24) return `${hours}h`
  const days = Math.floor(hours / 24)
  return `${days}d ${hours % 24}h`
})

const isOverdue = computed(() => {
  if (!currentTask.value?.started_at) return false
  const start = new Date(currentTask.value.started_at)
  const now = new Date()
  const days = (now - start) / (1000 * 60 * 60 * 24)
  return days > 7
})

const waitingFor = computed(() => {
  const task = currentTask.value
  if (!task) return null
  return task.for_role || 'assignee'
})

const crumbClass = (task) => ({
  'crumb-done': task.status === 'done',
  'crumb-active': task.status === 'in_progress',
  'crumb-pending': task.status === 'pending'
})
</script>

<style scoped>
.process-card {
  padding: 1.2rem;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sketch);
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--color-background);
}
.process-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-card-hover);
}
.card-completed {
  opacity: 0.6;
  border-left: 6px solid var(--status-done);
}

.card-row-1 {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.card-row-1 h3 {
  margin: 0;
  font-family: var(--font-1);
  font-size: 1.1rem;
}
.stage-timer {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  padding: 2px 10px;
  border-radius: var(--radius-sm);
  background: var(--color-hover);
  white-space: nowrap;
}
.stage-timer.overdue {
  color: var(--color-danger);
  background: var(--color-danger-light);
}

.card-row-2 {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 0.85rem;
}
.task-counter {
  opacity: 0.7;
  font-family: var(--font-1);
}
.waiting-badge {
  padding: 2px 10px;
  border-radius: var(--radius-sm);
  background: var(--color-info-light);
  color: var(--color-info);
  font-style: italic;
}

.card-row-3 {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed var(--color-muted-light);
}

.breadcrumbs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  font-size: 0.7rem;
}

.crumb {
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  background: var(--color-hover);
  white-space: nowrap;
}
.crumb-done {
  background: var(--status-done-bg);
  color: var(--status-done);
}
.crumb-active {
  background: var(--status-progress-bg);
  color: var(--status-progress);
  font-weight: bold;
}
.crumb-pending {
  background: var(--status-pending-bg);
  color: var(--status-pending);
}

.crumb-arrow {
  opacity: 0.4;
  font-size: 0.6rem;
}
</style>