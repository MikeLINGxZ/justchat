import { useState } from 'react'
import { AlertViewport } from '@/components/alert/AlertViewport'
import { AlertEventProvider } from '@/components/providers/AlertEventProvider'
import { AppSettingsSyncProvider } from '@/components/providers/AppSettingsSyncProvider'
import { FontSizeProvider } from '@/components/providers/FontSizeProvider'
import { ThemeProvider } from '@/components/providers/ThemeProvider'
import type { SupportedProvider } from '@/types/settings'
import { OnboardingStepForm } from './OnboardingStepForm'
import { OnboardingStepIntro } from './OnboardingStepIntro'
import { OnboardingStepRuntime } from './OnboardingStepRuntime'
import { OnboardingStepSelect } from './OnboardingStepSelect'
import { Stepper } from './Stepper'

export type OnboardingStep = 'intro' | 'select' | 'config' | 'runtime'

export function OnboardingApp() {
	return (
		<AppSettingsSyncProvider>
			<AlertEventProvider>
				<ThemeProvider>
					<FontSizeProvider>
						<OnboardingWizard />
						<AlertViewport />
					</FontSizeProvider>
				</ThemeProvider>
			</AlertEventProvider>
		</AppSettingsSyncProvider>
	)
}

function OnboardingWizard() {
	const [step, setStep] = useState<OnboardingStep>('intro')
	const [provider, setProvider] = useState<SupportedProvider | null>(null)

	return (
		<div className="flex h-screen flex-col overflow-hidden bg-background text-foreground">
			<div className="shrink-0 px-6 pb-6 pt-12">
				<Stepper current={step} />
			</div>
			<div className="flex min-h-0 flex-1 flex-col overflow-hidden px-6 pb-6">
				{step === 'intro' && (
					<OnboardingStepIntro onStart={() => setStep('select')} />
				)}
				{step === 'select' && (
					<OnboardingStepSelect
						selected={provider}
						onSelect={setProvider}
						onBack={() => setStep('intro')}
						onNext={() => setStep('config')}
					/>
				)}
				{step === 'config' && provider && (
					<OnboardingStepForm
						provider={provider}
						onBack={() => setStep('select')}
						onDone={() => setStep('runtime')}
					/>
				)}
				{step === 'runtime' && (
					<OnboardingStepRuntime onBack={() => setStep('config')} />
				)}
			</div>
		</div>
	)
}
