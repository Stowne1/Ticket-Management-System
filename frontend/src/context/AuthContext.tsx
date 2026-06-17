import { createContext, useContext, useState, type ReactNode } from 'react'

interface AuthContextType {
  token: string | null
  login: (token: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  // Initialise from localStorage so the user stays logged in on page refresh.
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'))

  function login(newToken: string) {
    localStorage.setItem('token', newToken)
    setToken(newToken)
  }

  function logout() {
    localStorage.removeItem('token')
    setToken(null)
  }

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

// useAuth is a convenience hook so components don't have to import useContext
// and AuthContext separately — they just call useAuth().
export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider')
  return ctx
}
