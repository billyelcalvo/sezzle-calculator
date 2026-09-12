import { calculate } from '../api/calculator'
import type { CalculateRequest } from '../domain/calculator'
export function executeCalculation(input: CalculateRequest): Promise<number> { return calculate(input) }
