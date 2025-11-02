import { useTasks, useTriggerTask } from '../lib/hooks/useQuery'
import { useNotificationStore } from '../stores/notification'
import { Clock, Play, Calendar } from 'lucide-react'

export default function TasksPage() {
  const { data: tasks, isLoading } = useTasks()
  const triggerTask = useTriggerTask()
  const addNotification = useNotificationStore((state) => state.addNotification)

  const handleTrigger = async (id: string, name: string) => {
    try {
      await triggerTask.mutateAsync(id)
      addNotification({
        type: 'success',
        message: `Task "${name}" triggered successfully`,
      })
    } catch (error) {
      addNotification({
        type: 'error',
        message: `Failed to trigger task "${name}"`,
      })
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running':
        return 'bg-green-100 text-green-800'
      case 'failed':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

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
        <h1 className="text-3xl font-bold text-[#8B7355]">Tasks</h1>
        <p className="text-gray-600 mt-1">Scheduled tasks and jobs</p>
      </div>

      <div className="space-y-4">
        {tasks?.map((task) => (
          <div
            key={task.id}
            className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] p-6 hover:shadow-md transition-shadow"
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-3">
                  <h3 className="font-semibold text-[#8B7355] text-lg">{task.name}</h3>
                  <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(task.status)}`}>
                    {task.status}
                  </span>
                </div>
                <p className="text-gray-600 text-sm mt-2">{task.description}</p>
                <div className="flex items-center gap-6 mt-4 text-sm text-gray-600">
                  <div className="flex items-center gap-2">
                    <Clock className="w-4 h-4" />
                    <span>Schedule: {task.schedule}</span>
                  </div>
                  {task.next_run && (
                    <div className="flex items-center gap-2">
                      <Calendar className="w-4 h-4" />
                      <span>Next: {new Date(task.next_run).toLocaleString()}</span>
                    </div>
                  )}
                </div>
              </div>
              <button
                onClick={() => handleTrigger(task.id, task.name)}
                disabled={triggerTask.isPending}
                className="ml-4 px-4 py-2 bg-[#8B7355] text-white rounded-lg hover:bg-[#6d5940] transition-colors flex items-center gap-2 disabled:opacity-50"
              >
                <Play className="w-4 h-4" />
                Run Now
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
