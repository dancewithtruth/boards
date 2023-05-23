import { sendPostRequest } from './base';
import { API_BASE_URL } from '../../constants';

export type LoginParams = {
  email: string;
  password: string;
};

type LoginResponse = {
  token: string;
};

export async function login(params: LoginParams): Promise<LoginResponse> {
  const url = `${API_BASE_URL}/auth/login`;
  return sendPostRequest<LoginResponse>(url, params);
}
