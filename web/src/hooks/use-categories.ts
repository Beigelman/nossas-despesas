import { useQuery } from '@tanstack/react-query'

import { Category, CategoryGroup } from '@/domain/category'
import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'

type GetCategoriesApiResponse = ApiResponse<
  {
    id: number
    name: string
    icon: string
    categories: {
      id: number
      name: string
      icon: string
    }[]
  }[]
>

function useCategories() {
  const { data, isError, isPending } = useQuery<CategoryGroup[]>({
    queryKey: ['categories'],
    queryFn: async (): Promise<CategoryGroup[]> => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetCategoriesApiResponse>('/category')

      const categoryGroups = data.map<CategoryGroup>((categoryGroup) => ({
        id: categoryGroup.id,
        name: categoryGroup.name,
        icon: categoryGroup.icon,
        categories: categoryGroup.categories.map<Category>((category) => ({
          id: category.id,
          name: category.name,
          icon: category.icon,
        })),
      }))
      return categoryGroups
    },
  })

  const getCategory = (categoryId: number): Category | undefined => {
    return data?.flatMap((categoryGroup) => categoryGroup.categories).find((category) => category.id === categoryId)
  }

  const getGroupCategory = (categoryId: number): CategoryGroup | undefined => {
    return data?.find((categoryGroup) => categoryGroup.id === categoryId)
  }

  const orderedCategories = data ? orderCategories(data) : []

  return {
    getCategory,
    getGroupCategory,
    categories: orderedCategories,
    isLoading: isPending,
    isError,
  }
}

function orderCategories(categoryGroups: CategoryGroup[]): CategoryGroup[] {
  return categoryGroups.map((categoryGroup) => {
    const other = categoryGroup.categories.find((category) => category.name === 'Outros')
    const orderedCategories = categoryGroup.categories
      .filter((c) => c.name !== 'Outros')
      .sort((a, b) => a.name.localeCompare(b.name))

    if (other) {
      orderedCategories.push(other)
    }

    return { ...categoryGroup, categories: orderedCategories }
  })
}

export { useCategories }
