<template>
  <div class="task-sidebar">
    <div class="sidebar-header">
      <h3>{{ task.title }}</h3>
      <button class="edit-icon" @click="editTask">✏️</button>
      <button class="close-btn" @click="$emit('close')">×</button>
    </div>
    
    <div class="sidebar-content">
      <div class="waiting-info">
        <span class="label">Role:</span>
        <span class="role-badge">{{ task.for_role || 'None' }}</span>
      </div>
      
      <div class="attachments">
        <div class="attachments-header">
          <span>📎 Attachments</span>
          <label class="add-attach-btn">
            + Add
            <input type="file" @change="handleFileSelect" style="display:none" multiple />
          </label>
        </div>
        <div class="attachment-list">
          <div v-for="(file, idx) in task.attachments" :key="idx" class="attachment-item">
            📄 {{ typeof file === 'string' ? file : file.name }}
            <button class="remove-file" @click="removeAttachment(idx)">×</button>
          </div>
          <div v-if="!task.attachments || task.attachments.length === 0" class="no-files">
            No attachments
          </div>
        </div>
      </div>
      
      <div class="description-editor">
        <div class="editor-toolbar">
          <button @click="formatText('bold')" title="Bold"><b>B</b></button>
          <button @click="formatText('italic')" title="Italic"><i>I</i></button>
          <button @click="formatText('heading')" title="Heading">H</button>
          <button @click="formatText('list')" title="List">•</button>
          <button @click="formatText('code')" title="Code">&lt;/&gt;</button>
        </div>
        <div class="editor-area">
          <textarea 
            v-model="editedDescription" 
            @blur="saveDescription"
            placeholder="Task description..."
            class="description-input"
            @input="updatePreview"
          ></textarea>
        </div>
        <div v-if="showPreview" class="preview-pane" v-html="renderedMarkdown"></div>
        <button @click="showPreview = !showPreview" class="preview-toggle">
          {{ showPreview ? '✏️ Edit' : '👁 Preview' }}
        </button>
      </div>

      <button class="view-all-attachments" @click="showAllAttachments">📎 All Attachments ({{ totalAttachments }})</button>
      
      <button class="complete-task-btn" @click="completeTask" v-if="task.status !== 'done'">
        Complete Task
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  task: { type: Object, required: true }
})

const emit = defineEmits(['close', 'update', 'complete'])

const editedDescription = ref(props.task.comment || props.task.description || '')
const showPreview = ref(false)

const saveDescription = async () => {
  try {
    await apiCall(`/tasks/${props.task.id}/comment`, {
      method: 'PATCH',
      body: JSON.stringify({ comment: editedDescription.value })
    })
  } catch (e) {
    console.error('Failed to save description', e)
  }
}

const handleFileSelect = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    await fetch(`/tasks/${props.task.id}/attachments`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('auth_token')}` },
      body: formData
    })
    // Обновить список вложений
    await loadAttachments()
  } catch (e) {
    alert('Failed to upload file')
  }
}

// Простая функция рендеринга markdown (без библиотеки)
const renderMarkdown = (text) => {
  if (!text) return ''
  let html = text
    // Bold
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    // Italic
    .replace(/\*(.+?)\*/g, '<em>$1</em>')
    // Heading
    .replace(/^### (.+)$/gm, '<h4>$1</h4>')
    .replace(/^## (.+)$/gm, '<h3>$1</h3>')
    .replace(/^# (.+)$/gm, '<h2>$1</h2>')
    // Code
    .replace(/`(.+?)`/g, '<code>$1</code>')
    // List
    .replace(/^- (.+)$/gm, '<li>$1</li>')
    // Underline
    .replace(/__(.+?)__/g, '<u>$1</u>')
    // Line breaks
    .replace(/\n/g, '<br>')
  
  // Wrap consecutive <li> in <ul>
  html = html.replace(/(<li>.*<\/li>)/s, (match) => {
    if (match.includes('<li>')) return '<ul>' + match + '</ul>'
    return match
  })
  
  return html
}

const renderedMarkdown = computed(() => renderMarkdown(editedDescription.value))

const updatePreview = () => {
  // live preview
}

const editTask = () => {
  const newTitle = prompt('Edit task title:', props.task.title)
  if (newTitle) {
    emit('update', { ...props.task, title: newTitle })
  }
}



const removeAttachment = (idx) => {
  const attachments = [...(props.task.attachments || [])]
  attachments.splice(idx, 1)
  emit('update', { ...props.task, attachments })
}

const formatText = (format) => {
  const textarea = document.querySelector('.description-input')
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = editedDescription.value.substring(start, end)
  
  let wrapped = ''
  switch (format) {
    case 'bold': wrapped = `**${selectedText || 'bold text'}**`; break
    case 'italic': wrapped = `*${selectedText || 'italic text'}*`; break
    case 'heading': wrapped = `\n## ${selectedText || 'Heading'}\n`; break
    case 'list': wrapped = `\n- ${selectedText || 'list item'}`; break
    case 'code': wrapped = `\`${selectedText || 'code'}\``; break
  }
  
  editedDescription.value = 
    editedDescription.value.substring(0, start) + 
    wrapped + 
    editedDescription.value.substring(end)
}

const completeTask = () => {
  if (confirm('Complete this task?')) {
    // Эмитим complete вместо update для вызова API
    emit('complete', props.task.id)
  }
}
</script>

<style scoped>
.task-sidebar {
  position: fixed;
  right: 0;
  top: 0;
  width: 400px;
  height: 100vh;
  background: var(--color-background);
  border-left: 3px solid var(--color-border);
  box-shadow: -5px 0 20px rgba(0, 0, 0, 0.15);
  z-index: 200;
  display: flex;
  flex-direction: column;
  font-family: var(--font-1);
}

.sidebar-header {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 2px solid var(--color-border);
  gap: 10px;
}

.sidebar-header h3 { flex: 1; margin: 0; font-size: 1.2rem; }

.edit-icon, .close-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 4px 8px;
  color: var(--color-text);
}
.close-btn { font-size: 1.5rem; }

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
  background: var(--color-hover);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: 8px;
}
.waiting-info .label { font-weight: bold; }
.role-badge {
  background: var(--color-info);
  color: white;
  padding: 2px 10px;
  border-radius: 8px;
  font-size: 0.85rem;
}

.attachments-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-weight: bold;
}

.add-attach-btn {
  background: none;
  border: 2px solid var(--color-success);
  color: var(--color-success);
  border-radius: var(--radius-sm);
  padding: 2px 10px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.15s;
}
.add-attach-btn:hover {
  background: var(--color-success);
  color: white;
}

.attachment-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.attachment-item {
  padding: 6px 10px;
  background: var(--color-hover);
  border-radius: 8px;
  font-size: 0.85rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.remove-file {
  background: none;
  border: none;
  color: var(--color-danger);
  cursor: pointer;
  font-size: 1rem;
}

.no-files {
  font-size: 0.8rem;
  opacity: 0.5;
  text-align: center;
  padding: 8px;
}

.editor-toolbar {
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
}

.editor-toolbar button {
  width: 32px;
  height: 32px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
  transition: all 0.15s;
}
.editor-toolbar button:hover {
  background: var(--color-hover);
}

.editor-area {
  margin-bottom: 8px;
}

.description-input {
  width: 100%;
  min-height: 120px;
  padding: 10px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  background: transparent;
  font-family: var(--font-1);
  color: var(--color-text);
  resize: vertical;
  line-height: 1.5;
}

.preview-toggle {
  background: none;
  border: 1px solid var(--color-muted);
  border-radius: var(--radius-sm);
  padding: 4px 12px;
  cursor: pointer;
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  align-self: flex-start;
}

.preview-pane {
  padding: 10px;
  border: 1px dashed var(--color-muted-light);
  border-radius: var(--radius-md);
  min-height: 60px;
  line-height: 1.5;
  font-size: 0.9rem;
}
.preview-pane :deep(strong) { font-weight: bold; }
.preview-pane :deep(em) { font-style: italic; }
.preview-pane :deep(h2), .preview-pane :deep(h3), .preview-pane :deep(h4) { margin: 8px 0 4px; }
.preview-pane :deep(code) { 
  background: var(--color-hover); 
  padding: 1px 5px; 
  border-radius: 4px; 
  font-family: var(--font-mono);
  font-size: 0.85em;
}
.preview-pane :deep(ul) { padding-left: 20px; }
.preview-pane :deep(li) { margin: 2px 0; }

.complete-task-btn {
  padding: 12px;
  background: var(--color-success);
  color: white;
  border: none;
  border-radius: var(--radius-sketch);
  cursor: pointer;
  font-family: var(--font-1);
  font-size: 1rem;
  margin-top: 10px;
  transition: all 0.2s;
}
.complete-task-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .task-sidebar {
    width: 100%;
    height: 60vh;
    top: auto;
    bottom: 0;
    border-left: none;
    border-top: 3px solid var(--color-border);
  }
}
</style>