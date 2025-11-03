import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { statusAPI, pluginsAPI, tasksAPI, configAPI, logsAPI } from '../api'

export const useStatus = () => {
  return useQuery({
    queryKey: ['status'],
    queryFn: async () => {
      const response = await statusAPI.get()
      return response.data
    },
    refetchOnWindowFocus: false,
  })
}

export const usePlugins = () => {
  return useQuery({
    queryKey: ['plugins'],
    queryFn: async () => {
      const response = await pluginsAPI.list()
      return response.data
    },
  })
}

export const useTogglePlugin = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (name: string) => pluginsAPI.toggle(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plugins'] })
    },
  })
}

export const useTasks = () => {
  return useQuery({
    queryKey: ['tasks'],
    queryFn: async () => {
      const response = await tasksAPI.list()
      return response.data
    },
  })
}

export const useTriggerTask = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => tasksAPI.trigger(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export const useConfig = () => {
  return useQuery({
    queryKey: ['config'],
    queryFn: async () => {
      const response = await configAPI.list()
      return response.data
    },
  })
}

export const useUpdateConfig = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: configAPI.update,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['config'] })
    },
  })
}

export const useLogs = (level?: string, limit?: number) => {
  return useQuery({
    queryKey: ['logs', level, limit],
    queryFn: async () => {
      const response = await logsAPI.list({ level, limit })
      return response.data
    },
    refetchInterval: 10000,
  })
}
