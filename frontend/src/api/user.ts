import { BASE_URL } from '@/constants';
import { sendGetRequest, sendPostRequest } from '.';

export type User = {
  id: string;
  name: string;
  email: string;
  is_guest: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateUserParams = {
  name: string;
  email?: string;
  password?: string;
  isGuest: boolean;
};

export type CreateUserResponse = {
  user: User;
  jwt_token: string;
};

export type ListUsersByFuzzyEmailResponse = {
  result: User[];
};

export async function getUserByJwt(jwtToken: string): Promise<User> {
  const url = `${BASE_URL}/users/me`;
  return sendGetRequest<User>(url, jwtToken);
}

export async function createUser(params: CreateUserParams): Promise<CreateUserResponse> {
  const url = `${BASE_URL}/users`;
  return sendPostRequest<CreateUserResponse>(url, params);
}

export async function listUsersByFuzzyEmail(email: string): Promise<ListUsersByFuzzyEmailResponse> {
  const url = `${BASE_URL}/users/search?email=${email}`;
  return sendGetRequest<ListUsersByFuzzyEmailResponse>(url);
}
