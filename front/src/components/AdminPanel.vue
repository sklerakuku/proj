<template>
  <div class="admin-panel">
    <h2>Admin Panel</h2>
    
    <!-- Табы -->
    <div class="tabs">
      <button @click="activeTab = 'users'" :class="{ active: activeTab === 'users' }">Users</button>
      <button @click="activeTab = 'templates'" :class="{ active: activeTab === 'templates' }">Templates</button>
      <button @click="activeTab = 'processes'" :class="{ active: activeTab === 'processes' }">Processes</button>
    </div>

    <!-- Таблица пользователей -->
    <div v-if="activeTab === 'users'" class="table-container">
      <h3>Users ({{ users.length }})</h3>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Role</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
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

    <!-- Таблица шаблонов -->
    <div v-if="activeTab === 'templates'" class="table-container">
      <h3>Templates ({{ templates.length }})</h3>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Description</th>
            <th>Tasks</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="template in templates" :key="template.id">
            <td>{{ template.id }}</td>
            <td>{{ template.name }}</td>
            <td>{{ template.description }}</td>
            <td>{{ template.tasks?.length || 0 }}</td>
            <td>
              <button @click="deleteTemplate(template.id)" class="btn-delete">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Таблица процессов -->
    <div v-if="activeTab === 'processes'" class="table-container">
      <h3>Processes ({{ processes.length }})</h3>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Status</th>
            <th>Tasks</th>
            <th>Started</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="process in processes" :key="process.id">
            <td>{{ process.id }}</td>
            <td>{{ process.title }}</td>
            <td>
              <span class="status-badge" :class="'status-' + process.status">{{ process.status }}</span>
            </td>
            <td>{{ process.tasks?.length || 0 }}</td>
            <td>{{ formatDate(process.started_at) }}</td>
            <td>
              <button @click="deleteProcess(process.id)" class="btn-delete">Delete</button>
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
  return response.json()
}

// Загрузка данных
const loadUsers = async () => {
  try { users.value = await apiCall('/admin/users') } catch (e) { console.error(e) }
}

const loadTemplates = async () => {
  try { templates.value = await apiCall('/templates') } catch (e) { console.error(e) }
}

const loadProcesses = async () => {
  try { processes.value = await apiCall('/processes') } catch (e) { console.error(e) }
}

// Пользователи
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

// Шаблоны
const deleteTemplate = async (id) => {
  if (!confirm('Delete this template?')) return
  try {
    await apiCall(`/admin/templates/${id}`, { method: 'DELETE' })
    await loadTemplates()
  } catch (e) {
    alert('Failed to delete template')
  }
}

// Процессы
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

h2 {
  font-family: var(--font-1);
  margin-bottom: 1.5rem;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 1.5rem;
}

.tabs button {
  padding: 8px 20px;
  border: 2px solid var(--color-text);
  border-radius: 255px 150px 225px 150px/150px 225px 150px 255px;
  background: transparent;
  cursor: pointer;
  font-family: var(--font-1);
  color: var(--color-text);
}

.tabs button.active {
  background: var(--color-text);
  color: var(--color-background);
}

.table-container {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-family: var(--font-1);
}

th {
  text-align: left;
  padding: 12px;
  border-bottom: 2px solid var(--color-text);
}

td {
  padding: 10px 12px;
  border-bottom: 1px solid rgba(0,0,0,0.1);
}

.edit-input {
  padding: 4px 8px;
  border: 2px solid var(--color-text);
  border-radius: 10px;
  background: transparent;
  font-family: var(--font-1);
}

.actions {
  display: flex;
  gap: 6px;
}

.btn-edit, .btn-save, .btn-cancel, .btn-delete {
  padding: 4px 12px;
  border: 1px solid var(--color-text);
  border-radius: 10px;
  cursor: pointer;
  font-size: 0.8rem;
  background: transparent;
}

.btn-save { border-color: #4caf50; color: #4caf50; }
.btn-delete { border-color: #f44336; color: #f44336; }
.btn-cancel { border-color: #999; color: #999; }

.status-badge {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 0.8rem;
}

.status-draft { background: #e0e0e0; }
.status-in_progress { background: #bbdefb; }
.status-done { background: #c8e6c9; }
</style>