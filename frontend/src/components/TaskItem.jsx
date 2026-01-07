const TaskItem = ({
  task,
  isEditing,
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
  const handleDragStart = (event) => {
    event.dataTransfer.setData('text/plain', String(task.id))
  }

  if (isEditing) {
    return (
      <li>
        <form className="form edit" onSubmit={onUpdate}>
          <input value={editTitle} onChange={(event) => onEditTitleChange(event.target.value)} />
          <textarea
            rows={3}
            value={editContent}
            onChange={(event) => onEditContentChange(event.target.value)}
          />
          <select value={editStatus} onChange={(event) => onEditStatusChange(event.target.value)}>
            <option value="todo">Yapılacak</option>
            <option value="in_progress">Devam ediyor</option>
            <option value="done">Tamamlandı</option>
          </select>
          <div className="actions">
            <button type="submit">Kaydet</button>
            <button type="button" className="ghost" onClick={onCancelEdit}>
              Vazgeç
            </button>
          </div>
        </form>
      </li>
    )
  }

  return (
    <li draggable={!isEditing} onDragStart={handleDragStart}>
      <div className="task-content">
        <span className="task-title">{task.title}</span>
        <span className="task-status">{task.status ?? 'todo'}</span>
        {task.content && <p>{task.content}</p>}
      </div>
      <div className="actions">
        <button className="ghost" onClick={onStartEdit}>
          Düzenle
        </button>
        <button className="link" onClick={onDelete}>
          Sil
        </button>
      </div>
    </li>
  )
}

export default TaskItem
