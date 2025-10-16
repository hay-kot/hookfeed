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

export interface DtosFeed {
  adapters?: string[];
  description?: string;
  id?: string;
  keys?: string[];
  middleware?: string[];
  name?: string;
  retention?: DtosRetention;
}

export interface DtosFeedMessage {
  createdAt?: string;
  id?: string;
  updatedAt?: string;
}

export type DtosFeedMessageCreate = object;

export type DtosFeedMessageUpdate = object;

export interface DtosPaginationResponseDtosFeedMessage {
  items?: DtosFeedMessage[];
  total?: number;
}

export interface DtosPasswordReset {
  /** @minLength 8 */
  password: string;
  token: string;
}

export interface DtosPasswordResetRequest {
  email: string;
}

export interface DtosRetention {
  maxAgeDays?: number;
  maxCount?: number;
}

export interface DtosStatusResponse {
  build?: string;
}

export interface DtosUser {
  createdAt?: string;
  email?: string;
  id?: string;
  subscriptionEndedDate?: string;
  subscriptionStartDate?: string;
  updatedAt?: string;
  username?: string;
}

export interface DtosUserAuthenticate {
  email: string;
  password: string;
}

export interface DtosUserRegister {
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

export interface DtosUserSession {
  expiresAt?: string;
  token?: string;
}

export interface DtosUserUpdate {
  email?: string;
  /**
   * @minLength 6
   * @maxLength 256
   */
  password?: string;
  /**
   * @minLength 6
   * @maxLength 128
   */
  username?: string;
}

export interface DtosWebhookResponse {
  feedId?: string;
  messageId?: string;
  success?: boolean;
}

export interface ServerErrorResp {
  data?: any;
  message?: string;
  requestId?: string;
  statusCode?: number;
}
