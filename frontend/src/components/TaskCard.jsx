import TaskList from './TaskList'

const TaskCard = ({
  onOpenCreate,
  tasks,
  loadingTasks,
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
  return (
    <section className="card">
      <div className="card-header">
        <h2>Görevlerin</h2>
        <button onClick={onOpenCreate}>Yeni Görev</button>
      </div>

      <TaskList
        tasks={tasks}
        loading={loadingTasks}
        editingTaskId={editingTaskId}
        editTitle={editTitle}
        editContent={editContent}
        editStatus={editStatus}
        onEditTitleChange={onEditTitleChange}
        onEditContentChange={onEditContentChange}
        onEditStatusChange={onEditStatusChange}
        onStartEdit={onStartEdit}
        onCancelEdit={onCancelEdit}
        onUpdate={onUpdate}
        onDelete={onDelete}
      />
    </section>
  )
}

export default TaskCard
