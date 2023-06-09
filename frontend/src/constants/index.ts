export * from './settings';
export * from './websocket';

export const BASE_URL = 'http://backend-core:8080';
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
export const WS_URL = 'ws://localhost:8080/ws';
