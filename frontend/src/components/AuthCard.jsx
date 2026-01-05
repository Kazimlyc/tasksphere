const AuthCard = ({ authMode, onAuthModeChange, authForm, onAuthFormChange, onSubmit }) => {
  return (
    <section className="card">
      <div className="tabs">
        <button className={authMode === 'login' ? 'active' : ''} onClick={() => onAuthModeChange('login')}>
          Giriş
        </button>
        <button
          className={authMode === 'register' ? 'active' : ''}
          onClick={() => onAuthModeChange('register')}
        >
          Kayıt
        </button>
      </div>

      <form className="form" onSubmit={onSubmit}>
        {authMode === 'register' && (
          <label>
            İsim
            <input
              type="text"
              value={authForm.name}
              required
              onChange={(event) => onAuthFormChange({ ...authForm, name: event.target.value })}
            />
          </label>
        )}
        <label>
          E-posta
          <input
            type="email"
            value={authForm.email}
            required
            onChange={(event) => onAuthFormChange({ ...authForm, email: event.target.value })}
          />
        </label>
        <label>
          Şifre
          <input
            type="password"
            value={authForm.password}
            required
            minLength={6}
            onChange={(event) => onAuthFormChange({ ...authForm, password: event.target.value })}
          />
        </label>

        <button type="submit">{authMode === 'login' ? 'Giriş yap' : 'Kayıt ol'}</button>
      </form>
    </section>
  )
}

export default AuthCard
