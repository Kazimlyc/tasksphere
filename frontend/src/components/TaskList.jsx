import TaskItem from './TaskItem'

const statusLabels = {
  todo: 'Yapılacak',
  in_progress: 'Devam ediyor',
  done: 'Tamamlandı',
}

const TaskList = ({
  tasks,
  loading,
  editingTaskId,
  editTitle,
  editContent,
  editStatus,
  onEditTitleChange,
  onEditContentChange,
  onEditStatusChange,
  onStartEdit,
  onCancelEdit,
  onUpdate,
  onDelete,
}) => {
  if (loading) {
    return <p>Yükleniyor...</p>
  }

  const renderColumn = (statusKey) => {
    const columnTasks = tasks.filter((task) => (task.status ?? 'todo') === statusKey)
    return (
      <div className="task-column" key={statusKey}>
        <div className="task-column-header">
          <h3>{statusLabels[statusKey]}</h3>
          <span>{columnTasks.length}</span>
        </div>
        {columnTasks.length === 0 ? (
          <p className="task-empty">Bu listede görev yok</p>
        ) : (
          <ul className="task-list">
            {columnTasks.map((task) => (
              <TaskItem
                key={task.id}
                task={task}
                isEditing={editingTaskId === task.id}
                editTitle={editTitle}
                editContent={editContent}
                editStatus={editStatus}
                onEditTitleChange={onEditTitleChange}
                onEditContentChange={onEditContentChange}
                onEditStatusChange={onEditStatusChange}
                onStartEdit={() => onStartEdit(task)}
                onCancelEdit={onCancelEdit}
                onUpdate={onUpdate}
                onDelete={() => onDelete(task.id)}
              />
            ))}
          </ul>
        )}
      </div>
    )
  }

  if (tasks.length === 0) {
    return <p>Henüz görev yok</p>
  }

  return (
    <div className="task-columns">
      {renderColumn('todo')}
      {renderColumn('in_progress')}
      {renderColumn('done')}
    </div>
  )
}

export default TaskList
