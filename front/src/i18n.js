import { reactive, computed, watch, ref } from 'vue'
import ru from './locales/ru.json'
import en from './locales/en.json'

const messages = { ru, en }
const currentLocale = ref(localStorage.getItem('locale') || 'ru')

export const t = (key) => {
  const keys = key.split('.')
  let translation = messages[currentLocale.value]
  
  for (const k of keys) {
    if (translation && translation[k]) {
      translation = translation[k]
    } else {
      return key
    }
  }
  return translation
}

export const setLocale = (locale) => {
  if (messages[locale]) {
    currentLocale.value = locale
    localStorage.setItem('locale', locale)
  }
}

export const useI18n = () => ({
  t,
  locale: computed(() => currentLocale.value),
  setLocale
})