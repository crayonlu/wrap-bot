import { useState, useRef, useEffect } from 'react'
import { ChevronDown } from 'lucide-react'

interface Option {
  value: string | number
  label: string
}

interface SelectorProps {
  value: string | number
  onChange: (value: string | number) => void
  options: Option[]
  placeholder?: string
  className?: string
}

export default function Selector({ value, onChange, options, placeholder, className = '' }: SelectorProps) {
  const [isOpen, setIsOpen] = useState(false)
  const selectorRef = useRef<HTMLDivElement>(null)

  const selectedOption = options.find(opt => opt.value === value)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (selectorRef.current && !selectorRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleSelect = (optionValue: string | number) => {
    onChange(optionValue)
    setIsOpen(false)
  }

  return (
    <div className={`selector ${className}`} ref={selectorRef}>
      <button
        type="button"
        className="selector__trigger"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span className="selector__value">
          {selectedOption?.label || placeholder || 'Select...'}
        </span>
        <ChevronDown className={`selector__icon ${isOpen ? 'selector__icon--open' : ''}`} />
      </button>
      
      {isOpen && (
        <div className="selector__dropdown">
          {options.map((option) => (
            <button
              key={option.value}
              type="button"
              className={`selector__option ${option.value === value ? 'selector__option--selected' : ''}`}
              onClick={() => handleSelect(option.value)}
            >
              {option.label}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
