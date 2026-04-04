<template>
  <div class="task-node" :class="statusClass" @click.stop="$emit('click')">
    <div class="task-header">
      <span class="task-title">{{ task.title }}</span>
      <button v-if="task.status !== 'done'" class="complete-btn" @click.stop="$emit('complete', task.id)">
        ✓ End
      </button>
    </div>
    
    <div class="task-description">{{ task.description }}</div>
    
    <div class="task-footer">
      <div class="task-actions">
        <button class="icon-btn" @click.stop="$emit('edit')">📝</button>
        <button class="icon-btn" @click.stop="$emit('attach')">📎</button>
      </div>
      <div class="task-assignee">
        <span class="assignee-name">{{ task.assignedTo }}</span>
      </div>
    </div>
    
    <div class="task-meta">
      <span class="role-badge">{{ task.role }}</span>
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
  },
  selected: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click', 'complete', 'edit', 'attach'])

const statusClass = computed(() => {
  const classes = {
    pending: 'task-pending',
    in_progress: 'task-progress',
    done: 'task-done'
  }
  return classes[props.task.status] || 'task-pending'
})
</script>

<style scoped>
.task-node {
  width: 190px;
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
  box-shadow: 4px 4px 0 rgba(0, 0, 0, 0.15);
}

.task-pending {
  border-left: 6px solid #9e9e9e;
}

.task-progress {
  border-left: 6px solid #ff9800;
}

.task-done {
  border-left: 6px solid #4caf50;
  opacity: 0.7;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.task-title {
  font-weight: bold;
  font-size: 0.85rem;
}

.complete-btn {
  padding: 2px 6px;
  background: #4caf50;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 0.65rem;
  cursor: pointer;
}

.task-description {
  font-size: 0.7rem;
  opacity: 0.7;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-actions {
  display: flex;
  gap: 6px;
}

.icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.8rem;
}

.task-assignee {
  font-size: 0.65rem;
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 10px;
}

.task-meta {
  display: flex;
  justify-content: space-between;
  margin-top: 6px;
  font-size: 0.6rem;
}

.role-badge {
  background: #2196f3;
  color: white;
  padding: 1px 6px;
  border-radius: 8px;
}
</style>