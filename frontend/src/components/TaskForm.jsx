const TaskForm = ({
  className = 'form inline',
  title,
  content,
  status,
  onTitleChange,
  onContentChange,
  onStatusChange,
  onSubmit,
}) => {
  return (
    <form className={className} onSubmit={onSubmit}>
      <input
        placeholder="Yeni görev başlığı"
        value={title}
        onChange={(event) => onTitleChange(event.target.value)}
      />
      <textarea
        placeholder="Kısa açıklama"
        rows={3}
        value={content}
        onChange={(event) => onContentChange(event.target.value)}
      />
      <select value={status} onChange={(event) => onStatusChange(event.target.value)}>
        <option value="todo">Yapılacak</option>
        <option value="in_progress">Devam ediyor</option>
        <option value="done">Tamamlandı</option>
      </select>
      <button type="submit">Ekle</button>
    </form>
  )
}

export default TaskForm
