import type { Ticket, TicketListResponse } from './types'

// All requests go through /api which Vite proxies to http://localhost:8080.
const BASE = '/api'

// Reads the JWT from localStorage and returns the Authorization header.
// Every protected endpoint requires this header.
function authHeader(): HeadersInit {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...authHeader(),
      ...options.headers,
    },
  })

  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error ?? `Request failed: ${res.status}`)
  }

  // 204 No Content has no body to parse.
  if (res.status === 204) return undefined as T
  return res.json()
}

// Auth
export const register = (email: string, password: string) =>
  request('/register', { method: 'POST', body: JSON.stringify({ email, password }) })

export const login = (email: string, password: string) =>
  request<{ token: string }>('/login', { method: 'POST', body: JSON.stringify({ email, password }) })

// Tickets
export const listTickets = (page: number, limit: number) =>
  request<TicketListResponse>(`/tickets?page=${page}&limit=${limit}`)

export const getTicket = (id: number) =>
  request<Ticket>(`/tickets/${id}`)

export const createTicket = (data: Omit<Ticket, 'id' | 'created_at' | 'updated_at'>) =>
  request<{ id: number; message: string }>('/tickets', { method: 'POST', body: JSON.stringify(data) })

export const updateTicket = (id: number, data: Omit<Ticket, 'id' | 'created_at' | 'updated_at'>) =>
  request(`/tickets/${id}`, { method: 'PUT', body: JSON.stringify(data) })

export const deleteTicket = (id: number) =>
  request(`/tickets/${id}`, { method: 'DELETE' })
