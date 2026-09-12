import { useRef, useState, type FormEvent } from 'react'
import { executeCalculation } from '../application/calculate'
import { operationOptions, toCalculateRequest, validateInput, type Operation } from '../domain/calculator'

function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const requestVersion = useRef(0)
  const unary = operation === 'sqrt'
  const selected = operationOptions.find(item => item.value === operation) ?? operationOptions[0]

  function resetCalculation() {
    // Ignore responses for inputs that have since changed.
    requestVersion.current += 1
    setResult(null)
    setError('')
    setLoading(false)
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    resetCalculation()
    const validationError = validateInput(operation, a, b)
    if (validationError) {
      setError(validationError)
      return
    }

    const version = requestVersion.current
    setLoading(true)
    try {
      const value = await executeCalculation(toCalculateRequest(operation, a, b))
      if (version === requestVersion.current) setResult(value)
    } catch (cause) {
      if (version === requestVersion.current) {
        setError(cause instanceof Error ? cause.message : 'No se pudo realizar el cálculo.')
      }
    } finally {
      if (version === requestVersion.current) setLoading(false)
    }
  }

  return (
    <main className="page">
      <section className="calculator" aria-labelledby="title">
        <p className="eyebrow">Sezzle · Prueba técnica</p>
        <h1 id="title">Calculadora</h1>
        <p className="intro">Realiza operaciones matemáticas de forma rápida y sencilla.</p>
        <div className="expression" aria-label="Vista previa de la operación">
          {unary && <strong>{selected.symbol}</strong>}
          <span>{a || 'a'}</span>
          {!unary && <><strong>{selected.symbol}</strong><span>{b || 'b'}</span></>}
          <strong>=</strong>
          <span className="expression-result">{result ?? '?'}</span>
        </div>
        <form onSubmit={submit} noValidate>
          <label htmlFor="operation">Operación</label>
          <select id="operation" value={operation} onChange={event => {
            setOperation(event.target.value as Operation)
            resetCalculation()
          }}>
            {operationOptions.map(item => <option key={item.value} value={item.value}>{item.label}</option>)}
          </select>
          <div className="inputs">
            <div>
              <label htmlFor="a">{unary ? 'Número' : 'Primer número'}</label>
              <input id="a" inputMode="decimal" value={a} onChange={event => {
                setA(event.target.value)
                resetCalculation()
              }} placeholder="Ej. 10" />
            </div>
            {!unary && <div>
              <label htmlFor="b">Segundo número</label>
              <input id="b" inputMode="decimal" value={b} onChange={event => {
                setB(event.target.value)
                resetCalculation()
              }} placeholder="Ej. 5" />
            </div>}
          </div>
          <button type="submit" disabled={loading}>{loading ? 'Calculando…' : 'Calcular'}</button>
        </form>
        {error && <p className="message error" role="alert">{error}</p>}
        {result !== null && <div className="message success" role="status">
          <span>Resultado</span><strong>{result}</strong>
        </div>}
      </section>
    </main>
  )
}
export default Calculator
