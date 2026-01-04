import { getTranslations } from 'next-intl/server'

import { SignInWithGoogleButton } from '@/components/signin-with-google-button'
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

export default async function SignInPage() {
  const t = await getTranslations()

  return (
    <div className="flex h-full w-full items-center justify-center bg-secondary">
      <div className="flex flex-col items-center justify-center">
        <Card className="min-w-[400px] shadow-lg">
          <CardHeader className="space-y-1">
            <CardTitle className="text-center text-2xl">{t('auth.welcome')}</CardTitle>
            <CardDescription className="text-center">{t('auth.signInDescription')}</CardDescription>
          </CardHeader>
          <CardFooter className="flex flex-col">
            <Separator className="my-3" />
            <SignInWithGoogleButton />
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}
