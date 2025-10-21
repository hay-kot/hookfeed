/**
 * Composable for generating mock feed message data for development/design
 * Based on DtosFeedMessage type from data-contracts.ts
 */

import type { DtosFeedMessage } from '~/lib/api/types/data-contracts'

// Helper to encode JSON as number[] (simulating backend behavior)
const encodeJSON = (obj: any): number[] => {
  const str = JSON.stringify(obj)
  return Array.from(new TextEncoder().encode(str))
}

// Helper to decode number[] back to JSON
export const decodeJSON = (encoded: number[]): any => {
  if (!encoded || encoded.length === 0) return null
  try {
    const str = new TextDecoder().decode(new Uint8Array(encoded))
    return JSON.parse(str)
  } catch {
    return null
  }
}

const levels: Array<'info' | 'warning' | 'error' | 'success' | 'debug'> = [
  'info',
  'warning',
  'error',
  'success',
  'debug',
]

const states: Array<'new' | 'acknowledged' | 'resolved' | 'archived'> = [
  'new',
  'acknowledged',
  'resolved',
  'archived',
]

const sampleTitles = [
  'Deployment completed successfully',
  'Database backup failed',
  'New user registration',
  'Payment received',
  'API rate limit exceeded',
  'Server health check warning',
  'Critical error in production',
  'Feature flag updated',
  'Scheduled task completed',
  'Security alert detected',
]

const sampleMessages = [
  'The deployment to production environment has completed successfully. All services are running normally.',
  'Automated database backup failed due to insufficient disk space. Please check server storage.',
  'A new user has registered on the platform. Email verification pending.',
  'Payment of $99.99 received from customer. Transaction ID: TXN-12345.',
  'API rate limit of 1000 requests per hour has been exceeded. Some requests were throttled.',
  'Server CPU usage has exceeded 80% for the last 5 minutes. Consider scaling.',
  'Uncaught exception in payment processing service. Immediate attention required.',
  'Feature flag "new-checkout-flow" has been enabled for 10% of users.',
  'Daily data aggregation task completed in 2.3 seconds.',
  'Suspicious login attempt detected from IP 192.168.1.100.',
]

const sampleMetadata = [
  { environment: 'production', region: 'us-east-1', version: '2.1.0' },
  { userId: 'user_123', action: 'login', ip: '192.168.1.100' },
  { orderId: 'ORD-98765', amount: 99.99, currency: 'USD' },
  { service: 'api-gateway', endpoint: '/api/v1/users', method: 'POST' },
  { timestamp: Date.now(), severity: 'high', affectedUsers: 150 },
]

const sampleRawRequests = [
  {
    event: 'deployment.completed',
    repository: 'my-app',
    branch: 'main',
    commit: 'abc123def456',
    user: 'john@example.com',
  },
  {
    type: 'error',
    error: 'ENOSPC: no space left on device',
    path: '/var/backups',
    timestamp: new Date().toISOString(),
  },
  {
    event: 'user.registered',
    userId: 'usr_abc123',
    email: 'newuser@example.com',
    plan: 'free',
  },
  {
    event: 'payment.received',
    transactionId: 'txn_xyz789',
    amount: 9999,
    currency: 'usd',
    customer: 'cus_123abc',
  },
  {
    event: 'rate_limit.exceeded',
    apiKey: 'sk_live_***xyz',
    requestCount: 1250,
    limit: 1000,
    window: '1h',
  },
]

const sampleRawHeaders = [
  {
    'content-type': 'application/json',
    'user-agent': 'GitHub-Hookshot/abc123',
    'x-github-event': 'push',
    'x-github-delivery': 'uuid-here',
  },
  {
    'content-type': 'application/json',
    'user-agent': 'Discord-Bot/1.0',
    'x-discord-signature': 'signature-here',
  },
  {
    'content-type': 'application/json',
    'user-agent': 'Stripe/1.0',
    'stripe-signature': 't=timestamp,v1=signature',
  },
  {
    'content-type': 'application/json',
    'user-agent': 'Custom-Service/2.1',
    'x-api-key': 'sk_***abc123',
  },
]

const sampleLogs = [
  ['Processing webhook from GitHub', 'Validated signature', 'Extracted event data', 'Stored in database'],
  ['Received backup request', 'Checking disk space', 'ERROR: Insufficient space', 'Backup aborted'],
  ['New user registration started', 'Email format validated', 'Password hashed', 'User created', 'Verification email sent'],
  ['Payment webhook received', 'Signature verified', 'Amount validated', 'Payment recorded', 'Receipt generated'],
  ['Rate limit check started', 'Request count: 1250', 'Limit: 1000', 'WARNING: Limit exceeded', 'Throttling applied'],
]

export const useMockData = () => {
  // Generate a single mock message
  const generateMessage = (index: number = 0, feedSlug: string = 'technology'): DtosFeedMessage => {
    const now = new Date()
    const receivedAt = new Date(now.getTime() - Math.random() * 7 * 24 * 60 * 60 * 1000) // Random time within last 7 days
    const processedAt = new Date(receivedAt.getTime() + Math.random() * 5000) // Processed 0-5 seconds after received

    const titleIndex = index % sampleTitles.length
    const level = levels[index % levels.length]
    const state = states[index % states.length]

    return {
      id: `msg-${Math.random().toString(36).substring(2, 11)}`,
      feedSlug,
      title: sampleTitles[titleIndex],
      message: sampleMessages[titleIndex],
      level,
      state,
      logs: sampleLogs[index % sampleLogs.length],
      metadata: encodeJSON(sampleMetadata[index % sampleMetadata.length]),
      rawRequest: encodeJSON(sampleRawRequests[index % sampleRawRequests.length]),
      rawHeaders: encodeJSON(sampleRawHeaders[index % sampleRawHeaders.length]),
      receivedAt: receivedAt.toISOString(),
      processedAt: processedAt.toISOString(),
      createdAt: receivedAt.toISOString(),
      updatedAt: receivedAt.toISOString(),
      stateChangedAt: state !== 'new' ? processedAt.toISOString() : undefined,
    }
  }

  // Generate multiple messages
  const generateMessages = (count: number = 10, feedSlug?: string): DtosFeedMessage[] => {
    return Array.from({ length: count }, (_, i) => generateMessage(i, feedSlug))
  }

  return {
    generateMessage,
    generateMessages,
    decodeJSON,
  }
}
