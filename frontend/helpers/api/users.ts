import { sendPostRequest } from './base';
import { API_BASE_URL } from '../../constants';

export type CreateUserParams = {
  name: string;
  email: string;
  password: string;
};

export type User = {
  id: string;
  name: string;
  email: string;
  is_guest: boolean;
  created_at: string;
  updated_at: string;
};

export async function createUser(params: CreateUserParams): Promise<User> {
  const url = `${API_BASE_URL}/users`;
  return sendPostRequest<User>(url, params);
}
