<template>
  <div class="task-node" :class="statusClass" @click.stop="$emit('click')">
    <div class="task-header">
      <span class="task-title">{{ task.title }}</span>
      <div class="task-actions">
        <button class="link-btn" @click.stop="$emit('link')" title="Create dependency">🔗</button>
        <button v-if="task.status !== 'done'" class="complete-btn" @click.stop="$emit('complete', task.id)">✓</button>
      </div>
    </div>
    <div class="task-footer">
      <span class="role-badge">{{ task.for_role }}</span>
      <span class="progress-badge">{{ task.status === 'done' ? '✓' : '◯' }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ task: { type: Object, required: true } })
defineEmits(['click', 'complete', 'link'])

const statusClass = computed(() => ({
  'task-pending': props.task.status === 'pending',
  'task-progress': props.task.status === 'in_progress',
  'task-done': props.task.status === 'done'
}))
</script>

<style scoped>
.task-node {
  padding: 10px;
  border: 3px solid var(--color-border);
  border-radius: var(--radius-sketch);
  background: var(--color-background);
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: var(--shadow-card);
  width: 200px;
}
.task-node:hover { transform: translateY(-2px); }
.task-pending { border-left: 6px solid var(--status-pending); }
.task-progress { border-left: 6px solid var(--status-progress); }
.task-done { border-left: 6px solid var(--status-done); opacity: 0.7; }
.task-header { display: flex; justify-content: space-between; align-items: center; }
.task-title { font-weight: bold; font-size: 0.9rem; }
.task-actions { display: flex; gap: 4px; }
.link-btn {
  background: none;
  border: 1px solid var(--color-muted);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 0.7rem;
  padding: 1px 4px;
}
.complete-btn {
  background: var(--color-success);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  padding: 1px 6px;
}
.task-footer { display: flex; justify-content: space-between; margin-top: 8px; font-size: 0.7rem; }
.role-badge { background: var(--color-info); color: white; padding: 1px 6px; border-radius: 8px; }
</style>