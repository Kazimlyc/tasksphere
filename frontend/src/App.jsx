import { useEffect, useMemo, useState } from 'react'
import './App.css'
import AuthCard from './components/AuthCard'
import CreateTaskModal from './components/CreateTaskModal'
import SidebarMenu from './components/SidebarMenu'
import StatusMessage from './components/StatusMessage'
import TaskCard from './components/TaskCard'

const sanitizeBase = (url) => url?.replace(/\/$/, '')
const buildApiCandidates = () => {
  const envUrl = sanitizeBase(import.meta.env.VITE_API_URL)
  const defaults = ['http://localhost:8081', 'http://localhost:8080']
  const unique = new Set([envUrl, ...defaults])
  return Array.from(unique).filter(Boolean)
}

const initialAuthState = { name: '', email: '', password: '' }

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
  const [taskStatus, setTaskStatus] = useState('todo')
  const [loadingTasks, setLoadingTasks] = useState(false)
  const [editingTaskId, setEditingTaskId] = useState(null)
  const [editTitle, setEditTitle] = useState('')
  const [editContent, setEditContent] = useState('')
  const [editStatus, setEditStatus] = useState('todo')
  const [profileName, setProfileName] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const apiBases = useMemo(() => buildApiCandidates(), [])

  useEffect(() => {
    if (token) {
      fetchTasks()
      fetchProfile()
    } else {
      setTasks([])
      setProfileName('')
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

  const fetchProfile = async () => {
    try {
      const data = await request('/me', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setProfileName(data?.name ?? data?.email ?? '')
    } catch (error) {
      if (error.userMessage?.toLowerCase().includes('unauthorized')) {
        clearToken()
        setToken('')
      }
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
    setTaskStatus('todo')
    setEditingTaskId(null)
    setEditTitle('')
    setEditContent('')
    setEditStatus('todo')
    setProfileName('')
    setMenuOpen(false)
    setCreateModalOpen(false)
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
          status: taskStatus,
        }),
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setTaskTitle('')
      setTaskContent('')
      setTaskStatus('todo')
      fetchTasks()
      return true
    } catch (error) {
      setStatus(error.userMessage ?? error.message)
    }
    return false
  }

  const handleCreateTaskAndClose = async (event) => {
    const created = await handleCreateTask(event)
    if (created) {
      setCreateModalOpen(false)
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

  const startEditTask = (task) => {
    setEditingTaskId(task.id)
    setEditTitle(task.title ?? '')
    setEditContent(task.content ?? '')
    setEditStatus(task.status ?? 'todo')
  }

  const cancelEditTask = () => {
    setEditingTaskId(null)
    setEditTitle('')
    setEditContent('')
    setEditStatus('todo')
  }

  const handleUpdateTask = async (event) => {
    event.preventDefault()
    if (!editingTaskId || !editTitle.trim()) return

    try {
      await request(`/tasks/${editingTaskId}`, {
        method: 'PUT',
        body: JSON.stringify({
          title: editTitle.trim(),
          content: editContent.trim(),
          status: editStatus,
        }),
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setTasks((prev) =>
        prev.map((task) =>
          task.id === editingTaskId
            ? {
                ...task,
                title: editTitle.trim(),
                content: editContent.trim(),
                status: editStatus,
              }
            : task,
        ),
      )
      cancelEditTask()
    } catch (error) {
      setStatus(error.userMessage ?? error.message)
    }
  }

  const handleMoveTask = async (taskId, nextStatus) => {
    const target = tasks.find((task) => task.id === taskId)
    if (!target || target.status === nextStatus) return

    const prevTasks = tasks
    setTasks((current) =>
      current.map((task) =>
        task.id === taskId ? { ...task, status: nextStatus } : task,
      ),
    )

    try {
      await request(`/tasks/${taskId}`, {
        method: 'PUT',
        body: JSON.stringify({
          title: target.title,
          content: target.content ?? '',
          status: nextStatus,
        }),
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
    } catch (error) {
      setTasks(prevTasks)
      setStatus(error.userMessage ?? error.message)
    }
  }

  return (
    <div className="page">
      <header>
        <div className="header-row">
          {token && (
            <button className="menu-button" onClick={() => setMenuOpen(true)}>
              Menü
            </button>
          )}
          <div>
            <h1>TaskSphere</h1>
            <p>Backend API ile tam entegre bir görev listesi</p>
          </div>
        </div>
      </header>

      <StatusMessage status={status} />

      {!token ? (
        <AuthCard
          authMode={authMode}
          onAuthModeChange={setAuthMode}
          authForm={authForm}
          onAuthFormChange={setAuthForm}
          onSubmit={handleAuth}
        />
      ) : (
        <TaskCard
          onOpenCreate={() => setCreateModalOpen(true)}
          onCreateForStatus={(status) => {
            setTaskStatus(status)
            setCreateModalOpen(true)
          }}
          onMoveTask={handleMoveTask}
          tasks={tasks}
          loadingTasks={loadingTasks}
          editingTaskId={editingTaskId}
          editTitle={editTitle}
          editContent={editContent}
          editStatus={editStatus}
          onEditTitleChange={setEditTitle}
          onEditContentChange={setEditContent}
          onEditStatusChange={setEditStatus}
          onStartEdit={startEditTask}
          onCancelEdit={cancelEditTask}
          onUpdate={handleUpdateTask}
          onDelete={handleDeleteTask}
        />
      )}

      {token && (
        <SidebarMenu
          isOpen={menuOpen}
          onClose={() => setMenuOpen(false)}
          userName={profileName || 'Kullanıcı'}
          onLogout={handleLogout}
        />
      )}

      {token && (
        <CreateTaskModal
          isOpen={createModalOpen}
          onClose={() => setCreateModalOpen(false)}
          title={taskTitle}
          content={taskContent}
          status={taskStatus}
          onTitleChange={setTaskTitle}
          onContentChange={setTaskContent}
          onStatusChange={setTaskStatus}
          onSubmit={handleCreateTaskAndClose}
        />
      )}
    </div>
  )
}

export default App
