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
  onCreateForStatus,
  onMoveTask,
}) => {
  return (
    <section className="card">
      <div className="card-header">
        <h2>Görevlerin</h2>
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
        onCreateForStatus={onCreateForStatus}
        onMoveTask={onMoveTask}
      />
    </section>
  )
}

export default TaskCard
