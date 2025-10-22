/**
 * Composable for managing feed messages
 * Handles fetching, filtering, searching, and state updates
 */

import type { DtosFeedMessage, DtosPaginationResponseDtosFeedMessage } from '~/lib/api/types/data-contracts'

export interface MessageFilters {
  feedSlug?: string
  priority?: number
  state?: 'new' | 'acknowledged' | 'resolved' | 'archived'
  since?: string
  until?: string
  q?: string
  skip?: number
  limit?: number
}

export const useFeedMessages = () => {
  const messages = useState<DtosFeedMessage[]>('feedMessages', () => [])
  const totalMessages = useState<number>('totalMessages', () => 0)
  const isLoadingMessages = useState<boolean>('isLoadingMessages', () => false)
  const currentFilters = useState<MessageFilters>('messageFilters', () => ({
    limit: 50,
    skip: 0,
  }))

  // Fetch messages with optional filters
  const fetchMessages = async (filters?: MessageFilters) => {
    const { apiClient } = useApiClient()

    try {
      isLoadingMessages.value = true

      // Merge filters with current filters
      const queryParams = {
        ...currentFilters.value,
        ...filters,
      }

      // Update current filters
      currentFilters.value = queryParams

      const params = new URLSearchParams()
      if (queryParams.feedSlug) params.append('feedSlug', queryParams.feedSlug)
      if (queryParams.priority) params.append('priority', queryParams.priority.toString())
      if (queryParams.state) params.append('state', queryParams.state)
      if (queryParams.since) params.append('since', queryParams.since)
      if (queryParams.until) params.append('until', queryParams.until)
      if (queryParams.q) params.append('q', queryParams.q)
      if (queryParams.skip !== undefined) params.append('skip', queryParams.skip.toString())
      if (queryParams.limit) params.append('limit', queryParams.limit.toString())

      const response = await apiClient.get<DtosPaginationResponseDtosFeedMessage>(
        `/feed-messages?${params.toString()}`
      )

      messages.value = response.data?.items || []
      totalMessages.value = response.data?.total || 0
    } catch (error) {
      console.error('Failed to fetch messages:', error)
      messages.value = []
      totalMessages.value = 0
    } finally {
      isLoadingMessages.value = false
    }
  }

  // Update a single message state
  const updateMessageState = async (messageId: string, state: 'new' | 'acknowledged' | 'resolved' | 'archived') => {
    const { apiClient } = useApiClient()

    try {
      await apiClient.patch(`/feed-messages/${messageId}/state`, { state })

      // Update local state
      const message = messages.value.find(m => m.id === messageId)
      if (message) {
        message.state = state
        message.stateChangedAt = new Date().toISOString()
      }

      return { success: true }
    } catch (error) {
      console.error('Failed to update message state:', error)
      return { success: false, error }
    }
  }

  // Bulk update message states
  const bulkUpdateMessageStates = async (
    feedSlug: string,
    messageIds: string[],
    state: 'new' | 'acknowledged' | 'resolved' | 'archived'
  ) => {
    const { apiClient } = useApiClient()

    try {
      await apiClient.post(`/feeds/${feedSlug}/messages/bulk-state`, {
        messageIds,
        state,
      })

      // Update local state
      messages.value.forEach(message => {
        if (messageIds.includes(message.id || '')) {
          message.state = state
          message.stateChangedAt = new Date().toISOString()
        }
      })

      return { success: true }
    } catch (error) {
      console.error('Failed to bulk update message states:', error)
      return { success: false, error }
    }
  }

  // Delete a single message
  const deleteMessage = async (messageId: string) => {
    const { apiClient } = useApiClient()

    try {
      await apiClient.delete(`/feed-messages/${messageId}`)

      // Remove from local state
      messages.value = messages.value.filter(m => m.id !== messageId)
      totalMessages.value = Math.max(0, totalMessages.value - 1)

      return { success: true }
    } catch (error) {
      console.error('Failed to delete message:', error)
      return { success: false, error }
    }
  }

  // Get a single message by ID
  const getMessageById = async (messageId: string) => {
    const { apiClient } = useApiClient()

    try {
      const response = await apiClient.get<DtosFeedMessage>(`/feed-messages/${messageId}`)
      return response.data
    } catch (error) {
      console.error('Failed to fetch message:', error)
      return null
    }
  }

  // Clear filters and reset to defaults
  const clearFilters = () => {
    currentFilters.value = {
      limit: 50,
      skip: 0,
    }
  }

  // Load more messages (pagination)
  const loadMore = async () => {
    if (isLoadingMessages.value) return

    currentFilters.value.skip = (currentFilters.value.skip || 0) + (currentFilters.value.limit || 50)
    await fetchMessages()
  }

  return {
    messages: readonly(messages),
    totalMessages: readonly(totalMessages),
    isLoadingMessages: readonly(isLoadingMessages),
    currentFilters: readonly(currentFilters),
    fetchMessages,
    updateMessageState,
    bulkUpdateMessageStates,
    deleteMessage,
    getMessageById,
    clearFilters,
    loadMore,
  }
}
