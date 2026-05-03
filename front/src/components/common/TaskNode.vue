<template>
  <div class="task-node" :class="statusClass" @click.stop="$emit('click')">
    <div class="task-header">
      <span class="task-title">{{ task.title }}</span>
      <button v-if="task.status !== 'done'" class="complete-btn" @click.stop="$emit('complete', task.id)">
        ✓
      </button>
    </div>
    <div class="task-footer">
      <span class="role-badge">{{ task.for_role }}</span>
      <span class="progress-badge">{{ task.status === 'done' ? '✓' : '◯' }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

defineEmits(['click', 'complete'])

const statusClass = computed(() => {
  switch (props.task.status) {
    case 'pending': return 'task-pending'
    case 'in_progress': return 'task-progress'
    case 'done': return 'task-done'
    default: return 'task-pending'
  }
})
</script>

<style scoped>
.task-node {
  padding: 10px;
  border: 3px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  background: var(--color-background);
  cursor: pointer;
  transition: all 0.1s ease;
  box-shadow: 2px 2px 0 rgba(0, 0, 0, 0.1);
}

.task-node:hover {
  transform: translateY(-2px);
}

.task-pending { border-left: 6px solid #9e9e9e; }
.task-progress { border-left: 6px solid #ff9800; }
.task-done { border-left: 6px solid #4caf50; opacity: 0.7; }

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-title { font-weight: bold; font-size: 0.9rem; }

.complete-btn {
  background: #4caf50;
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  font-size: 0.7rem;
}

.role-badge { background: #2196f3; color: white; padding: 1px 6px; border-radius: 8px; }
</style>