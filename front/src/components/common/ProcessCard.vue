<template>
  <div class="process-card" @click="$emit('click')">
    <div class="card-header">
      <h3>{{ process.title || 'Untitled Process' }}</h3>
      <span class="status-badge" :class="statusClass">{{ process.status }}</span>
    </div>
    
    <div class="card-body">
      <div class="progress-info">
        <span>{{ completedTasks }} / {{ totalTasks }} tasks</span>
      </div>
      <div class="timer" v-if="process.started_at">
        <span>⏱️ Started: {{ formatDate(process.started_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  process: {
    type: Object,
    required: true
  }
})

defineEmits(['click'])

const statusClass = computed(() => {
  switch (props.process.status) {
    case 'done': return 'status-done'
    case 'in_progress': return 'status-progress'
    case 'draft': return 'status-draft'
    default: return 'status-pending'
  }
})

const totalTasks = computed(() => props.process.tasks?.length || 0)
const completedTasks = computed(() => props.process.tasks?.filter(t => t.status === 'done').length || 0)

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  try {
    return new Date(dateStr).toLocaleDateString()
  } catch {
    return dateStr
  }
}
</script>

<style scoped>
.process-card {
  padding: 1.2rem;
  border: 2px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  cursor: pointer;
  transition: all 0.1s ease;
  background: var(--color-background);
}

.process-card:hover {
  transform: translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.card-header h3 {
  margin: 0;
  font-family: var(--font-1);
}

.status-badge {
  font-family: var(--font-1);
  font-size: 0.8rem;
  padding: 4px 10px;
  border-radius: 15px;
  background: rgba(0,0,0,0.05);
}

.status-draft { color: #888; }
.status-progress { color: #ff9800; }
.status-done { color: #4caf50; }

.card-body {
  font-size: 0.9rem;
}

.progress-info, .timer {
  margin: 4px 0;
}
</style>