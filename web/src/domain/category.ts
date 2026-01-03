type Category = {
  id: number
  name: string
  icon: string
}

type CategoryGroup = {
  id: number
  name: string
  icon: string
  categories: Category[]
}

const DEFAULT_CATEGORY_ID = 16

export { DEFAULT_CATEGORY_ID }
export type { Category, CategoryGroup }
