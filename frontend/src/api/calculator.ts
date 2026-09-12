import type { CalculateRequest } from '../domain/calculator'
export type { CalculateRequest } from '../domain/calculator'
type ResponseData = { result?: number; error?: string }
export async function calculate(input: CalculateRequest): Promise<number> {
  const operands = input.operation === 'sqrt' ? [input.a] : [input.a, input.b]
  if (!operands.every(Number.isFinite)) throw new Error('Ingresa números válidos.')
  let response: Response
  try { response = await fetch('http://localhost:8080/api/calculate', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(input) }) }
  catch { throw new Error('No se pudo conectar con el servidor.') }
  if (!(response.headers.get('content-type') ?? '').includes('application/json')) throw new Error(`Respuesta inesperada del servidor (${response.status}).`)
  const data = await response.json() as ResponseData
  if (!response.ok) throw new Error(data.error ?? 'No se pudo realizar el cálculo.')
  if (typeof data.result !== 'number' || !Number.isFinite(data.result)) throw new Error('El servidor devolvió un resultado inválido.')
  return data.result
}
