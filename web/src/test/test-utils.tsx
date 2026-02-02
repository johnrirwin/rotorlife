import React, { ReactElement } from 'react'
import { render, RenderOptions } from '@testing-library/react'
import { AuthProvider } from '../contexts/AuthContext'

// Custom render that wraps components with providers
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  initialAuthState?: {
    isAuthenticated: boolean
    user: { id: string; email: string; displayName: string } | null
  }
}

function AllTheProviders({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      {children}
    </AuthProvider>
  )
}

const customRender = (
  ui: ReactElement,
  options?: CustomRenderOptions
) => render(ui, { wrapper: AllTheProviders, ...options })

// Re-export everything from testing-library
export * from '@testing-library/react'
export { customRender as render }
