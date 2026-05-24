export const ErrorCode = {
  NoError: 0,
  NotFoundFrom: 1,
  NotFoundTo: 2,
  NotFoundRoute: 3,
  ServerError: 4,
} as const;
export type ErrorCodeType = (typeof ErrorCode)[keyof typeof ErrorCode];
