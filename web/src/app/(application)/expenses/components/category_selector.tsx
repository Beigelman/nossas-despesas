import { useState } from 'react'

import { CategoryIcon } from '@/components/category-icon'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Category, CategoryGroup } from '@/domain/category'
import { cn } from '@/lib/utils'

type CategorySelectorProps = {
  open?: boolean
  setOpen?: (open: boolean) => void
  handleClick: (value: number) => void
  getCategory: (id: number) => Category | undefined
  selectedICategoryID: number
  iconSize?: number
  isPredicting?: boolean
  isPredicted?: boolean
  categories: CategoryGroup[]
}

function CategorySelectorMenu({
  open,
  setOpen,
  handleClick,
  categories,
  getCategory,
  selectedICategoryID,
  iconSize,
  isPredicting,
  isPredicted,
}: CategorySelectorProps) {
  const [internalOpen, setInternalOpen] = useState(false)

  return (
    <DropdownMenu open={open ?? internalOpen} onOpenChange={setOpen ?? setInternalOpen}>
      <div className="relative inline-flex items-center justify-center">
        <DropdownMenuTrigger>
          <CategoryIcon
            name={getCategory(selectedICategoryID)?.icon ?? ''}
            size={iconSize ?? 34}
            hover
            isPredicted={isPredicted}
            isPredicting={isPredicting}
            className="p-3 shadow shadow-zinc-200 dark:shadow-zinc-800"
          />
        </DropdownMenuTrigger>
      </div>
      <DropdownMenuContent className="w-57">
        {categories.map((categoryGroup) => (
          <DropdownMenuSub key={categoryGroup.name}>
            <DropdownMenuSubTrigger>
              <CategoryIcon name={categoryGroup.icon} size={17} hover className="p-2" blendBackground />
              <span>{categoryGroup.name}</span>
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                {categoryGroup.categories.map((category) => (
                  <DropdownMenuItem
                    key={category.name}
                    onClick={() => handleClick(category.id)}
                    className={cn(
                      category.id === selectedICategoryID && 'bg-accent text-accent-foreground',
                      'cursor-pointer',
                    )}
                  >
                    <CategoryIcon name={category.icon} size={17} className="p-2" blendBackground />
                    {category.name}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export { CategorySelectorMenu }
