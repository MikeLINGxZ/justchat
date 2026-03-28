package tasker

import "sync"

// Manager 管理按 taskUuid 区分的后台任务。
var Manager *manager

func init() {
	Manager = &manager{
		tasks: make(map[string]*Runtime),
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
	mu    sync.Mutex
	tasks map[string]*Runtime
}

func (m *manager) StartTask(task Runtime, fn func(stop <-chan struct{})) {
	userStop := make(chan struct{}, 1)
	task.stopCh = userStop

	m.mu.Lock()
	m.tasks[task.TaskUUID] = &task
	m.mu.Unlock()

	go func() {
		defer func() {
			m.mu.Lock()
			delete(m.tasks, task.TaskUUID)
			m.mu.Unlock()
			close(userStop)
		}()
		fn(userStop)
	}()
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
