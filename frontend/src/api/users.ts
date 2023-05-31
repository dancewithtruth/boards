import { sendGetRequest, sendPostRequest } from './base';
import { API_BASE_URL } from '@/constants';

export type CreateUserParams = {
  name: string;
  email?: string;
  password?: string;
  isGuest: boolean;
};

export type User = {
  id: string;
  name: string;
  email: string;
  is_guest: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateUserResponse = {
  user: User;
  jwt_token: string;
};

export async function createUser(params: CreateUserParams): Promise<CreateUserResponse> {
  const url = `${API_BASE_URL}/users`;
  return sendPostRequest<CreateUserResponse>(url, params);
}

export async function getUserByJwt(jwtToken: string): Promise<User> {
  const url = `${API_BASE_URL}/users/me`;
  return sendGetRequest<User>(url, jwtToken);
}
