<template>
  <div class="task-sidebar">
    <div class="sidebar-header">
      <h3>{{ task.title }}</h3>
      <button class="edit-icon" @click="editTask">✏️</button>
      <button class="close-btn" @click="$emit('close')">×</button>
    </div>
    
    <div class="sidebar-content">
      <div class="waiting-info">
        <span class="label">Waiting for:</span>
        <span>{{ task.waitingFor || 'None' }}</span>
        <span class="timer" v-if="task.waitingTimer">⏱️ {{ task.waitingTimer }}</span>
      </div>
      
      <div class="attachments">
        <div class="attachments-header">
          <span>📎 Attachments</span>
          <button @click="addAttachment">+ Add</button>
        </div>
        <div class="attachment-list">
          <div v-for="file in task.attachments" :key="file" class="attachment-item">
            📄 {{ file }}
          </div>
        </div>
      </div>
      
      <div class="description-editor">
        <div class="editor-toolbar">
          <button @click="formatText('bold')"><b>B</b></button>
          <button @click="formatText('italic')"><i>I</i></button>
          <button @click="formatText('underline')"><u>U</u></button>
        </div>
        <textarea 
          v-model="editedDescription" 
          placeholder="Task description..."
          class="description-input"
        ></textarea>
      </div>
      
      <button class="complete-task-btn" @click="completeTask" v-if="task.status !== 'done'">
        ✓ Complete Task
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'update'])

const editedDescription = ref(props.task.description)

const editTask = () => {
  // Открыть режим редактирования названия
  const newTitle = prompt('Edit task title:', props.task.title)
  if (newTitle) {
    emit('update', { ...props.task, title: newTitle })
  }
}

const addAttachment = () => {
  const file = prompt('Enter file name:')
  if (file) {
    const attachments = [...(props.task.attachments || []), file]
    emit('update', { ...props.task, attachments })
  }
}

const formatText = (format) => {
  const textarea = document.querySelector('.description-input')
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = editedDescription.value.substring(start, end)
  
  let wrapped = ''
  if (format === 'bold') wrapped = `**${selectedText}**`
  if (format === 'italic') wrapped = `*${selectedText}*`
  if (format === 'underline') wrapped = `__${selectedText}__`
  
  editedDescription.value = editedDescription.value.substring(0, start) + wrapped + editedDescription.value.substring(end)
}

const completeTask = () => {
  if (confirm('Complete this task?')) {
    emit('update', { ...props.task, status: 'done' })
  }
}
</script>

<style scoped>
.task-sidebar {
  position: fixed;
  right: 0;
  top: 0;
  width: 380px;
  height: 100vh;
  background: var(--color-background);
  border-left: 3px solid var(--color-text);
  box-shadow: -5px 0 20px rgba(0, 0, 0, 0.1);
  z-index: 200;
  display: flex;
  flex-direction: column;
  font-family: var(--font-1);
}

.sidebar-header {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 2px solid var(--color-text);
  gap: 10px;
}

.sidebar-header h3 {
  flex: 1;
  margin: 0;
  font-size: 1.2rem;
}

.edit-icon, .close-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 4px 8px;
}

.close-btn {
  font-size: 1.5rem;
}

.sidebar-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.waiting-info {
  padding: 12px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 15px;
}

.waiting-info .label {
  font-weight: bold;
  margin-right: 8px;
}

.timer {
  margin-left: 10px;
  font-family: monospace;
}

.attachments-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.attachments-header button {
  background: none;
  border: 1px solid var(--color-text);
  border-radius: 12px;
  padding: 2px 8px;
  cursor: pointer;
}

.attachment-item {
  padding: 6px;
  margin: 4px 0;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 8px;
  font-size: 0.85rem;
}

.editor-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.editor-toolbar button {
  width: 32px;
  height: 32px;
  border: 1px solid var(--color-text);
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
}

.description-input {
  width: 100%;
  min-height: 150px;
  padding: 10px;
  border: 2px solid var(--color-text);
  border-radius: 15px;
  background: transparent;
  font-family: inherit;
  resize: vertical;
}

.complete-task-btn {
  padding: 12px;
  background: #4caf50;
  color: white;
  border: none;
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  cursor: pointer;
  font-family: var(--font-1);
  margin-top: 20px;
}

@media (max-width: 768px) {
  .task-sidebar {
    width: 100%;
    height: 50vh;
    top: auto;
    bottom: 0;
    border-left: none;
    border-top: 3px solid var(--color-text);
  }
}
</style>