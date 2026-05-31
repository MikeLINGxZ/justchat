import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { SkillsList } from '@/components/settings/skills/SkillsList'
import { SkillsListItem } from '@/components/settings/skills/SkillsListItem'
import { SkillsDetailPane } from '@/components/settings/skills/SkillsDetailPane'
import i18n from '@/i18n'
import { getSettingsInitialState, useSettingsStore } from '@/store/settingsStore'
import type { SkillFull, SkillItem } from '@/types/skills'

const {
  getSkillMock,
  listSkillsMock,
  toggleSkillMock,
  updateSkillMock,
  deleteSkillMock,
  openAddSkillMock,
} = vi.hoisted(() => ({
  getSkillMock: vi.fn(),
  listSkillsMock: vi.fn(),
  toggleSkillMock: vi.fn(),
  updateSkillMock: vi.fn(),
  deleteSkillMock: vi.fn(),
  openAddSkillMock: vi.fn(),
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills', () => ({
  Skills: {
    ListSkills: listSkillsMock,
    ToggleSkill: toggleSkillMock,
    GetSkill: getSkillMock,
    UpdateSkill: updateSkillMock,
    DeleteSkill: deleteSkillMock,
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window', () => ({
  Window: {
    OpenAddSkill: openAddSkillMock,
  },
}))

describe('Skills settings', () => {
  beforeEach(async () => {
    useSettingsStore.setState(getSettingsInitialState())
    getSkillMock.mockReset()
    listSkillsMock.mockReset()
    toggleSkillMock.mockReset()
    updateSkillMock.mockReset()
    deleteSkillMock.mockReset()
    openAddSkillMock.mockReset()
    listSkillsMock.mockResolvedValue({ skills: [] })
    openAddSkillMock.mockResolvedValue(undefined)
    await i18n.changeLanguage('en')
  })

  it('renders the enable control as a right-side switch without selecting the row', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    const onToggle = vi.fn()
    const onDelete = vi.fn()
    const item: SkillItem = {
      name: 'local-skill',
      description: 'Local skill description',
      source: 'user',
      disabled: false,
      has_body: true,
    }

    render(
      <SkillsListItem
        item={item}
        selected={false}
        onSelect={onSelect}
        onToggle={onToggle}
        onDelete={onDelete}
      />
    )

    const toggle = screen.getByRole('switch', { name: 'Enabled' })
    expect(toggle).toHaveAttribute('aria-checked', 'true')

    await user.click(toggle)

    expect(onToggle).toHaveBeenCalledWith(item)
    expect(onSelect).not.toHaveBeenCalled()
  })

  it('shows the skill name as the detail title and styles the built-in notice as a card', async () => {
    const builtinSkill: SkillFull = {
      name: 'install-cli-from-docs',
      description: 'Built in helper',
      source: 'builtin',
      disabled: false,
      has_body: true,
      body: '# Builtin',
    }

    useSettingsStore.setState({
      ...getSettingsInitialState(),
      selectedSkillName: builtinSkill.name,
      skills: [builtinSkill],
    })
    getSkillMock.mockResolvedValue({ skill: builtinSkill })

    render(<SkillsDetailPane />)

    expect(await screen.findByRole('heading', { name: /install-cli-from-docs/ })).toBeInTheDocument()
    expect(screen.getByText('Built-in skill')).toBeInTheDocument()
    expect(screen.getByText('Built-in skill').closest('.rounded-xl')).toBeInTheDocument()
    await waitFor(() => {
      expect(getSkillMock).toHaveBeenCalled()
    })
  })

  it('opens a dedicated add-skill window instead of replacing the detail pane', async () => {
    const user = userEvent.setup()

    render(<SkillsList />)

    await user.click(screen.getByRole('button', { name: 'New Skill' }))

    expect(openAddSkillMock).toHaveBeenCalledWith({})
    expect(useSettingsStore.getState().selectedSkillName).toBeNull()
  })

  it('sends the edited skill name when saving an existing skill', async () => {
    const user = userEvent.setup()
    const existingSkill: SkillFull = {
      name: 'editable-skill',
      description: 'Before',
      source: 'user',
      disabled: false,
      has_body: true,
      body: 'Before body',
    }
    const renamedSkill: SkillFull = {
      ...existingSkill,
      name: 'renamed-skill',
      description: 'After',
      body: 'After body',
    }

    useSettingsStore.setState({
      ...getSettingsInitialState(),
      selectedSkillName: existingSkill.name,
      skills: [existingSkill],
    })
    getSkillMock.mockResolvedValue({ skill: existingSkill })
    updateSkillMock.mockResolvedValue({ skill: renamedSkill })

    render(<SkillsDetailPane />)

    const textboxes = await screen.findAllByRole('textbox')
    await user.clear(textboxes[0])
    await user.type(textboxes[0], renamedSkill.name)
    await user.clear(textboxes[1])
    await user.type(textboxes[1], renamedSkill.description)
    await user.clear(textboxes[2])
    await user.type(textboxes[2], renamedSkill.body)
    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(updateSkillMock).toHaveBeenCalledWith(expect.objectContaining({
      name: existingSkill.name,
      new_name: renamedSkill.name,
      description: renamedSkill.description,
      body: renamedSkill.body,
    }))
  })

  it('shows a validation message and skips saving when the edited skill name is invalid', async () => {
    const user = userEvent.setup()
    const existingSkill: SkillFull = {
      name: 'editable-skill',
      description: 'Before',
      source: 'user',
      disabled: false,
      has_body: true,
      body: 'Before body',
    }

    useSettingsStore.setState({
      ...getSettingsInitialState(),
      selectedSkillName: existingSkill.name,
      skills: [existingSkill],
    })
    getSkillMock.mockResolvedValue({ skill: existingSkill })

    render(<SkillsDetailPane />)

    const textboxes = await screen.findAllByRole('textbox')
    await user.clear(textboxes[0])
    await user.type(textboxes[0], 'test-skillaa11aa啵啵')
    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByText('Name must include English letters or numbers and will be saved in lowercase kebab-case.')).toBeInTheDocument()
    expect(updateSkillMock).not.toHaveBeenCalled()
  })
})
