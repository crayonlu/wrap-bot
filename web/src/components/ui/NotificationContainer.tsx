import { useNotificationStore } from '../../stores/notification'
import { X, CheckCircle, AlertCircle, Info, AlertTriangle } from 'lucide-react'

const icons = {
  success: CheckCircle,
  error: AlertCircle,
  info: Info,
  warning: AlertTriangle,
}

export default function NotificationContainer() {
  const { notifications, removeNotification } = useNotificationStore()

  return (
    <div className="notification__container">
      {notifications.map((notification) => {
        const Icon = icons[notification.type]
        return (
          <div
            key={notification.id}
            className={`notification__item notification__item--${notification.type}`}
          >
            <div className="notification__icon">
              <Icon />
            </div>
            <p className="notification__message">{notification.message}</p>
            <button
              onClick={() => removeNotification(notification.id)}
              className="notification__close"
            >
              <X />
            </button>
          </div>
        )
      })}
    </div>
  )
}
