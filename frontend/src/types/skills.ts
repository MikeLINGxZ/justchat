export type SkillSource = 'builtin' | 'user' | 'ai'

export type SkillItem = {
  name: string
  description: string
  source: SkillSource
  disabled: boolean
  has_body: boolean
}

export type SkillFull = SkillItem & {
  body: string
}
