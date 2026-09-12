import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import Calculator from './Calculator'

function jsonResponse(data: { result: number } | { error: string }, status = 200) {
  return new Response(JSON.stringify(data), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

function deferredResponse() {
  let resolve!: (response: Response) => void
  const promise = new Promise<Response>(resolvePromise => { resolve = resolvePromise })
  return { promise, resolve }
}

function enterOperands() {
  fireEvent.change(screen.getByLabelText('Primer número'), { target: { value: '10' } })
  fireEvent.change(screen.getByLabelText('Segundo número'), { target: { value: '5' } })
}

describe('Calculator', () => {
  afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals() })
  it('renders initially', () => { render(<Calculator />); expect(screen.getByRole('heading', { name: 'Calculadora' })).toBeInTheDocument() })
  it('shows success', async () => { vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ result: 15 }), { status: 200, headers: { 'Content-Type': 'application/json' } }))); render(<Calculator />); fireEvent.change(screen.getByLabelText('Primer número'), { target: { value: '10' } }); fireEvent.change(screen.getByLabelText('Segundo número'), { target: { value: '5' } }); fireEvent.click(screen.getByRole('button', { name: 'Calcular' })); await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('15')) })
  it('validates empty inputs', () => { render(<Calculator />); fireEvent.click(screen.getByRole('button', { name: 'Calcular' })); expect(screen.getByRole('alert')).toHaveTextContent('Completa todos los campos requeridos.') })
  it('shows backend errors', async () => { vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ error: 'division by zero' }), { status: 400, headers: { 'Content-Type': 'application/json' } }))); render(<Calculator />); fireEvent.change(screen.getByLabelText('Primer número'), { target: { value: '10' } }); fireEvent.change(screen.getByLabelText('Segundo número'), { target: { value: '2' } }); fireEvent.change(screen.getByLabelText('Operación'), { target: { value: 'power' } }); fireEvent.click(screen.getByRole('button', { name: 'Calcular' })); await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('division by zero')) })

  it.each(['Primer número', 'Segundo número'])('clears the result when editing %s', async label => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ result: 15 })))
    render(<Calculator />)
    enterOperands()
    fireEvent.click(screen.getByRole('button', { name: 'Calcular' }))
    expect(await screen.findByRole('status')).toHaveTextContent('15')

    fireEvent.change(screen.getByLabelText(label), { target: { value: '20' } })

    expect(screen.queryByRole('status')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Vista previa de la operación')).toHaveTextContent(/=\?$/)
  })

  it('clears validation errors when editing an operand', () => {
    render(<Calculator />)
    fireEvent.click(screen.getByRole('button', { name: 'Calcular' }))
    expect(screen.getByRole('alert')).toBeInTheDocument()

    fireEvent.change(screen.getByLabelText('Primer número'), { target: { value: '25' } })

    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('shows the square root symbol before its operand and sends only a', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ result: 5 }))
    vi.stubGlobal('fetch', fetchMock)
    render(<Calculator />)
    enterOperands()
    fireEvent.change(screen.getByLabelText('Operación'), { target: { value: 'sqrt' } })
    fireEvent.change(screen.getByLabelText('Número'), { target: { value: '25' } })

    expect(screen.queryByLabelText('Segundo número')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Vista previa de la operación')).toHaveTextContent('√25=?')
    fireEvent.click(screen.getByRole('button', { name: 'Calcular' }))

    expect(await screen.findByRole('status')).toHaveTextContent('5')
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:8080/api/calculate', expect.objectContaining({
      body: JSON.stringify({ operation: 'sqrt', a: 25 }),
    }))
    expect(screen.getByLabelText('Vista previa de la operación')).toHaveTextContent('√25=5')
  })

  describe.each([
    { label: 'Primer número', value: '20' },
    { label: 'Segundo número', value: '20' },
    { label: 'Operación', value: 'sqrt' },
  ])('after editing $label during a request', ({ label, value }) => {
    it.each(['success', 'error'])('ignores the previous %s response', async outcome => {
      const pending = deferredResponse()
      vi.stubGlobal('fetch', vi.fn().mockReturnValue(pending.promise))
      render(<Calculator />)
      enterOperands()
      fireEvent.click(screen.getByRole('button', { name: 'Calcular' }))
      expect(screen.getByRole('button', { name: 'Calculando…' })).toBeDisabled()

      fireEvent.change(screen.getByLabelText(label), { target: { value } })
      expect(screen.getByRole('button', { name: 'Calcular' })).toBeEnabled()
      await act(async () => {
        pending.resolve(outcome === 'success'
          ? jsonResponse({ result: 15 })
          : jsonResponse({ error: 'Old error' }, 400))
      })

      expect(screen.queryByRole('status')).not.toBeInTheDocument()
      expect(screen.queryByRole('alert')).not.toBeInTheDocument()
      expect(screen.getByLabelText('Vista previa de la operación')).toHaveTextContent(/=\?$/)
    })
  })

  it.each(['old first', 'new first'])('keeps the latest calculation when responses arrive %s', async order => {
    const oldRequest = deferredResponse()
    const newRequest = deferredResponse()
    vi.stubGlobal('fetch', vi.fn()
      .mockReturnValueOnce(oldRequest.promise)
      .mockReturnValueOnce(newRequest.promise))
    render(<Calculator />)
    enterOperands()
    fireEvent.click(screen.getByRole('button', { name: 'Calcular' }))
    fireEvent.change(screen.getByLabelText('Primer número'), { target: { value: '20' } })
    fireEvent.click(screen.getByRole('button', { name: 'Calcular' }))

    if (order === 'old first') {
      await act(async () => { oldRequest.resolve(jsonResponse({ result: 15 })) })
      expect(screen.getByRole('button', { name: 'Calculando…' })).toBeDisabled()
      expect(screen.queryByRole('status')).not.toBeInTheDocument()
      await act(async () => { newRequest.resolve(jsonResponse({ result: 25 })) })
    } else {
      await act(async () => { newRequest.resolve(jsonResponse({ result: 25 })) })
      await act(async () => { oldRequest.resolve(jsonResponse({ result: 15 })) })
    }

    expect(screen.getByRole('status')).toHaveTextContent('25')
    expect(screen.getByLabelText('Vista previa de la operación')).toHaveTextContent('20+5=25')
    expect(screen.getByRole('button', { name: 'Calcular' })).toBeEnabled()
  })
})
