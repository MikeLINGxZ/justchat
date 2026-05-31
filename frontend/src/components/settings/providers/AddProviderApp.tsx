import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Window } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window'
import { AlertViewport } from '@/components/alert/AlertViewport'
import { ThemeProvider } from '@/components/providers/ThemeProvider'
import { FontSizeProvider } from '@/components/providers/FontSizeProvider'
import { AppSettingsSyncProvider } from '@/components/providers/AppSettingsSyncProvider'
import { AlertEventProvider } from '@/components/providers/AlertEventProvider'
import { AddProviderStepSelect } from '@/components/settings/providers/AddProviderStepSelect'
import { AddProviderStepForm } from '@/components/settings/providers/AddProviderStepForm'
import { AddProviderStepDone } from '@/components/settings/providers/AddProviderStepDone'
import { cn } from '@/lib/utils'
import type { SupportedProvider } from '@/types/settings'

type WizardStep = 1 | 2 | 3

export function AddProviderApp() {
  return (
    <AppSettingsSyncProvider>
      <AlertEventProvider>
        <ThemeProvider>
          <FontSizeProvider>
            <AddProviderWizard />
            <AlertViewport />
          </FontSizeProvider>
        </ThemeProvider>
      </AlertEventProvider>
    </AppSettingsSyncProvider>
  )
}

function AddProviderWizard() {
  const { t } = useTranslation()
  const [step, setStep] = useState<WizardStep>(1)
  const [selectedProvider, setSelectedProvider] = useState<SupportedProvider | null>(null)

  const steps = [
    t('settingsPage.addProvider.stepSelect'),
    t('settingsPage.addProvider.stepForm'),
    t('settingsPage.addProvider.stepDone'),
  ]

  return (
    <div className="flex h-screen flex-col overflow-hidden bg-background text-foreground">
      {/* Title */}
      <div className="shrink-0 px-6 pb-3 pt-12">
        <h1 className="text-lg font-semibold">{t('settingsPage.addProvider.title')}</h1>
      </div>

      {/* Stepper */}
      <div className="shrink-0 px-6 pb-4">
        <div className="flex items-center gap-2">
          {steps.map((label, idx) => {
            const num = (idx + 1) as WizardStep
            const active = num === step
            const done = num < step
            return (
              <div key={label} className="flex items-center gap-2">
                {idx > 0 && <div className="h-px w-8 bg-border" />}
                <div className="flex items-center gap-1.5">
                  <div
                    className={cn(
                      'flex h-5 w-5 items-center justify-center rounded-full text-xs font-semibold',
                      active
                        ? 'bg-primary text-primary-foreground'
                        : done
                          ? 'bg-primary/30 text-primary'
                          : 'bg-muted text-muted-foreground',
                    )}
                  >
                    {num}
                  </div>
                  <span
                    className={cn(
                      'text-xs',
                      active ? 'font-medium text-foreground' : 'text-muted-foreground',
                    )}
                  >
                    {label}
                  </span>
                </div>
              </div>
            )
          })}
        </div>
      </div>

      {/* Content */}
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden px-6">
        {step === 1 && (
          <div className="min-h-0 overflow-y-auto">
            <AddProviderStepSelect
              selected={selectedProvider}
              onSelect={setSelectedProvider}
            />
          </div>
        )}
        {step === 2 && selectedProvider && (
          <AddProviderStepForm
            provider={selectedProvider}
            onDone={() => setStep(3)}
          />
        )}
        {step === 3 && <AddProviderStepDone />}
      </div>

      {/* Footer */}
      <div className="shrink-0 border-t border-border px-6 pb-6 pt-4">
        {step === 1 && (
          <div className="flex justify-between">
            <button
              type="button"
              onClick={() => { void Window.CloseAddProvider({}) }}
              className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
            >
              {t('settingsPage.addProvider.cancel')}
            </button>
            <button
              type="button"
              disabled={!selectedProvider}
              onClick={() => setStep(2)}
              className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground transition-opacity disabled:opacity-40"
            >
              {t('settingsPage.addProvider.next')}
            </button>
          </div>
        )}
        {step === 2 && (
          <div className="flex justify-between">
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => setStep(1)}
                className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
              >
                {t('settingsPage.addProvider.prev')}
              </button>
              <button
                type="button"
                onClick={() => { void Window.CloseAddProvider({}) }}
                className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
              >
                {t('settingsPage.addProvider.cancel')}
              </button>
            </div>
            <button
              type="submit"
              form="add-provider-form"
              className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground"
            >
              {t('settingsPage.addProvider.add')}
            </button>
          </div>
        )}
        {step === 3 && (
          <div className="flex justify-end">
            <button
              type="button"
              onClick={() => { void Window.CloseAddProvider({}) }}
              className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground"
            >
              {t('settingsPage.addProvider.done')}
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
