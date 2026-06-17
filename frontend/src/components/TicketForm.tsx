import { useState } from 'react'
import type { Ticket } from '../types'

interface Props {
  // If initial is provided the form is in edit mode; otherwise create mode.
  initial?: Ticket
  onSubmit: (data: { title: string; description: string; status: Ticket['status'] }) => Promise<void>
  onCancel: () => void
}

export default function TicketForm({ initial, onSubmit, onCancel }: Props) {
  const [title, setTitle] = useState(initial?.title ?? '')
  const [description, setDescription] = useState(initial?.description ?? '')
  const [status, setStatus] = useState<Ticket['status']>(initial?.status ?? 'open')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      await onSubmit({ title, description, status })
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="card" onSubmit={handleSubmit}>
      <h3>{initial ? 'Edit Ticket' : 'New Ticket'}</h3>

      {error && <p className="error">{error}</p>}

      <div className="form-group">
        <label>Title</label>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Short summary"
          required
        />
      </div>

      <div className="form-group">
        <label>Description</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Describe the issue"
          rows={3}
          required
        />
      </div>

      <div className="form-group">
        <label>Status</label>
        <select value={status} onChange={(e) => setStatus(e.target.value as Ticket['status'])}>
          <option value="open">Open</option>
          <option value="in_progress">In Progress</option>
          <option value="closed">Closed</option>
        </select>
      </div>

      <div className="form-actions">
        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? 'Saving…' : initial ? 'Update' : 'Create'}
        </button>
        <button type="button" className="btn btn-secondary" onClick={onCancel}>
          Cancel
        </button>
      </div>
    </form>
  )
}
