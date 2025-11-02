import { useStatus } from '../lib/hooks/useQuery'
import { Activity, MessageSquare, Send } from 'lucide-react'

export default function DashboardPage() {
  const { data: status, isLoading } = useStatus()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#8B7355]"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-[#8B7355]">Dashboard</h1>
        <p className="text-gray-600 mt-1">Bot status and statistics</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-[#FAF8F5] rounded-lg">
              <Activity className="w-6 h-6 text-[#8B7355]" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Status</p>
              <p className="text-2xl font-bold text-[#8B7355]">
                {status?.online ? 'Online' : 'Offline'}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-[#FAF8F5] rounded-lg">
              <MessageSquare className="w-6 h-6 text-[#8B7355]" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Messages Received</p>
              <p className="text-2xl font-bold text-[#8B7355]">
                {status?.stats.message_received || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-[#FAF8F5] rounded-lg">
              <Send className="w-6 h-6 text-[#8B7355]" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Messages Sent</p>
              <p className="text-2xl font-bold text-[#8B7355]">
                {status?.stats.message_sent || 0}
              </p>
            </div>
          </div>
        </div>
      </div>

      {status?.account && (
        <div className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] p-6">
          <h2 className="text-xl font-semibold text-[#8B7355] mb-4">Account Info</h2>
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-600">Nickname:</span>
              <span className="font-medium text-[#8B7355]">{status.account.nickname}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">User ID:</span>
              <span className="font-medium text-[#8B7355]">{status.account.user_id}</span>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
