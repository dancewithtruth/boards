import { User, sendPostRequest } from './index';
import { BASE_URL } from '@/constants';

export type LoginParams = {
  email: string;
  password: string;
};

type LoginResponse = {
  token: string;
};

export async function login(params: LoginParams): Promise<LoginResponse> {
  const url = `${BASE_URL}/auth/login`;
  return sendPostRequest<LoginResponse>(url, params);
}
