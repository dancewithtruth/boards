export * from './settings';
export * from './websocket';

export const BASE_URL = typeof window === 'undefined' ? process.env.BACKEND_URL : process.env.NEXT_PUBLIC_BACKEND_URL;
export const COOKIE_NAME_JWT_TOKEN = 'jwt_token';
export const POST_WIDTH = 275;
export const DEFAULT_POST_HEIGHT = 114;
export const BOARD_SPACE_ADD = 150;
export const POST_COLORS: { [key: string]: string } = {
  LIGHT_PINK: '#F5E6E8',
  LIGHT_GREEN: '#E7ECD9',
  LIGHT_LAVENDER: '#E5E1F1',
  LIGHT_PEACH: '#FCE6C9',
  LIGHT_AQUA: '#D8E2DC',
};
export const ENV = {
  DEVELOPMENT: 'development',
  PRODUCTION: 'production',
};
export const WS_URL = process.env.NEXT_PUBLIC_WS_URL || `ws://localhost:8080/ws`
export const INVITE_STATUS = {
  PENDING: 'PENDING',
  ACCEPTED: 'ACCEPTED',
  IGNORED: 'IGNORED',
  CANCELLED: 'CANCELLED',
};
export const MEMBERSHIP_ROLES = {
  ADMIN: 'ADMIN',
  MEMBER: 'MEMBER',
};
