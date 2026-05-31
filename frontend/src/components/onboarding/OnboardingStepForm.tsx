import { useTranslation } from 'react-i18next'
import { AddProviderStepForm } from '@/components/settings/providers/AddProviderStepForm'
import type { SupportedProvider } from '@/types/settings'

type Props = {
	provider: SupportedProvider
	onBack: () => void
	onDone: () => void
}

export function OnboardingStepForm(props: Props) {
	const { t } = useTranslation()

	return (
		<div className="flex h-full min-h-0 flex-col gap-4">
			<div className="space-y-1">
				<h2 className="text-lg font-semibold">{t('onboarding.config.title', { name: props.provider.name })}</h2>
				<p className="text-xs text-muted-foreground">
					{props.provider.description || t('onboarding.config.hintFallback')}
				</p>
			</div>

			<div className="min-h-0 flex-1 overflow-hidden">
				<AddProviderStepForm
					provider={props.provider}
					mode="onboarding"
					onDone={props.onDone}
				/>
			</div>

			<div className="flex shrink-0 items-center justify-between border-t border-border pt-4">
				<button
					type="button"
					onClick={props.onBack}
					className="rounded-xl border border-border px-5 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
				>
					{t('onboarding.actions.back')}
				</button>
				<button
					type="submit"
					form="add-provider-form"
					className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground"
				>
					{t('onboarding.actions.next')}
				</button>
			</div>
		</div>
	)
}
