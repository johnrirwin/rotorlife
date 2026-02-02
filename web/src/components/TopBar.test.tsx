import { describe, it, expect, vi } from 'vitest'
import { screen, fireEvent } from '@testing-library/react'
import { render } from '../test/test-utils'
import { TopBar } from './TopBar'

describe('TopBar', () => {
  const defaultProps = {
    query: '',
    onQueryChange: vi.fn(),
    fromDate: '2024-01-01',
    toDate: '2024-01-31',
    onFromDateChange: vi.fn(),
    onToDateChange: vi.fn(),
    sort: 'newest' as const,
    onSortChange: vi.fn(),
    onRefresh: vi.fn(),
    isRefreshing: false,
    totalCount: 42,
  }

  it('renders search input', () => {
    render(<TopBar {...defaultProps} />)
    expect(screen.getByPlaceholderText('Search news...')).toBeInTheDocument()
  })

  it('displays query value in search input', () => {
    render(<TopBar {...defaultProps} query="drone racing" />)
    const input = screen.getByPlaceholderText('Search news...') as HTMLInputElement
    expect(input.value).toBe('drone racing')
  })

  it('calls onQueryChange when typing in search', () => {
    const onQueryChange = vi.fn()
    render(<TopBar {...defaultProps} onQueryChange={onQueryChange} />)

    const input = screen.getByPlaceholderText('Search news...')
    fireEvent.change(input, { target: { value: 'fpv' } })

    expect(onQueryChange).toHaveBeenCalledWith('fpv')
  })

  it('shows clear button when query is not empty', () => {
    const { container } = render(<TopBar {...defaultProps} query="test" />)

    // The clear button should be visible
    const buttons = container.querySelectorAll('button')
    const clearButton = Array.from(buttons).find(btn => 
      btn.querySelector('svg path[d*="M6 18L18 6"]')
    )
    expect(clearButton).toBeInTheDocument()
  })

  it('does not show clear button in search input when query is empty', () => {
    render(<TopBar {...defaultProps} query="" />)
    
    // The search input should not have the clear button
    // There may be other close icons (like clear dates), so we check within the search container
    const searchInput = screen.getByPlaceholderText('Search news...')
    const searchContainer = searchInput.closest('.relative')
    const clearButton = searchContainer?.querySelector('button')
    expect(clearButton).not.toBeInTheDocument()
  })

  it('calls onQueryChange with empty string when clear clicked', () => {
    const onQueryChange = vi.fn()
    const { container } = render(<TopBar {...defaultProps} query="test" onQueryChange={onQueryChange} />)

    const buttons = container.querySelectorAll('button')
    const clearButton = Array.from(buttons).find(btn => 
      btn.querySelector('svg path[d*="M6 18L18 6"]')
    )
    if (clearButton) {
      fireEvent.click(clearButton)
      expect(onQueryChange).toHaveBeenCalledWith('')
    }
  })

  it('renders date inputs', () => {
    render(<TopBar {...defaultProps} />)
    expect(screen.getByText('From:')).toBeInTheDocument()
    expect(screen.getByText('To:')).toBeInTheDocument()
  })

  it('calls onFromDateChange when from date changes', () => {
    const onFromDateChange = vi.fn()
    render(<TopBar {...defaultProps} onFromDateChange={onFromDateChange} />)

    const fromInput = screen.getByDisplayValue('2024-01-01')
    fireEvent.change(fromInput, { target: { value: '2024-02-01' } })

    expect(onFromDateChange).toHaveBeenCalledWith('2024-02-01')
  })

  it('calls onToDateChange when to date changes', () => {
    const onToDateChange = vi.fn()
    render(<TopBar {...defaultProps} onToDateChange={onToDateChange} />)

    const toInput = screen.getByDisplayValue('2024-01-31')
    fireEvent.change(toInput, { target: { value: '2024-02-28' } })

    expect(onToDateChange).toHaveBeenCalledWith('2024-02-28')
  })

  it('renders sort buttons', () => {
    render(<TopBar {...defaultProps} />)
    expect(screen.getByText('Newest')).toBeInTheDocument()
    expect(screen.getByText('Top Score')).toBeInTheDocument()
  })

  it('highlights active sort button', () => {
    render(<TopBar {...defaultProps} sort="newest" />)
    const newestButton = screen.getByText('Newest')
    expect(newestButton).toHaveClass('bg-primary-600')
  })

  it('calls onSortChange when sort button clicked', () => {
    const onSortChange = vi.fn()
    render(<TopBar {...defaultProps} onSortChange={onSortChange} />)

    fireEvent.click(screen.getByText('Top Score'))
    expect(onSortChange).toHaveBeenCalledWith('score')
  })

  it('displays total count', () => {
    render(<TopBar {...defaultProps} totalCount={123} />)
    expect(screen.getByText(/123/)).toBeInTheDocument()
  })

  it('calls onRefresh when refresh button clicked', () => {
    const onRefresh = vi.fn()
    render(<TopBar {...defaultProps} onRefresh={onRefresh} />)

    // Find the refresh button by its text or nearby text
    const refreshButton = screen.getByRole('button', { name: /refresh/i })
    fireEvent.click(refreshButton)

    expect(onRefresh).toHaveBeenCalled()
  })

  it('shows refreshing state', () => {
    const { container } = render(<TopBar {...defaultProps} isRefreshing={true} />)

    // Should show spinning animation
    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
  })
})
