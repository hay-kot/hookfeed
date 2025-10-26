/**
 * Composable for managing global application state
 * Fetches and manages feed data from the API
 */

import type { DtosFeed } from '~/lib/api/types/data-contracts'

export interface Feed {
  id: string
  name: string
  slug: string
  description: string
  category?: string
  unreadCount: number
}

export const useAppState = () => {
  const feeds = useState<Feed[]>('feeds', () => [])
  const isLoadingFeeds = useState<boolean>('isLoadingFeeds', () => false)

  // Fetch feeds from the API
  const fetchFeeds = async () => {
    const { apiClient } = useApiClient()

    try {
      isLoadingFeeds.value = true
      const response = await apiClient.get<DtosFeed[]>('/feeds')

      // Map backend DTOs to frontend Feed interface
      feeds.value = (response.data || []).map(dto => ({
        id: dto.id || '',
        name: dto.name || '',
        slug: dto.id || '', // Backend uses ID as slug
        description: dto.description || '',
        category: dto.category,
        unreadCount: 0, // Will be calculated separately
      }))
    } catch (error) {
      console.error('Failed to fetch feeds:', error)
      feeds.value = []
    } finally {
      isLoadingFeeds.value = false
    }
  }

  // Computed: Check if there are any new messages across all feeds
  const hasNewMessages = computed(() => {
    return feeds.value.some(feed => feed.unreadCount > 0)
  })

  // Computed: Total unread count
  const totalUnreadCount = computed(() => {
    return feeds.value.reduce((sum, feed) => sum + feed.unreadCount, 0)
  })

  // Get feed by ID
  const getFeedById = (id: string) => {
    return feeds.value.find(feed => feed.id === id)
  }

  // Get feed by slug
  const getFeedBySlug = (slug: string) => {
    return feeds.value.find(feed => feed.slug === slug)
  }

  // Mark feed as read (set unreadCount to 0)
  const markFeedAsRead = (feedId: string) => {
    const feed = getFeedById(feedId)
    if (feed) {
      feed.unreadCount = 0
    }
  }

  // Mark all feeds as read
  const markAllAsRead = () => {
    feeds.value.forEach(feed => {
      feed.unreadCount = 0
    })
  }

  return {
    feeds: readonly(feeds),
    isLoadingFeeds: readonly(isLoadingFeeds),
    hasNewMessages,
    totalUnreadCount,
    getFeedById,
    getFeedBySlug,
    markFeedAsRead,
    markAllAsRead,
    fetchFeeds,
  }
}
