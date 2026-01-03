import { icons, LucideProps } from 'lucide-react'

interface IconProps extends LucideProps {
  name: string
}

const Icon = ({ name, ...props }: IconProps) => {
  const LucideIcon = icons[name as keyof typeof icons]
  if (!LucideIcon) {
    return <div style={{ background: '#ddd', width: 24, height: 24 }} />
  }

  return <LucideIcon {...props} />
}

export { Icon }
export type { IconProps }
