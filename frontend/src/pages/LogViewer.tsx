import { useState, useEffect } from 'react'
import axios from 'axios'
import { format } from 'date-fns'

export default function LogViewer() {
  const [logs, setLogs] = useState<any[]>([])

  useEffect(() => {
    fetchLogs()
  }, [])

  const fetchLogs = async () => {
    try {
      const res = await axios.get('/api/logs')
      setLogs(res.data)
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="space-y-4 border rounded-lg bg-white p-4 shadow-sm">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Logs Explorer</h1>
        <button onClick={fetchLogs} className="bg-slate-200 px-4 py-2 rounded hover:bg-slate-300">Refresh</button>
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Timestamp</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tenant</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Source</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Event Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Source IP</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {logs.map((log: any) => (
              <tr key={log.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {format(new Date(log.timestamp), 'yyyy-MM-dd HH:mm:ss')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{log.tenant}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{log.source}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{log.event_type}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                    log.action === 'deny' ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'
                  }`}>
                    {log.action || 'unknown'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{log.src_ip}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {logs.length === 0 && <div className="text-center p-8 text-gray-500">No logs found</div>}
      </div>
    </div>
  )
}
