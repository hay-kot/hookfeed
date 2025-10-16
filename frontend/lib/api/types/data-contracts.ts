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

export interface PasswordReset {
  /** @minLength 8 */
  password: string;
  token: string;
}

export interface PasswordResetRequest {
  email: string;
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
