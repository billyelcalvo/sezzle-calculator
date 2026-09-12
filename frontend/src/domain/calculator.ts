export type Operation = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt' | 'percentage'

export type CalculateRequest = { operation: 'sqrt'; a: number } | { operation: Exclude<Operation, 'sqrt'>; a: number; b: number }

export type OperationOption = { value: Operation; label: string; symbol: string; unary: boolean }

export const operationOptions: OperationOption[] = [
  { value: 'add', label: 'Suma', symbol: '+', unary: false }, { value: 'subtract', label: 'Resta', symbol: '−', unary: false }, { value: 'multiply', label: 'Multiplicación', symbol: '×', unary: false }, { value: 'divide', label: 'División', symbol: '÷', unary: false }, { value: 'power', label: 'Potencia', symbol: '^', unary: false }, { value: 'sqrt', label: 'Raíz cuadrada', symbol: '√', unary: true }, { value: 'percentage', label: 'Porcentaje', symbol: '%', unary: false },
]
export function validateInput(operation: Operation, a: string, b: string): string | null {
   const unary = operation === 'sqrt'; if (!a.trim() || (!unary && !b.trim())) return 'Completa todos los campos requeridos.';
   const first = Number(a); const second = Number(b); if (!Number.isFinite(first) || (!unary && !Number.isFinite(second))) return 'Ingresa números válidos.';
   if (unary && first < 0) return 'La raíz cuadrada requiere un número mayor o igual a cero.';
   if (operation === 'divide' && second === 0) return 'No se puede dividir por cero.';
   return null
   }
export function toCalculateRequest(operation: Operation, a: string, b: string): CalculateRequest {
  const first = Number(a);
  const second = Number(b);
   return operation === 'sqrt' ? { operation, a: first } : { operation, a: first, b: second }
  }
