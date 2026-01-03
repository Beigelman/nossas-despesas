import { ArrowRightIcon, SearchIcon } from 'lucide-react'
import { useRef, useState } from 'react'

import { useSearch } from '@/hooks/use-search'

import { Button } from './ui/button'
import { Input } from './ui/input'

function SearchBar() {
  const inputRef = useRef(null)
  const setSearch = useSearch((state) => state.setSearch)
  const [searchInput, setSearchInput] = useState('')
  const [isInputVisible, setInputVisible] = useState(false)

  const handleInputChange: React.ChangeEventHandler<HTMLInputElement> = (event) => {
    const newValue = event.target.value
    setSearchInput(newValue)
    // Check if the input is cleared
    if (newValue === '') {
      setSearch(newValue)
    }
  }

  const handleButtonClick = () => {
    setInputVisible((value) => (searchInput.length > 0 ? true : !value))
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    inputRef?.current?.focus()
  }

  return (
    <div className="flex items-center justify-center gap-1">
      <Button className="h-7 w-7 p-0" variant="ghost" onClick={handleButtonClick}>
        <SearchIcon size={16} />
      </Button>
      <div
        className={`transition-width flex items-center gap-1 overflow-hidden duration-500
        ${isInputVisible ? 'w-48 p-1' : 'w-0'}`}
      >
        <Input
          type="search"
          ref={inputRef}
          value={searchInput}
          onChange={handleInputChange}
          placeholder="Buscas despesas..."
          onBlur={() => setInputVisible((value) => (searchInput.length > 0 ? true : !value))}
          onFocus={() => setInputVisible(true)}
          onKeyDown={(e) => e.key === 'Enter' && setSearch(searchInput)}
          className="h-7"
        />
        <Button className="h-7 w-7 p-0" variant="ghost" onClick={() => setSearch(searchInput)}>
          <ArrowRightIcon size={16} />
        </Button>
      </div>
    </div>
  )
}

export { SearchBar }
