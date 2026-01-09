import { useEffect } from 'react'
import { useTasks, useTriggerTask } from '../lib/hooks/useQuery'
import { useWebSocket } from '../lib/hooks/useWebSocket'
import { useWebSocketStore } from '../stores/websocket'
import toast from 'react-hot-toast'
import { Clock, Play, Calendar, Wifi, WifiOff } from 'lucide-react'

export default function TasksPage() {
  const { data: initialTasks, isLoading } = useTasks()
  const triggerTask = useTriggerTask()
  const { connected } = useWebSocketStore()
  const tasks = useWebSocketStore((state) => state.tasks)
  const setTasks = useWebSocketStore((state) => state.setTasks)

  // Initialize tasks from API response
  useEffect(() => {
    if (initialTasks && initialTasks.length > 0) {
      setTasks(initialTasks)
    }
  }, [initialTasks, setTasks])

  // WebSocket integration
  useWebSocket({
    enabled: true,
    onMessage: (message) => {
      if (message.type === 'tasks') {
        const updatedTasks = message.data as any[]
        setTasks(updatedTasks)
      }
    },
  })

  const handleTrigger = async (id: string, name: string) => {
    try {
      await triggerTask.mutateAsync(id)
      toast.success(`Task "${name}" triggered successfully`)
    } catch (error) {
      toast.error(`Failed to trigger task "${name}"`)
    }
  }

  const getStatusClass = (status: string) => {
    switch (status) {
      case 'running':
        return 'logs__badge--info'
      case 'failed':
        return 'logs__badge--error'
      default:
        return 'logs__badge--debug'
    }
  }

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  return (
    <div className="tasks">
      <div className="tasks__header">
        <div>
          <h1>Tasks</h1>
          <p>Scheduled tasks and jobs</p>
        </div>
        <div className="tasks__connection-status">
          {connected ? (
            <span className="tasks__status tasks__status--connected">
              <Wifi size={16} />
              Connected
            </span>
          ) : (
            <span className="tasks__status tasks__status--disconnected">
              <WifiOff size={16} />
              Disconnected
            </span>
          )}
        </div>
      </div>

      <div className="plugins__grid">
        {tasks?.map((task) => (
          <div key={task.id} className="plugins__card">
            <div className="plugins__card-header">
              <h3 className="plugins__card-title">{task.name}</h3>
              <span className={`logs__badge ${getStatusClass(task.status)}`}>
                {task.status}
              </span>
            </div>
            <p className="plugins__card-description">{task.description || task.schedule}</p>
            <div className="dashboard__info-row">
              <span><Clock style={{display: 'inline', width: '1rem', height: '1rem', verticalAlign: 'middle', marginRight: '0.5rem'}} />{task.schedule}</span>
            </div>
            {task.next_run && (
              <div className="dashboard__info-row">
                <span><Calendar style={{display: 'inline', width: '1rem', height: '1rem', verticalAlign: 'middle', marginRight: '0.5rem'}} />{new Date(task.next_run).toLocaleString()}</span>
              </div>
            )}
            <button
              onClick={() => handleTrigger(task.id, task.name)}
              disabled={triggerTask.isPending}
              className="login-page__button"
              style={{marginTop: '1rem'}}
            >
              <Play style={{width: '1rem', height: '1rem'}} />
              Run Now
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}
