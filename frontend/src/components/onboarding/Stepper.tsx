import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import type { OnboardingStep } from './OnboardingApp'

type Props = {
	current: OnboardingStep
}

const order: OnboardingStep[] = ['intro', 'select', 'config', 'runtime']

const labelKeys: Record<OnboardingStep, string> = {
	intro: 'onboarding.steps.intro',
	select: 'onboarding.steps.select',
	config: 'onboarding.steps.config',
	runtime: 'onboarding.steps.runtime',
}

export function Stepper(props: Props) {
	const { t } = useTranslation()
	const currentIdx = order.indexOf(props.current)

	return (
		<div className="flex items-center gap-3">
			{order.map((step, idx) => {
				const active = idx === currentIdx
				const done = idx < currentIdx
				return (
					<div key={step} className="flex items-center gap-3">
						{idx > 0 && <div className="h-px w-10 bg-border" />}
						<div className="flex items-center gap-2">
							<div
								className={cn(
									'flex h-6 w-6 items-center justify-center rounded-full text-[11px] font-semibold',
									active
										? 'bg-primary text-primary-foreground'
										: done
											? 'bg-primary/30 text-primary'
											: 'bg-muted text-muted-foreground',
								)}
							>
								{idx + 1}
							</div>
							<span
								className={cn(
									'text-sm',
									active ? 'font-medium text-foreground' : 'text-muted-foreground',
								)}
							>
								{t(labelKeys[step])}
							</span>
						</div>
					</div>
				)
			})}
		</div>
	)
}
