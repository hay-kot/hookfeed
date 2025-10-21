/**
 * Composable for managing global application state
 * Contains placeholder feed data and notification state
 */

export interface Feed {
  id: string
  name: string
  slug: string
  description: string
  unreadCount: number
}

export const useAppState = () => {
  // Placeholder feed data - will be replaced with API calls later
  const feeds = useState<Feed[]>('feeds', () => [
    {
      id: '1',
      name: 'Technology',
      slug: 'technology',
      description: 'Tech news and updates',
      unreadCount: 3,
    },
    {
      id: '2',
      name: 'News',
      slug: 'news',
      description: 'General news feed',
      unreadCount: 0,
    },
    {
      id: '3',
      name: 'Alerts',
      slug: 'alerts',
      description: 'Important system alerts',
      unreadCount: 5,
    },
    {
      id: '4',
      name: 'Development',
      slug: 'development',
      description: 'Development notifications',
      unreadCount: 1,
    },
  ])

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
    hasNewMessages,
    totalUnreadCount,
    getFeedById,
    getFeedBySlug,
    markFeedAsRead,
    markAllAsRead,
  }
}
