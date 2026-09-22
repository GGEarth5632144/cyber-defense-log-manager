import { useState, useEffect } from 'react'
import axios from 'axios'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'

export default function Dashboard() {
  const [stats, setStats] = useState<any>(null)
  const [timeline, setTimeline] = useState<any[]>([])
  const [error, setError] = useState<string | null>(null)
  
  const [timeRange, setTimeRange] = useState('24h')
  const [tenant, setTenant] = useState('all')
  const [isAdmin, setIsAdmin] = useState(false)
  
  useEffect(() => {
    // Check if user is admin
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        const u = JSON.parse(storedUser)
        if (u.role === 'admin') setIsAdmin(true)
      } catch (e) {}
    }
  }, [])

  useEffect(() => {
    setError(null)
    const params = { time: timeRange, tenant: tenant }
    
    axios.get('/api/stats', { params })
      .then(res => setStats(res.data))
      .catch(err => {
        console.error(err)
        setError(err.message)
      })
      
    axios.get('/api/timeline', { params })
      .then(res => setTimeline(res.data))
      .catch(console.error)
  }, [timeRange, tenant])

  if (error) return <div className="text-red-500">Failed to load dashboard: {error}. Are you logged in?</div>
  if (!stats) return <div className="text-slate-500 font-semibold animate-pulse">Loading Dashboard...</div>

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];
  
  const topIps = stats.top_ips || []
  const topUsers = stats.top_users || []
  const topEventTypes = stats.top_event_types || []
  const timelineData = timeline || []

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        
        <div className="flex space-x-4 bg-white p-2 rounded-lg shadow-sm border">
          {isAdmin && (
            <select 
              value={tenant} 
              onChange={e => setTenant(e.target.value)}
              className="bg-slate-50 border rounded p-1 text-sm"
            >
              <option value="all">All Tenants</option>
              <option value="demoA">demoA</option>
              <option value="demoB">demoB</option>
            </select>
          )}
          
          <select 
            value={timeRange} 
            onChange={e => setTimeRange(e.target.value)}
            className="bg-slate-50 border rounded p-1 text-sm"
          >
            <option value="1h">Last 1 Hour</option>
            <option value="24h">Last 24 Hours</option>
            <option value="7d">Last 7 Days</option>
          </select>
        </div>
      </div>
      
      <div className="grid grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow border border-slate-100 flex flex-col items-center justify-center">
          <div className="text-slate-500">Total Logs</div>
          <div className="text-4xl font-bold">{stats.total_logs || 0}</div>
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow border border-slate-100">
        <h2 className="text-xl font-semibold mb-4">Event Timeline</h2>
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={timelineData}>
              <XAxis dataKey="time" tickFormatter={(v) => v ? v.split(' ')[1] : ''} />
              <YAxis />
              <Tooltip />
              <Bar dataKey="count" fill="#3b82f6" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow border border-slate-100">
          <h2 className="text-xl font-semibold mb-4">Top Source IPs</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie data={topIps} dataKey="count" nameKey="name" cx="50%" cy="50%" outerRadius={80} label>
                  {topIps.map((_: any, index: number) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow border border-slate-100">
          <h2 className="text-xl font-semibold mb-4">Top Users</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={topUsers} layout="vertical">
                <XAxis type="number" />
                <YAxis dataKey="name" type="category" width={100} />
                <Tooltip />
                <Bar dataKey="count" fill="#10b981" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow border border-slate-100">
          <h2 className="text-xl font-semibold mb-4">Top Event Types</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={topEventTypes} layout="vertical">
                <XAxis type="number" />
                <YAxis dataKey="name" type="category" width={100} />
                <Tooltip />
                <Bar dataKey="count" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  )
}
