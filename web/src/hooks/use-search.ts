import { create } from 'zustand'

interface SearchState {
  search: string
  setSearch: (value: string) => void
}

const useSearch = create<SearchState>((set) => ({
  search: '',
  setSearch: (value: string) => set(() => ({ search: value })),
}))

export { useSearch }
