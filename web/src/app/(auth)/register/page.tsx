'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { signIn } from 'next-auth/react'
import { useTranslations } from 'next-intl'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

import { SignInWithGoogleButton } from '@/components/signin-with-google-button'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { signUpFormSchema } from '@/schemas'

export default function SignUpPage() {
  const form = useForm<z.infer<typeof signUpFormSchema>>({
    resolver: zodResolver(signUpFormSchema),
  })

  const router = useRouter()
  const t = useTranslations()

  async function onSubmit(data: z.infer<typeof signUpFormSchema>) {
    const result = await signIn('signup', {
      name: data.name,
      email: data.email,
      password: data.password,
      passwordConfirmation: data.passwordConfirmation,
      redirect: false,
    })

    if (result?.error) {
      toast.error(result.error)
      return
    }

    router.replace('/expenses')
  }

  return (
    <div className="flex h-full w-full items-center justify-center bg-secondary">
      <div className="flex flex-col items-center justify-center">
        <Card className="min-w-[400px] shadow-lg">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <CardHeader className="space-y-1">
                <CardTitle className="text-center text-2xl">{t('auth.createAccount')}</CardTitle>
                <CardDescription className="text-center">{t('auth.createAccountDescription')}</CardDescription>
              </CardHeader>
              <CardContent className="grid gap-4">
                <div className="grid gap-2">
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('auth.name')}</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="grid gap-2">
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('auth.email')}</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="grid gap-2">
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('auth.password')}</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="grid gap-2">
                  <FormField
                    control={form.control}
                    name="passwordConfirmation"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('auth.confirmPassword')}</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </CardContent>
              <CardFooter className="flex flex-col">
                <Button type="submit" className="w-full">
                  {t('auth.register')}
                </Button>
                <p className="mt-3 text-center text-xs text-gray-700">
                  {t('auth.alreadyHaveAccount')}{' '}
                  <Link href={'/login'}>
                    <span className=" text-blue-600 hover:underline">{t('auth.signInLink')}</span>
                  </Link>
                </p>
                <Separator className="my-3" />
                <SignInWithGoogleButton />
              </CardFooter>
            </form>
          </Form>
        </Card>
      </div>
    </div>
  )
}
