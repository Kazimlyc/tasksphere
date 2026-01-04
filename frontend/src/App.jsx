import { useEffect, useMemo, useState } from 'react'
import './App.css'

const sanitizeBase = (url) => url?.replace(/\/$/, '')
const buildApiCandidates = () => {
  const envUrl = sanitizeBase(import.meta.env.VITE_API_URL)
  const defaults = ['http://localhost:8081', 'http://localhost:8080']
  const unique = new Set([envUrl, ...defaults])
  return Array.from(unique).filter(Boolean)
}

const initialAuthState = { email: '', password: '' }

const readToken = () => {
  try {
    return localStorage.getItem('token') ?? ''
  } catch {
    return ''
  }
}

const clearToken = () => {
  try {
    localStorage.removeItem('token')
  } catch {
    // No-op: storage might be blocked in some browsers.
  }
}

function App() {
  const [authMode, setAuthMode] = useState('login')
  const [authForm, setAuthForm] = useState(initialAuthState)
  const [token, setToken] = useState(() => readToken())
  const [status, setStatus] = useState('')
  const [tasks, setTasks] = useState([])
  const [taskTitle, setTaskTitle] = useState('')
  const [taskContent, setTaskContent] = useState('')
  const [loadingTasks, setLoadingTasks] = useState(false)
  const apiBases = useMemo(() => buildApiCandidates(), [])

  useEffect(() => {
    if (token) {
      fetchTasks()
    } else {
      setTasks([])
    }
  }, [token])

  const request = async (path, options = {}) => {
    let networkError = null

    for (const base of apiBases) {
      try {
        return await performRequest(base, path, options)
      } catch (error) {
        if (error.isNetworkError) {
          networkError = error
          continue
        }
        throw error
      }
    }

    throw (
      networkError ??
      Object.assign(new Error('API sunucusuna ulaşılamadı'), {
        userMessage: 'API sunucusuna ulaşılamadı. Backend ayakta mı?',
      })
    )
  }

  const performRequest = async (base, path, options) => {
    const headers = { 'Content-Type': 'application/json', ...(options.headers ?? {}) }
    const url = `${base}${path}`

    let response
    try {
      response = await fetch(url, {
        ...options,
        headers,
      })
    } catch (error) {
      const err = Object.assign(
        new Error(`API (${url}) erişilemedi: ${error.message}`),
        {
          isNetworkError: true,
          userMessage:
            'API sunucusuna ulaşılamadı. Lütfen backend adresi ile frontend/.env içindeki VITE_API_URL değerini eşleştirin.',
        },
      )
      throw err
    }

    let body = null
    try {
      body = await response.json()
    } catch {
      body = null
    }

    if (!response.ok) {
      throw Object.assign(
        new Error(body?.error ?? 'Beklenmeyen bir hata oluştu'),
        { userMessage: body?.error ?? 'Beklenmeyen bir hata oluştu' },
      )
    }
    return body
  }

  const fetchTasks = async () => {
    setLoadingTasks(true)
    try {
      const data = await request('/tasks', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setTasks(Array.isArray(data) ? data : [])
    } catch (error) {
      if (error.userMessage?.toLowerCase().includes('unauthorized')) {
        clearToken()
        setToken('')
      }
      setStatus(error.userMessage ?? error.message)
    } finally {
      setLoadingTasks(false)
    }
  }

  const handleAuth = async (event) => {
    event.preventDefault()
    setStatus('')

    const endpoint = authMode === 'login' ? '/login' : '/register'

    try {
      const data = await request(endpoint, {
        method: 'POST',
        body: JSON.stringify(authForm),
      })

      if (authMode === 'login') {
        try {
          localStorage.setItem('token', data.token)
        } catch {
          // Ignore storage errors; user stays logged in for this session.
        }
        setToken(data.token)
        setStatus('Giriş başarılı')
      } else {
        setStatus('Kayıt tamamlandı, şimdi giriş yapabilirsiniz')
        setAuthMode('login')
      }
      setAuthForm(initialAuthState)
    } catch (error) {
      setStatus(error.userMessage ?? error.message)
    }
  }

  const handleLogout = () => {
    clearToken()
    setToken('')
    setTaskTitle('')
    setTaskContent('')
    setStatus('Çıkış yapıldı')
  }

  const handleCreateTask = async (event) => {
    event.preventDefault()
    if (!taskTitle.trim()) return

    try {
      await request('/tasks', {
        method: 'POST',
        body: JSON.stringify({
          title: taskTitle.trim(),
          content: taskContent.trim(),
        }),
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setTaskTitle('')
      setTaskContent('')
      fetchTasks()
    } catch (error) {
      setStatus(error.userMessage ?? error.message)
    }
  }

  const handleDeleteTask = async (taskId) => {
    try {
      await request(`/tasks/${taskId}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setTasks((prev) => prev.filter((task) => task.id !== taskId))
    } catch (error) {
      setStatus(error.userMessage ?? error.message)
    }
  }

  return (
    <div className="page">
      <header>
        <h1>TaskSphere</h1>
        <p>Backend API ile tam entegre bir görev listesi</p>
      </header>

      {status && <div className="status">{status}</div>}

      {!token ? (
        <section className="card">
          <div className="tabs">
            <button
              className={authMode === 'login' ? 'active' : ''}
              onClick={() => setAuthMode('login')}
            >
              Giriş
            </button>
            <button
              className={authMode === 'register' ? 'active' : ''}
              onClick={() => setAuthMode('register')}
            >
              Kayıt
            </button>
          </div>

          <form className="form" onSubmit={handleAuth}>
            <label>
              E-posta
              <input
                type="email"
                value={authForm.email}
                required
                onChange={(event) => setAuthForm((prev) => ({ ...prev, email: event.target.value }))}
              />
            </label>
            <label>
              Şifre
              <input
                type="password"
                value={authForm.password}
                required
                minLength={6}
                onChange={(event) => setAuthForm((prev) => ({ ...prev, password: event.target.value }))}
              />
            </label>

            <button type="submit">{authMode === 'login' ? 'Giriş yap' : 'Kayıt ol'}</button>
          </form>
        </section>
      ) : (
        <section className="card">
          <div className="card-header">
            <h2>Görevlerin</h2>
            <button onClick={handleLogout}>Çıkış</button>
          </div>

          <form className="form inline" onSubmit={handleCreateTask}>
            <input
              placeholder="Yeni görev başlığı"
              value={taskTitle}
              onChange={(event) => setTaskTitle(event.target.value)}
            />
            <textarea
              placeholder="Kısa açıklama"
              rows={3}
              value={taskContent}
              onChange={(event) => setTaskContent(event.target.value)}
            />
            <button type="submit">Ekle</button>
          </form>

          {loadingTasks ? (
            <p>Yükleniyor...</p>
          ) : tasks.length === 0 ? (
            <p>Henüz görev yok</p>
          ) : (
            <ul className="task-list">
              {tasks.map((task) => (
                <li key={task.id}>
                  <div className="task-content">
                    <span className="task-title">{task.title}</span>
                    {task.content && <p>{task.content}</p>}
                  </div>
                  <button className="link" onClick={() => handleDeleteTask(task.id)}>
                    Sil
                  </button>
                </li>
              ))}
            </ul>
          )}
        </section>
      )}
    </div>
  )
}

export default App
