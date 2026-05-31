const skillNamePattern = /^[a-z0-9][a-z0-9-]{0,63}$/

export function isValidSkillName(value: string): boolean {
  return skillNamePattern.test(value.trim())
}

export function toSkillName(value: string): string {
  const name = value.trim()
  return isValidSkillName(name) ? name : ''
}
