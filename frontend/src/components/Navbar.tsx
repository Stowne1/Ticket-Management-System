import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { logout } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <span className="navbar-brand">Ticket Manager</span>
      <button className="btn btn-secondary" onClick={handleLogout}>
        Log out
      </button>
    </nav>
  )
}
