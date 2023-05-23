const API_BASE_URL = 'http://localhost:8080';

export type CreateUserParams = {
  name: string;
  email: string;
  password: string;
};

type CreateUserResponse = {
  id: string;
};

export const createUser = async (params: CreateUserParams): Promise<string> => {
  const url = `${API_BASE_URL}/users`;

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      throw new Error('Sign up failed');
    }

    const data: CreateUserResponse = await response.json();
    const { id } = data;

    return id;
  } catch (error) {
    console.error('Error signing up:', error);
    // Handle error
    // ...
    throw error;
  }
};
