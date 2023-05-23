interface ApiError {
  status: number;
  message: string;
}

export async function sendPostRequest<T>(url: string, body: object): Promise<T> {
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      const errorData: ApiError = await response.json();
      throw new Error(errorData.message);
    }

    const data: T = await response.json();

    return data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}
