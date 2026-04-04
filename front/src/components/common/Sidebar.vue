<template>
  <aside class="search-filter">
    <div class="search-box">
      <input 
        type="text" 
        v-model="localSearch" 
        @input="emitSearch"
        placeholder="Search processes..."
        class="search-input"
      />
      <span class="search-icon">🔍</span>
    </div>
    
    <div class="filters">
      <h4>Priority</h4>
      <label>
        <input type="checkbox" value="high" v-model="localFilters.priority" @change="emitFilters" />
        <span class="priority-dot high"></span> High
      </label>
      <label>
        <input type="checkbox" value="medium" v-model="localFilters.priority" @change="emitFilters" />
        <span class="priority-dot medium"></span> Medium
      </label>
      <label>
        <input type="checkbox" value="low" v-model="localFilters.priority" @change="emitFilters" />
        <span class="priority-dot low"></span> Low
      </label>
    </div>
    
    <div class="calendar">
      <div class="calendar-header">
        <button @click="prevMonth">←</button>
        <span>{{ monthName }} {{ year }}</span>
        <button @click="nextMonth">→</button>
      </div>
      <div class="calendar-grid">
        <div class="weekday">Mon</div>
        <div class="weekday">Tue</div>
        <div class="weekday">Wed</div>
        <div class="weekday">Thu</div>
        <div class="weekday">Fri</div>
        <div class="weekday">Sat</div>
        <div class="weekday">Sun</div>
        
        <div 
          v-for="day in calendarDays" 
          :key="day.date"
          class="calendar-day"
          :class="{ selected: day.selected, 'has-event': day.hasEvent }"
          @click="selectDay(day)"
        >
          {{ day.day }}
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['search', 'filter-change', 'date-change'])

const localSearch = ref('')
const localFilters = ref({ priority: [] })

const currentDate = ref(new Date())
const selectedDate = ref(null)

const monthName = computed(() => currentDate.value.toLocaleString('default', { month: 'long' }))
const year = computed(() => currentDate.value.getFullYear())

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  const firstDay = new Date(year, month, 1)
  const startWeekday = firstDay.getDay() || 7
  
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  
  const days = []
  
  for (let i = 1; i < startWeekday; i++) {
    days.push({ day: '', date: null, selected: false, hasEvent: false })
  }
  
  for (let i = 1; i <= daysInMonth; i++) {
    const date = new Date(year, month, i)
    const selected = selectedDate.value && selectedDate.value.toDateString() === date.toDateString()
    const hasEvent = [5, 12, 15, 20, 25].includes(i)
    days.push({ day: i, date: date, selected: selected, hasEvent: hasEvent })
  }
  
  return days
})

const emitSearch = () => {
  emit('search', localSearch.value)
}

const emitFilters = () => {
  emit('filter-change', localFilters.value)
}

const selectDay = (day) => {
  if (day.date) {
    selectedDate.value = day.date
    emit('date-change', day.date)
  }
}

const prevMonth = () => {
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
}

const nextMonth = () => {
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
}
</script>

<style scoped>
.search-filter {
  padding: 1.5rem;
  border: 2px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  background: var(--color-background);
}

.search-box {
  position: relative;
  margin-bottom: 1.5rem;
}

.search-input {
  width: 100%;
  padding: 12px 40px 12px 16px;
  border: 2px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  background: transparent;
  color: var(--color-text);
  font-family: var(--font-1);
}

.search-input:focus {
  outline: none;
  border-color: #4a90e2;
}

.search-icon {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  opacity: 0.6;
}

.filters {
  margin-bottom: 1.5rem;
}

.filters h4, .calendar h4 {
  margin-bottom: 1rem;
  font-family: var(--font-1);
}

.filters label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 10px 0;
  cursor: pointer;
}

.priority-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.priority-dot.high { background: #f44336; }
.priority-dot.medium { background: #ff9800; }
.priority-dot.low { background: #4caf50; }

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.calendar-header button {
  background: none;
  border: 2px solid var(--color-text);
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  cursor: pointer;
  font-size: 1rem;
  padding: 4px 12px;
  color: var(--color-text);
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.weekday {
  text-align: center;
  font-size: 0.7rem;
  font-weight: bold;
  padding: 6px;
}

.calendar-day {
  text-align: center;
  padding: 8px 4px;
  cursor: pointer;
  border-radius: 50%;
  font-size: 0.8rem;
  position: relative;
  transition: all 0.1s ease;
}

.calendar-day:hover {
  background: rgba(0, 0, 0, 0.05);
}

.calendar-day.selected {
  background: var(--color-text);
  color: var(--color-background);
}

.calendar-day.has-event::after {
  content: '';
  position: absolute;
  bottom: 2px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background: #4caf50;
  border-radius: 50%;
}
</style>