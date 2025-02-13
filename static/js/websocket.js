// WebSocket connection
let ws;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 3000; // 3 seconds

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket connected');
        reconnectAttempts = 0;
        updateConnectionStatus('connected');
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected');
        updateConnectionStatus('disconnected');
        handleReconnect();
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        updateConnectionStatus('error');
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleWebSocketMessage(data);
    };
}

function handleReconnect() {
    if (reconnectAttempts < maxReconnectAttempts) {
        reconnectAttempts++;
        console.log(`Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts})...`);
        setTimeout(connectWebSocket, reconnectDelay);
    } else {
        console.error('Max reconnection attempts reached');
        showNotification('Connection lost. Please refresh the page.', 'error');
    }
}

function updateConnectionStatus(status) {
    const statusIndicator = document.getElementById('connection-status');
    if (statusIndicator) {
        statusIndicator.className = `connection-status ${status}`;
        statusIndicator.title = `WebSocket: ${status}`;
    }
}

function handleWebSocketMessage(data) {
    switch (data.event) {
        case 'task.created':
            addTaskToList(data.data);
            showNotification('New task created', 'success');
            break;
        case 'task.updated':
            updateTaskInList(data.data);
            break;
        case 'task.deleted':
            removeTaskFromList(data.data.id);
            showNotification('Task deleted', 'info');
            break;
        case 'task.status':
            updateTaskStatus(data.data);
            break;
        case 'task.progress':
            updateTaskProgress(data.data);
            break;
        default:
            console.log('Unknown event:', data.event);
    }
}

function addTaskToList(task) {
    const taskList = document.getElementById('taskList');
    if (!taskList) return;

    const taskElement = createTaskElement(task);
    taskList.insertBefore(taskElement, taskList.firstChild);
    updateTaskCount();
}

function updateTaskInList(task) {
    const existingTask = document.getElementById(`task-${task.id}`);
    if (!existingTask) return;

    const newTaskElement = createTaskElement(task);
    existingTask.replaceWith(newTaskElement);
}

function removeTaskFromList(taskId) {
    const taskElement = document.getElementById(`task-${taskId}`);
    if (taskElement) {
        taskElement.remove();
        updateTaskCount();
    }
}

function updateTaskStatus(data) {
    const taskElement = document.getElementById(`task-${data.id}`);
    if (!taskElement) return;

    const statusBadge = taskElement.querySelector('.status-badge');
    if (statusBadge) {
        statusBadge.className = `status-badge ${data.status.toLowerCase()}`;
        statusBadge.textContent = data.status;
    }

    if (data.status === 'completed') {
        showNotification(`Task ${data.id} completed successfully`, 'success');
    } else if (data.status === 'failed') {
        showNotification(`Task ${data.id} failed: ${data.error}`, 'error');
    }
}

function updateTaskProgress(data) {
    const taskElement = document.getElementById(`task-${data.id}`);
    if (!taskElement) return;

    const progressBar = taskElement.querySelector('.progress-bar');
    if (progressBar) {
        progressBar.style.width = `${data.percentage}%`;
        progressBar.setAttribute('aria-valuenow', data.percentage);
        progressBar.textContent = `${data.percentage}%`;
    }
}

function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;

    const container = document.getElementById('notification-container');
    if (container) {
        container.appendChild(notification);
        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => notification.remove(), 300);
        }, 3000);
    }
}

// Connect when the page loads
document.addEventListener('DOMContentLoaded', connectWebSocket); 