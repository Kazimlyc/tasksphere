import TaskForm from './TaskForm'

const CreateTaskModal = ({
  isOpen,
  onClose,
  title,
  content,
  status,
  onTitleChange,
  onContentChange,
  onStatusChange,
  onSubmit,
}) => {
  return (
    <>
      <button
        className={`modal-overlay ${isOpen ? 'open' : ''}`}
        onClick={onClose}
        aria-label="Yeni görev ekranını kapat"
      />
      <div className={`modal-panel ${isOpen ? 'open' : ''}`}>
        <div className="modal-header">
          <h3>Yeni Görev</h3>
          <button className="ghost" onClick={onClose}>
            Kapat
          </button>
        </div>
        <TaskForm
          className="form"
          title={title}
          content={content}
          status={status}
          onTitleChange={onTitleChange}
          onContentChange={onContentChange}
          onStatusChange={onStatusChange}
          onSubmit={onSubmit}
        />
      </div>
    </>
  )
}

export default CreateTaskModal
