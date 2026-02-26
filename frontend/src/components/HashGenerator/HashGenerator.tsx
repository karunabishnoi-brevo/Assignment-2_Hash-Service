import { useState, useCallback, type FormEvent, type ChangeEvent } from 'react'
import { generateHash, type HashResponse } from '../../services/api'
import './HashGenerator.css'

const ALPHANUMERIC_REGEX = /^[a-zA-Z0-9]*$/

function HashGenerator() {
  const [input, setInput] = useState('')
  const [result, setResult] = useState<HashResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleInputChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    if (ALPHANUMERIC_REGEX.test(value)) {
      setInput(value)
      setError('')
    }
  }, [])

  const handleSubmit = useCallback(
    async (e: FormEvent<HTMLFormElement>) => {
      e.preventDefault()

      const trimmed = input.trim()
      if (!trimmed) {
        setError('Please enter a value.')
        return
      }

      setLoading(true)
      setError('')
      setResult(null)

      try {
        const data = await generateHash(trimmed)
        setResult(data)
      } catch (err) {
        const message =
          err instanceof Error ? err.message : 'An unexpected error occurred.'
        setError(message)
      } finally {
        setLoading(false)
      }
    },
    [input],
  )

  const isSubmitDisabled = input.trim().length === 0 || loading

  return (
    <div className="hash-generator">
      <h1 className="hash-generator__title">Hash Generator</h1>
      <p className="hash-generator__subtitle">
        Enter an alphanumeric string to generate a unique 10-character hash.
      </p>

      <form className="hash-generator__form" onSubmit={handleSubmit}>
        <label htmlFor="hash-input" className="hash-generator__label">
          Input
        </label>
        <div className="hash-generator__input-row">
          <input
            id="hash-input"
            type="text"
            className="hash-generator__input"
            value={input}
            onChange={handleInputChange}
            placeholder="e.g. abc123"
            autoComplete="off"
            aria-describedby={error ? 'hash-error' : undefined}
            aria-invalid={error ? true : undefined}
            disabled={loading}
          />
          <button
            type="submit"
            className="hash-generator__button"
            disabled={isSubmitDisabled}
            aria-busy={loading}
          >
            {loading ? 'Generating...' : 'Generate'}
          </button>
        </div>
      </form>

      {error && (
        <p id="hash-error" className="hash-generator__error" role="alert">
          {error}
        </p>
      )}

      {result && (
        <div className="hash-generator__result" aria-live="polite">
          <span className="hash-generator__result-label">Generated Hash</span>
          <output className="hash-generator__hash">{result.hash}</output>
          <span className="hash-generator__result-input">
            Input: <code>{result.input}</code>
          </span>
        </div>
      )}
    </div>
  )
}

export default HashGenerator
