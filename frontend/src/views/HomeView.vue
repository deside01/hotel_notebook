<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const records = ref([])
const search = ref('')
const error = ref('')
const loading = ref(false)
const saving = ref(false)
const editingId = ref(null)
const emptyForm = () => ({ full_name: '', phone: '', room: '', check_in: '', check_out: '', note: '' })
const form = ref(emptyForm())

const visibleRecords = computed(() => {
  const query = search.value.trim().toLowerCase()
  return records.value.filter((record) =>
    [record.full_name, record.phone, record.room].some((value) => value.toLowerCase().includes(query)),
  )
})

async function loadRecords() {
  loading.value = true
  error.value = ''
  try {
    records.value = (await api('/guest_records')) || []
  } catch (err) {
    if (err.status === 401) router.replace('/')
    else error.value = err.message
  } finally {
    loading.value = false
  }
}

function edit(record) {
  editingId.value = record.id
  form.value = {
    full_name: record.full_name,
    phone: record.phone,
    room: record.room,
    check_in: record.check_in,
    check_out: record.check_out,
    note: record.note,
  }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function cancelEdit() {
  editingId.value = null
  form.value = emptyForm()
  error.value = ''
}

async function save() {
  error.value = ''
  if (form.value.check_out <= form.value.check_in) {
    error.value = 'Дата выезда должна быть позже даты заезда'
    return
  }
  saving.value = true
  try {
    await api(editingId.value ? `/guest_records/${editingId.value}` : '/guest_records', {
      method: editingId.value ? 'PATCH' : 'POST',
      body: JSON.stringify(form.value),
    })
    cancelEdit()
    await loadRecords()
  } catch (err) {
    error.value = err.message
  } finally {
    saving.value = false
  }
}

async function remove(record) {
  if (!confirm(`Удалить запись «${record.full_name}»?`)) return
  error.value = ''
  try {
    await api(`/guest_records/${record.id}`, { method: 'DELETE' })
    if (editingId.value === record.id) cancelEdit()
    await loadRecords()
  } catch (err) {
    error.value = err.message
  }
}

onMounted(loadRecords)
</script>

<template>
  <main class="container">
    <header class="page-header">
      <div>
        <h1>Журнал гостей</h1>
        <p class="muted">Записи о заселении</p>
      </div>
    </header>

    <section class="panel">
      <h2>{{ editingId ? 'Редактировать запись' : 'Новая запись' }}</h2>
      <form class="record-form" @submit.prevent="save">
        <label>ФИО <input v-model.trim="form.full_name" required /></label>
        <label>Телефон <input v-model.trim="form.phone" type="tel" maxlength="15" required /></label>
        <label>Номер комнаты <input v-model.trim="form.room" required /></label>
        <label>Дата заезда <input v-model="form.check_in" type="date" required /></label>
        <label>Дата выезда <input v-model="form.check_out" type="date" required /></label>
        <label class="wide">Примечание <textarea v-model="form.note" rows="2" /></label>
        <div class="actions wide">
          <button type="submit" :disabled="saving">{{ saving ? 'Сохранение…' : editingId ? 'Сохранить' : 'Добавить' }}</button>
          <button v-if="editingId" type="button" class="secondary" @click="cancelEdit">Отмена</button>
        </div>
      </form>
    </section>

    <p v-if="error" class="error" role="alert">{{ error }}</p>

    <section class="panel list-panel">
      <div class="list-header">
        <h2>Записи</h2>
        <input v-model="search" type="search" placeholder="Поиск по ФИО, телефону, комнате" aria-label="Поиск записей" />
      </div>
      <p v-if="loading" class="muted">Загрузка…</p>
      <p v-else-if="!visibleRecords.length" class="muted">{{ search ? 'Ничего не найдено' : 'Записей пока нет' }}</p>
      <div v-else class="table-wrap">
        <table>
          <thead><tr><th>Гость</th><th>Телефон</th><th>Комната</th><th>Заезд</th><th>Выезд</th><th>Примечание</th><th>Действия</th></tr></thead>
          <tbody>
            <tr v-for="record in visibleRecords" :key="record.id">
              <td>{{ record.full_name }}</td>
              <td>{{ record.phone }}</td>
              <td>{{ record.room }}</td>
              <td>{{ record.check_in }}</td>
              <td>{{ record.check_out }}</td>
              <td>{{ record.note || '—' }}</td>
              <td class="row-actions">
                <button type="button" class="link-button" @click="edit(record)">Изменить</button>
                <button type="button" class="link-button danger" @click="remove(record)">Удалить</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </main>
</template>
