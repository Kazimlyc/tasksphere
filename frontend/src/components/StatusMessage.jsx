const StatusMessage = ({ status }) => {
  if (!status) return null
  return <div className="status">{status}</div>
}

export default StatusMessage
