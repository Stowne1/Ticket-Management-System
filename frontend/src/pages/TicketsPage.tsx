import { useEffect, useState } from 'react'
import { listTickets, createTicket, updateTicket, deleteTicket } from '../api'
import type { Ticket } from '../types'
import TicketForm from '../components/TicketForm'
import Navbar from '../components/Navbar'

const PAGE_SIZE = 10

// Badge colours by status
const statusClass: Record<Ticket['status'], string> = {
  open: 'badge badge-open',
  in_progress: 'badge badge-in-progress',
  closed: 'badge badge-closed',
}

const statusLabel: Record<Ticket['status'], string> = {
  open: 'Open',
  in_progress: 'In Progress',
  closed: 'Closed',
}

export default function TicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Which ticket is being edited (null = none)
  const [editing, setEditing] = useState<Ticket | null>(null)
  // Whether the create form is open
  const [creating, setCreating] = useState(false)

  const totalPages = Math.ceil(total / PAGE_SIZE)

  async function fetchTickets() {
    setLoading(true)
    setError('')
    try {
      const data = await listTickets(page, PAGE_SIZE)
      setTickets(data.tickets)
      setTotal(data.total)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to load tickets')
    } finally {
      setLoading(false)
    }
  }

  // Re-fetch whenever the page number changes.
  useEffect(() => {
    fetchTickets()
  }, [page])

  async function handleCreate(data: { title: string; description: string; status: Ticket['status'] }) {
    await createTicket(data)
    setCreating(false)
    fetchTickets()
  }

  async function handleUpdate(data: { title: string; description: string; status: Ticket['status'] }) {
    if (!editing) return
    await updateTicket(editing.id, data)
    setEditing(null)
    fetchTickets()
  }

  async function handleDelete(id: number) {
    if (!confirm('Delete this ticket?')) return
    await deleteTicket(id)
    fetchTickets()
  }

  return (
    <div>
      <Navbar />

      <main className="container">
        <div className="page-header">
          <h1>Tickets</h1>
          <button className="btn btn-primary" onClick={() => { setCreating(true); setEditing(null) }}>
            + New Ticket
          </button>
        </div>

        {creating && (
          <TicketForm
            onSubmit={handleCreate}
            onCancel={() => setCreating(false)}
          />
        )}

        {error && <p className="error">{error}</p>}

        {loading ? (
          <p className="muted">Loading…</p>
        ) : tickets.length === 0 ? (
          <p className="muted">No tickets yet. Create one above.</p>
        ) : (
          <div className="ticket-list">
            {tickets.map((ticket) => (
              <div key={ticket.id} className="card ticket-card">
                {editing?.id === ticket.id ? (
                  <TicketForm
                    initial={ticket}
                    onSubmit={handleUpdate}
                    onCancel={() => setEditing(null)}
                  />
                ) : (
                  <>
                    <div className="ticket-header">
                      <span className="ticket-title">{ticket.title}</span>
                      <span className={statusClass[ticket.status]}>{statusLabel[ticket.status]}</span>
                    </div>
                    <p className="ticket-description">{ticket.description}</p>
                    <div className="ticket-actions">
                      <button className="btn btn-secondary" onClick={() => { setEditing(ticket); setCreating(false) }}>
                        Edit
                      </button>
                      <button className="btn btn-danger" onClick={() => handleDelete(ticket.id)}>
                        Delete
                      </button>
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        )}

        {totalPages > 1 && (
          <div className="pagination">
            <button
              className="btn btn-secondary"
              onClick={() => setPage((p) => p - 1)}
              disabled={page === 1}
            >
              Previous
            </button>
            <span className="muted">Page {page} of {totalPages}</span>
            <button
              className="btn btn-secondary"
              onClick={() => setPage((p) => p + 1)}
              disabled={page === totalPages}
            >
              Next
            </button>
          </div>
        )}
      </main>
    </div>
  )
}
