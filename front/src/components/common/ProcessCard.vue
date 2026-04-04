<template>
  <div class="process-card" :class="priorityClass" @click="$emit('click')">
    <div class="card-header">
      <h3>{{ process.name }}
      <span class="stage-badge"> > {{ process.currentStage }}</span></h3>
    </div>
    
    <div class="progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
      </div>
      <span class="progress-text">{{ process.completedTasks }}/{{ process.totalTasks }}</span>
    </div>
    
    <div class="card-footer">
      <div class="waiting-info" v-if="process.waitingFor">
        <span class="label">⏳ Waiting:</span>
        <span>{{ process.waitingFor }}</span>
      </div>
      <div class="timer">
        <span class="label">⏱️</span>
        <span>{{ process.timer }}</span>
      </div>
    </div>
    
    <div class="stages-preview">
      <template v-for="(stage, idx) in previewStages" :key="idx">
        <span class="stage-dot" :class="{ active: isActiveStage(stage) }">
          {{ stage }}
          <span v-if="!isLastStage(stage)" class="stage-arrow">→</span>
        </span>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  process: {
    type: Object,
    required: true
  },
  completed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])

const priorityClass = computed(() => {
  if (props.completed) return 'card-completed'
  const classes = {
    high: 'card-high',
    medium: 'card-medium',
    low: 'card-low'
  }
  return classes[props.process.priority] || 'card-medium'
})

const progressPercent = computed(() => {
  if (props.process.totalTasks === 0) return 0
  return (props.process.completedTasks / props.process.totalTasks) * 100
})

const previewStages = computed(() => {
  const stages = props.process.stages
  const active = props.process.activeStage
  
  if (!stages || stages.length === 0) return []
  if (stages.length <= 4) return stages
  
  const start = Math.max(0, active - 1)
  const end = Math.min(stages.length, active + 2)
  return stages.slice(start, end)
})

const isActiveStage = (stage) => {
  return stage === props.process.currentStage
}

const isLastStage = (stage) => {
  const stages = props.process.stages
  return stage === stages[stages.length - 1]
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
  min-width: 0;
  word-wrap: break-word;
}

.process-card:hover {
  transform: translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.1);
}

.card-high {
  border-left: 6px solid #f44336;
}

.card-medium {
  border-left: 6px solid #ff9800;
}

.card-low {
  border-left: 6px solid #4caf50;
}

.card-completed {
  border-left: 6px solid #9e9e9e;
  opacity: 0.7;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 8px;
}

.card-header h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: bold;
}

.stage-badge {
  font-size: 0.7rem;
  padding: 4px 10px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
}

.progress {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #4caf50;
  transition: width 0.3s;
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
}

.progress-text {
  font-size: 0.75rem;
  font-weight: bold;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 0.7rem;
  flex-wrap: wrap;
  gap: 8px;
}

.label {
  opacity: 0.6;
  margin-right: 4px;
}

.stages-preview {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed var(--color-text);
}

.stage-dot {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
}

.stage-dot.active {
  background: #2196f3;
  color: white;
}

.stage-arrow {
  opacity: 0.5;
}
</style>