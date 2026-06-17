export interface Ticket {
  id: number
  title: string
  description: string
  status: 'open' | 'in_progress' | 'closed'
  created_at: string
  updated_at: string
}

export interface TicketListResponse {
  tickets: Ticket[]
  total: number
  page: number
  limit: number
}
