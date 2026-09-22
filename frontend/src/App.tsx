import { useState } from 'react'
import { BrowserRouter, Routes, Route, Navigate, useNavigate, Link } from 'react-router-dom'
import { Activity, ShieldAlert, List, LogOut } from 'lucide-react'
import Dashboard from './pages/Dashboard'
import LogViewer from './pages/LogViewer'
import Alerts from './pages/Alerts'
import axios from 'axios'

axios.defaults.baseURL = 'http://localhost:8080'
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

function Login({ setAuth }: { setAuth: (val: boolean) => void }) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const res = await axios.post('/login', { username, password })
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      axios.defaults.headers.common['Authorization'] = `Bearer ${res.data.token}`
      setAuth(true)
      navigate('/')
    } catch (err) {
      setError('Invalid credentials')
    }
  }

  return (
    <div className="flex h-screen items-center justify-center bg-slate-100">
      <div className="w-96 rounded-lg bg-white p-8 shadow-md">
        <h1 className="mb-6 text-2xl font-bold">LogManager Login</h1>
        {error && <div className="mb-4 text-red-500">{error}</div>}
        <form onSubmit={handleLogin}>
          <input className="mb-4 w-full rounded border p-2" type="text" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
          <input className="mb-4 w-full rounded border p-2" type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
          <button className="w-full rounded bg-blue-600 p-2 text-white hover:bg-blue-700" type="submit">Login</button>
        </form>
        <p className="mt-4 text-sm text-gray-500">welcome to log manager</p>
      </div>
    </div>
  )
}

function Layout({ children, setAuth }: { children: React.ReactNode, setAuth: (val: boolean) => void }) {
  const navigate = useNavigate()
  
  const handleLogout = () => {
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
    setAuth(false)
    navigate('/login')
  }

  return (
    <div className="flex h-screen bg-slate-50">
      <div className="w-64 bg-slate-900 text-white flex flex-col">
        <div className="p-4 text-xl font-bold flex items-center gap-2 border-b border-slate-700">
          <Activity /> LogManager
        </div>
        <nav className="flex-1 p-4 flex flex-col gap-2">
          <Link to="/" className="flex items-center gap-2 p-2 hover:bg-slate-800 rounded">
            <Activity size={18} /> Dashboard
          </Link>
          <Link to="/logs" className="flex items-center gap-2 p-2 hover:bg-slate-800 rounded">
            <List size={18} /> Logs Explorer
          </Link>
          <Link to="/alerts" className="flex items-center gap-2 p-2 hover:bg-slate-800 rounded">
            <ShieldAlert size={18} /> Alerts
          </Link>
        </nav>
        <div className="p-4 border-t border-slate-700">
          <button onClick={handleLogout} className="flex items-center gap-2 p-2 text-red-400 w-full hover:bg-slate-800 rounded">
            <LogOut size={18} /> Logout
          </button>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-8">
        {children}
      </div>
    </div>
  )
}

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('token'))

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login setAuth={setIsAuthenticated} />} />
        <Route path="/" element={isAuthenticated ? <Layout setAuth={setIsAuthenticated}><Dashboard /></Layout> : <Navigate to="/login" />} />
        <Route path="/logs" element={isAuthenticated ? <Layout setAuth={setIsAuthenticated}><LogViewer /></Layout> : <Navigate to="/login" />} />
        <Route path="/alerts" element={isAuthenticated ? <Layout setAuth={setIsAuthenticated}><Alerts /></Layout> : <Navigate to="/login" />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
