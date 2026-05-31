import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Config } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/config'
import { cn } from '@/lib/utils'
import type { SupportedProvider } from '@/types/settings'

type Props = {
	selected: SupportedProvider | null
	onSelect: (provider: SupportedProvider) => void
	onBack: () => void
	onNext: () => void
}

export function OnboardingStepSelect(props: Props) {
	const { t } = useTranslation()
	const [providers, setProviders] = useState<SupportedProvider[]>([])
	const [loadFailed, setLoadFailed] = useState(false)

	useEffect(() => {
		let cancelled = false

		Config.SupportedProviderList({})
			.then((result) => {
				if (cancelled) {
					return
				}

				if (result?.supported_providers) {
					setProviders(
						result.supported_providers.map((provider) => ({
							type: provider.type,
							icon: provider.icon,
							name: provider.name,
							description: provider.description,
							base_url: provider.base_url,
						})),
					)
				}
			})
			.catch(() => {
				if (!cancelled) {
					setLoadFailed(true)
				}
			})

		return () => {
			cancelled = true
		}
	}, [])

	return (
		<div className="flex h-full min-h-0 flex-col gap-4">
			<div className="space-y-1">
				<h2 className="text-lg font-semibold">{t('onboarding.steps.select')}</h2>
				<p className="text-xs text-muted-foreground">
					{loadFailed ? t('onboarding.select.loadFailed') : t('onboarding.select.hint')}
				</p>
			</div>

			<div className="min-h-0 flex-1 overflow-y-auto pr-1">
				{providers.length === 0 && !loadFailed && (
					<div className="py-8 text-center text-xs text-muted-foreground">...</div>
				)}
				{providers.length === 0 && loadFailed && (
					<div className="py-8 text-center text-xs text-muted-foreground">
						{t('onboarding.select.empty')}
					</div>
				)}
				<div className="space-y-2">
					{providers.map((provider) => (
						<button
							key={provider.type}
							type="button"
							aria-label={provider.name}
							onClick={() => props.onSelect(provider)}
							className={cn(
								'flex w-full items-start gap-4 rounded-2xl border p-4 text-left transition-colors',
								props.selected?.type === provider.type
									? 'border-primary bg-primary/10'
									: 'border-border bg-background hover:bg-accent',
							)}
						>
							{provider.icon && (
								<img
									src={provider.icon}
									alt={provider.name}
									className="h-10 w-10 rounded-xl object-contain"
									onError={(event) => {
										event.currentTarget.style.display = 'none'
									}}
								/>
							)}
							<div className="min-w-0 flex-1">
								<div className="text-sm font-semibold text-foreground">{provider.name}</div>
								<div className="mt-0.5 text-xs text-muted-foreground">{provider.description}</div>
								{provider.base_url && (
									<div className="mt-1 text-xs text-muted-foreground/70">{provider.base_url}</div>
								)}
							</div>
						</button>
					))}
				</div>
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
					type="button"
					disabled={!props.selected}
					onClick={props.onNext}
					className="rounded-xl bg-primary px-5 py-2 text-sm font-medium text-primary-foreground transition-opacity disabled:opacity-40"
				>
					{t('onboarding.actions.next')}
				</button>
			</div>
		</div>
	)
}
