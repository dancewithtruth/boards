export * from './settings';
export * from './websocket';

export const BASE_URL = typeof window === 'undefined' ? process.env.BACKEND_URL : process.env.NEXT_PUBLIC_BACKEND_URL;
export const COOKIE_NAME_JWT_TOKEN = 'jwt_token';
export const POST_WIDTH = 275;
export const POST_HEIGHT = 100;
export const BOARD_SPACE_ADD = 150;
export const POST_COLORS: { [key: string]: string } = {
  LIGHT_PINK: '#F5E6E8',
  LIGHT_GREEN: '#E7ECD9',
  LIGHT_LAVENDER: '#E5E1F1',
  LIGHT_PEACH: '#FCE6C9',
  LIGHT_AQUA: '#D8E2DC',
};
// TODO: Make env var more robust
export const WS_URL = process.env.ENV === 'development' ? 'ws://localhost:8080/ws' : 'ws://api.useboards.com/ws';
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
