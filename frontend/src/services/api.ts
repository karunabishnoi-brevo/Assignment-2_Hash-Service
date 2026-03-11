const BASE_URL = 'http://localhost:8080'

export interface HashResponse {
  input: string
  hash: string
}

export async function generateHash(input: string): Promise<HashResponse> {
  const response = await fetch(`${BASE_URL}/api/hash`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ input }),
  })
  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error || 'Failed to generate hash')
  }
  return response.json()
}
