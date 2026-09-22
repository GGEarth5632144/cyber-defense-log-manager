import { useState, useEffect } from 'react'
import axios from 'axios'
import { format } from 'date-fns'

export default function Alerts() {
  const [alerts, setAlerts] = useState<any[]>([])

  useEffect(() => {
    fetchAlerts()
  }, [])

  const fetchAlerts = async () => {
    try {
      const res = await axios.get('/api/alerts')
      setAlerts(res.data)
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="space-y-4 border rounded-lg bg-white p-4 shadow-sm">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">System Alerts</h1>
        <button onClick={fetchAlerts} className="bg-slate-200 px-4 py-2 rounded hover:bg-slate-300">Refresh</button>
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Timestamp</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tenant</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Message</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Severity</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {alerts.map((alert: any) => (
              <tr key={alert.id} className="hover:bg-red-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {format(new Date(alert.timestamp), 'yyyy-MM-dd HH:mm:ss')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{alert.tenant}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600 font-semibold">{alert.message}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
                    High ({alert.severity})
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {alerts.length === 0 && <div className="text-center p-8 text-gray-500">No active alerts</div>}
      </div>
    </div>
  )
}
