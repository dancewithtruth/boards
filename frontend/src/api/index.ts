export * from './auth';
export * from './user';
export * from './board';

export class APIError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = 'APIError';
  }
}

export async function sendPostRequest<T>(url: string, body: object, authToken?: string): Promise<T> {
  try {
    const headers: Record<string, string> = {};

    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...headers,
      },
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      const errorData: APIErrorResponse = await response.json();
      throw new APIError(errorData.message, errorData.status);
    }

    const data: T = await response.json();

    return data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}

export async function sendGetRequest<T>(url: string, authToken?: string): Promise<T> {
  try {
    const headers: Record<string, string> = {};

    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    const response = await fetch(url, {
      method: 'GET',
      headers: headers,
      cache: 'no-store',
    });

    if (!response.ok) {
      const errorData: APIErrorResponse = await response.json();
      throw new APIError(errorData.message, errorData.status);
    }

    const data: T = await response.json();

    return data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}
