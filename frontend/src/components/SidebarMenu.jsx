const SidebarMenu = ({ isOpen, onClose, userName, onLogout }) => {
  return (
    <>
      <button
        className={`sidebar-overlay ${isOpen ? 'open' : ''}`}
        onClick={onClose}
        aria-label="Menüyü kapat"
      />
      <aside className={`sidebar-panel ${isOpen ? 'open' : ''}`}>
        <div className="sidebar-panel-header">
          <p className="sidebar-title">Profil</p>
          <button className="ghost" onClick={onClose}>
            Kapat
          </button>
        </div>
        <p className="sidebar-name">{userName}</p>
        <button className="ghost" onClick={onLogout}>
          Çıkış
        </button>
      </aside>
    </>
  )
}

export default SidebarMenu
