/* post-processed by ./scripts/process-types.go */
/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface Feed {
  adapters: string[];
  description: string;
  id: string;
  keys: string[];
  middleware: string[];
  name: string;
  retention: Retention;
}

export interface FeedMessage {
  createdAt: Date | string;
  feedSlug: string;
  id: string;
  logs: string[];
  message: string;
  metadata: number[];
  priority: number;
  processedAt: string;
  rawHeaders: number[];
  rawQueryParams: number[];
  rawRequest: number[];
  receivedAt: string;
  state: string;
  stateChangedAt: string;
  title: string;
  updatedAt: Date | string;
}

export interface FeedMessageBulkDelete {
  filter: FeedMessageDeleteFilter;
  messageIds: string[];
}

export interface FeedMessageBulkUpdateState {
  /** @minItems 1 */
  messageIds: string[];
  state: "new" | "acknowledged" | "resolved" | "archived";
}

export interface FeedMessageCreate {
  feedSlug: string;
  logs: string[];
  message: string;
  metadata: number[];
  /**
   * @min 1
   * @max 5
   */
  priority: number;
  processedAt: string;
  rawHeaders: number[];
  rawQueryParams: number[];
  rawRequest: number[];
  receivedAt: string;
  state: "new" | "acknowledged" | "resolved" | "archived";
  title: string;
}

export interface FeedMessageDeleteFilter {
  olderThan: string;
  /**
   * @min 1
   * @max 5
   */
  priority: number;
}

export interface FeedMessageUpdateState {
  state: "new" | "acknowledged" | "resolved" | "archived";
}

export interface PaginationResponseDtosFeedMessage {
  items: FeedMessage[];
  total: number;
}

export interface PasswordReset {
  /** @minLength 8 */
  password: string;
  token: string;
}

export interface PasswordResetRequest {
  email: string;
}

export interface Retention {
  maxAgeDays: number;
  maxCount: number;
}

export interface StatusResponse {
  build: string;
}

export interface User {
  createdAt: Date | string;
  email: string;
  id: string;
  subscriptionEndedDate: string;
  subscriptionStartDate: string;
  updatedAt: Date | string;
  username: string;
}

export interface UserAuthenticate {
  email: string;
  password: string;
}

export interface UserRegister {
  email: string;
  /**
   * @minLength 6
   * @maxLength 256
   */
  password: string;
  /**
   * @minLength 6
   * @maxLength 128
   */
  username: string;
}

export interface UserSession {
  expiresAt: string;
  token: string;
}

export interface UserUpdate {
  email: string;
  /**
   * @minLength 6
   * @maxLength 256
   */
  password: string;
  /**
   * @minLength 6
   * @maxLength 128
   */
  username: string;
}

export interface WebhookResponse {
  feedId: string;
  messageId: string;
  success: boolean;
}

export interface ServerErrorResp {
  data: any;
  message: string;
  requestId: string;
  statusCode: number;
}
