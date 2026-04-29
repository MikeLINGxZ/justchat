package tasker

import (
	"fmt"
	"sync"
)

// Manager 管理按 taskUuid 区分的后台任务。
var Manager *manager

func init() {
	Manager = &manager{
		tasks:      make(map[string]*Runtime),
		maxWorkers: 5,
	}
}

type Runtime struct {
	TaskUUID             string
	ChatUUID             string
	AssistantMessageUUID string
	EventKey             string
	stopCh               chan struct{}
}

type manager struct {
	mu            sync.Mutex
	tasks         map[string]*Runtime
	maxWorkers    int
	activeWorkers int
	concurrencyCh chan struct{}
}

func (m *manager) SetMaxWorkers(max int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if max <= 0 {
		max = 5
	}
	m.maxWorkers = max
	if m.concurrencyCh != nil {
		m.concurrencyCh = make(chan struct{}, max)
	}
}

func (m *manager) StartTask(task Runtime, fn func(stop <-chan struct{})) {
	userStop := make(chan struct{}, 1)
	task.stopCh = userStop

	m.mu.Lock()
	m.tasks[task.TaskUUID] = &task
	if m.concurrencyCh == nil {
		m.concurrencyCh = make(chan struct{}, m.maxWorkers)
	}
	m.mu.Unlock()

	go func() {
		m.mu.Lock()
		m.activeWorkers++
		m.mu.Unlock()
		defer func() {
			m.mu.Lock()
			m.activeWorkers--
			delete(m.tasks, task.TaskUUID)
			m.mu.Unlock()
			close(userStop)
		}()

		select {
		case m.concurrencyCh <- struct{}{}:
		default:
		}
		defer func() { <-m.concurrencyCh }()

		fn(userStop)
	}()
}

func (m *manager) TryStartTask(task Runtime, fn func(stop <-chan struct{})) error {
	m.mu.Lock()
	if len(m.tasks) >= m.maxWorkers {
		m.mu.Unlock()
		return fmt.Errorf("too many concurrent tasks (max: %d)", m.maxWorkers)
	}
	m.mu.Unlock()

	m.StartTask(task, fn)
	return nil
}

func (m *manager) CancelAndReplace(chatUUID string, newTask Runtime, fn func(stop <-chan struct{})) error {
	m.mu.Lock()
	for taskUUID, existingTask := range m.tasks {
		if existingTask.ChatUUID == chatUUID {
			m.mu.Unlock()
			m.StopTask(taskUUID)
			m.mu.Lock()
			break
		}
	}
	m.mu.Unlock()

	m.StartTask(newTask, fn)
	return nil
}

func (m *manager) ActiveTaskCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.tasks)
}

func (m *manager) StopTask(taskUUID string) {
	m.mu.Lock()
	task, ok := m.tasks[taskUUID]
	m.mu.Unlock()
	if !ok || task == nil || task.stopCh == nil {
		return
	}
	select {
	case task.stopCh <- struct{}{}:
	default:
	}
}

func (m *manager) StopByEventKey(eventKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, task := range m.tasks {
		if task == nil || task.EventKey != eventKey || task.stopCh == nil {
			continue
		}
		select {
		case task.stopCh <- struct{}{}:
		default:
		}
		return
	}
}

func (m *manager) GetTaskRuntime(taskUUID string) (*Runtime, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task, ok := m.tasks[taskUUID]
	if !ok || task == nil {
		return nil, false
	}
	copyTask := *task
	copyTask.stopCh = nil
	return &copyTask, true
}

func (m *manager) ListRunningTasks() []Runtime {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := make([]Runtime, 0, len(m.tasks))
	for _, task := range m.tasks {
		if task == nil {
			continue
		}
		copyTask := *task
		copyTask.stopCh = nil
		res = append(res, copyTask)
	}
	return res
}
