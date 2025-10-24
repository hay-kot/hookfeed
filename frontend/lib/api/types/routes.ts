export type RCPRoute =
  | `/hooks/${string}/`
  | `/feed-messages/`
  | `/feed-messages/${string}/`
  | `/feed-messages/${string}/state/`
  | `/feeds/`
  | `/feeds/${string}/messages/bulk-delete/`
  | `/feeds/${string}/messages/bulk-state/`
  | `/info/`
  | `/users/login/`
  | `/users/register/`
  | `/users/request-password-reset/`
  | `/users/reset-password/`
  | `/users/self/`
  | `/${string}/`;
