<template>
  <div class="admin-panel">
    <h2>Admin Panel</h2>
    
    <div class="tabs">
      <button @click="activeTab = 'users'" :class="{ active: activeTab === 'users' }">Users</button>
      <button @click="activeTab = 'templates'" :class="{ active: activeTab === 'templates' }">Templates</button>
      <button @click="activeTab = 'processes'" :class="{ active: activeTab === 'processes' }">Processes</button>
    </div>

    <!-- Users -->
    <div v-if="activeTab === 'users'" class="table-container">
      <h3>Users ({{ users.length }})</h3>
      <table>
        <thead>
          <tr><th>ID</th><th>Username</th><th>Role</th><th>Created</th><th>Actions</th></tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>
              <input v-if="editingUser?.id === user.id" v-model="editForm.username" class="edit-input" />
              <span v-else>{{ user.username }}</span>
            </td>
            <td>
              <select v-if="editingUser?.id === user.id" v-model="editForm.role" class="edit-input">
                <option value="admin">Admin</option>
                <option value="worker">Worker</option>
                <option value="manager">Manager</option>
              </select>
              <span v-else>{{ user.role }}</span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td class="actions">
              <template v-if="editingUser?.id === user.id">
                <button @click="saveUser(user.id)" class="btn-save">Save</button>
                <button @click="cancelEdit" class="btn-cancel">Cancel</button>
              </template>
              <template v-else>
                <button @click="editUser(user)" class="btn-edit">Edit</button>
                <button @click="deleteUser(user.id)" class="btn-delete">Delete</button>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Templates -->
    <div v-if="activeTab === 'templates'" class="table-container">
      <div class="tab-header">
        <h3>Templates ({{ templates.length }})</h3>
        <button class="btn-add" @click="showAddTemplate = true">+ New Template</button>
      </div>
      <table>
        <thead>
          <tr><th>ID</th><th>Name</th><th>Description</th><th>Actions</th></tr>
        </thead>
        <tbody>
          <tr v-for="tpl in templates" :key="tpl.id">
            <td>{{ tpl.id }}</td>
            <td>{{ tpl.name }}</td>
            <td>{{ tpl.description }}</td>
            <td>
              <button @click="deleteTemplate(tpl.id)" class="btn-delete">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- В AdminPanel.vue, внутри модалки создания шаблона -->
      <div v-if="showAddTemplate" class="modal-overlay" @click.self="showAddTemplate = false">
        <div class="modal modal-large">
          <h3>Create Template</h3>
          <form @submit.prevent="createTemplate">
            <input v-model="newTemplate.name" placeholder="Template name" class="sketch-input" required />
            <textarea v-model="newTemplate.description" placeholder="Description" class="sketch-input" rows="2"></textarea>
            
            <!-- Список задач шаблона -->
            <div class="template-tasks-section">
              <h4>Tasks</h4>
              <div v-for="(task, idx) in newTemplate.tasks" :key="idx" class="template-task-row">
                <input v-model="task.title" placeholder="Task title" class="sketch-input small" />
                <select v-model="task.for_role" class="sketch-input small">
                  <option value="worker">Worker</option>
                  <option value="manager">Manager</option>
                  <option value="admin">Admin</option>
                </select>
                <label><input type="checkbox" v-model="task.is_file_required" /> File required</label>
                <input type="number" v-model="task.plan_hours" placeholder="Hours" class="sketch-input tiny" />
                <button type="button" class="btn-remove" @click="removeTask(idx)">✕</button>
              </div>
              <button type="button" class="btn-add-task" @click="addTask">+ Add Task</button>
            </div>
            
            <!-- Зависимости между задачами -->
            <div class="dependencies-section" v-if="newTemplate.tasks.length > 1">
              <h4>Dependencies</h4>
              <div v-for="(task, idx) in newTemplate.tasks" :key="idx" class="dependency-row">
                <span class="task-name">{{ task.title || `Task ${idx+1}` }}</span>
                <span>depends on:</span>
                <select v-model="task.depends_on" multiple class="sketch-input small">
                  <option v-for="(other, oidx) in newTemplate.tasks" :key="oidx" :value="oidx" :disabled="oidx === idx">
                    {{ other.title || `Task ${oidx+1}` }}
                  </option>
                </select>
              </div>
            </div>
            
            <div class="modal-buttons">
              <button type="submit" class="btn-save">Create Template</button>
              <button type="button" class="btn-cancel" @click="showAddTemplate = false">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Processes -->
    <div v-if="activeTab === 'processes'" class="table-container">
      <h3>Processes ({{ processes.length }})</h3>
      <table>
        <thead>
          <tr><th>ID</th><th>Title</th><th>Status</th><th>Started</th><th>Actions</th></tr>
        </thead>
        <tbody>
          <tr v-for="proc in processes" :key="proc.id">
            <td>{{ proc.id }}</td>
            <td>{{ proc.title }}</td>
            <td>
              <select v-model="proc.status" @change="updateProcessStatus(proc)" class="edit-input">
                <option value="draft">Draft</option>
                <option value="in_progress">In Progress</option>
                <option value="done">Done</option>
              </select>
            </td>
            <td>{{ formatDate(proc.started_at) }}</td>
            <td>
              <button @click="deleteProcess(proc.id)" class="btn-delete">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const activeTab = ref('users')

const users = ref([])
const templates = ref([])
const processes = ref([])

const editingUser = ref(null)
const editForm = ref({ username: '', role: '' })
const showAddTemplate = ref(false)

const apiCall = async (url, options = {}) => {
  const token = localStorage.getItem('auth_token')
  if (!token) {
    router.push('/login')
    throw new Error('Not authenticated')
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers
    }
  })

  if (response.status === 401) {
    localStorage.removeItem('auth_token')
    router.push('/login')
    throw new Error('Session expired')
  }

  if (!response.ok) throw new Error('Request failed')
  return response.json().catch(() => ({}))
}

const loadUsers = async () => {
  try { users.value = await apiCall('/admin/users') } catch (e) { console.error(e) }
}

const loadTemplates = async () => {
  try { templates.value = await apiCall('/templates') } catch (e) { console.error(e) }
}

const loadProcesses = async () => {
  try { processes.value = await apiCall('/processes') } catch (e) { console.error(e) }
}

// Users
const editUser = (user) => {
  editingUser.value = { ...user }
  editForm.value = { username: user.username, role: user.role }
}

const cancelEdit = () => {
  editingUser.value = null
  editForm.value = { username: '', role: '' }
}

const saveUser = async (id) => {
  try {
    await apiCall(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(editForm.value)
    })
    await loadUsers()
    cancelEdit()
  } catch (e) {
    alert('Failed to update user')
  }
}

const deleteUser = async (id) => {
  if (!confirm('Delete this user?')) return
  try {
    await apiCall(`/admin/users/${id}`, { method: 'DELETE' })
    await loadUsers()
  } catch (e) {
    alert('Failed to delete user')
  }
}

// В AdminPanel.vue
const newTemplate = ref({ 
  name: '', 
  description: '', 
  tasks: [] 
})

const addTask = () => {
  newTemplate.value.tasks.push({
    title: '',
    for_role: 'worker',
    is_file_required: false,
    plan_hours: 0,
    depends_on: []
  })
}

const removeTask = (idx) => {
  newTemplate.value.tasks.splice(idx, 1)
}

const createTemplate = async () => {
  try {
    // Преобразуем индексы зависимостей в ID задач шаблона
    const tasksWithDeps = newTemplate.value.tasks.map((task, idx) => ({
      title: task.title,
      for_role: task.for_role,
      is_file_required: task.is_file_required,
      plan_hours: task.plan_hours,
      depends_on: task.depends_on.map(depIdx => depIdx + 1) // временные ID
    }))
    
    await apiCall('/templates', {
      method: 'POST',
      body: JSON.stringify({
        name: newTemplate.value.name,
        description: newTemplate.value.description,
        tasks: tasksWithDeps
      })
    })
    showAddTemplate.value = false
    newTemplate.value = { name: '', description: '', tasks: [] }
    await loadTemplates()
  } catch (e) {
    alert('Failed to create template: ' + e.message)
  }
}

const deleteTemplate = async (id) => {
  if (!confirm('Delete this template?')) return
  try {
    await apiCall(`/admin/templates/${id}`, { method: 'DELETE' })
    await loadTemplates()
  } catch (e) {
    alert('Failed to delete template')
  }
}

// Processes
const updateProcessStatus = async (proc) => {
  try {
    // Используем PATCH /tasks/ для обновления, или напрямую в БД
    // MVP: просто визуально меняем
  } catch (e) {
    alert('Failed to update status')
  }
}

const deleteProcess = async (id) => {
  if (!confirm('Delete this process?')) return
  try {
    await apiCall(`/admin/processes/${id}`, { method: 'DELETE' })
    await loadProcesses()
  } catch (e) {
    alert('Failed to delete process')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  try { return new Date(dateStr).toLocaleDateString() } catch { return dateStr }
}

onMounted(() => {
  loadUsers()
  loadTemplates()
  loadProcesses()
})
</script>

<style scoped>
.admin-panel {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

h2 { font-family: var(--font-1); margin-bottom: 1.5rem; }

.tabs { display: flex; gap: 8px; margin-bottom: 1.5rem; }

.tabs button {
  padding: 8px 20px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-btn);
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
}
.tabs button.active {
  background: var(--color-text);
  color: var(--color-background);
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}
.tab-header h3 { margin: 0; }

.btn-add {
  padding: 6px 16px;
  border: 2px solid var(--color-success);
  border-radius: var(--radius-btn);
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-success);
  transition: all 0.15s ease;
}
.btn-add:hover {
  background: var(--color-success);
  color: white;
}

.table-container { overflow-x: auto; }

table {
  width: 100%;
  border-collapse: collapse;
  font-family: var(--font-1);
}

th {
  text-align: left;
  padding: 12px;
  border-bottom: 2px solid var(--color-border);
}

td {
  padding: 10px 12px;
  border-bottom: 1px solid var(--color-muted-light);
}

.edit-input {
  padding: 4px 8px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: transparent;
  font-family: var(--font-1);
  color: var(--color-text);
}

.actions { display: flex; gap: 6px; }

.btn-edit, .btn-save, .btn-cancel, .btn-delete {
  padding: 4px 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 0.8rem;
  background: transparent;
  font-family: var(--font-1);
  color: var(--color-text);
}
.btn-save { border-color: var(--color-success); color: var(--color-success); }
.btn-delete { border-color: var(--color-danger); color: var(--color-danger); }
.btn-cancel { border-color: var(--color-muted); color: var(--color-muted); }

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
  max-width: 90vw;
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

/* Строки для темной темы */
td, th { color: var(--color-text); }
</style>