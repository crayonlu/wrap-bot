import { useTasks, useTriggerTask } from '../lib/hooks/useQuery'
import toast from 'react-hot-toast'
import { Clock, Play, Calendar } from 'lucide-react'

export default function TasksPage() {
  const { data: tasks, isLoading } = useTasks()
  const triggerTask = useTriggerTask()

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
        <h1>Tasks</h1>
        <p>Scheduled tasks and jobs</p>
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
