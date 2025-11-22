import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import BacklogView from './views/BacklogView.vue'
import PlayingView from './views/PlayingView.vue'
import HistoryView from './views/HistoryView.vue'
import CalendarView from './views/CalendarView.vue'

const routes = [
  { path: '/', redirect: '/playing' },
  { path: '/backlog', component: BacklogView },
  { path: '/playing', component: PlayingView },
  { path: '/history', component: HistoryView },
  { path: '/calendar', component: CalendarView },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.mount('#app')
