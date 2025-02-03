import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './registration_component.css';
import privates from "../../../assets/private.svg";
import look from "../../../assets/eye-look-icon.svg";

const Registration = ({ setIsAuthenticated }) => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [agree, setAgree] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      navigate('/auth/profile');
    }
  }, [navigate]);

  const validatePassword = (password) => {
    const minLength = 6;
    const hasUpperCase = /[A-Z]/.test(password);
    const hasLowerCase = /[a-z]/.test(password);
    const hasNumber = /[0-9]/.test(password);

    if (password.length < minLength) {
      return `Пароль должен быть не короче ${minLength} символов`;
    }
    if (!hasUpperCase) {
      return 'Пароль должен содержать хотя бы одну заглавную букву';
    }
    if (!hasLowerCase) {
      return 'Пароль должен содержать хотя бы одну строчную букву';
    }
    if (!hasNumber) {
      return 'Пароль должен содержать хотя бы одну цифру';
    }
    return '';
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(''); // Сброс ошибки перед новым запросом

    const passwordError = validatePassword(password);
    if (passwordError) {
      setError(passwordError);
      return;
    }

    if (password !== confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }

    if (!agree) {
      setError('Вы должны согласиться с обработкой данных');
      return;
    }

    const response = await fetch('http://localhost:5000/register', {
      method: 'POST',
      credentials: "include",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username,
        email,
        password,
      }),
    });

    if (response.ok) {
      alert('Регистрация успешна!');
      navigate('/auth/profile');
    } else {
      const errorData = await response.json();
      setError(errorData.message || 'Ошибка при регистрации');
    }
  };

  return (
    <div className="registration-component">
      <h1>Регистрация</h1>

      <form onSubmit={handleSubmit} className="registration-form">
        <div className='input'>
          <input
            type="text"
            placeholder="Логин"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <div className="password-input">
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Повторите пароль"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
            />
            <button type="button" onClick={() => setShowPassword(!showPassword)}>
              {showPassword ? (
                <img src={privates} alt="Скрыть" />
              ) : (
                <img src={look} alt="Показать" />
              )}
            </button>
          </div>
        </div>
      {error && <p className="error">{error}</p>}



        <div className="agreement">
          <input
            type="checkbox"
            checked={agree}
            onChange={(e) => setAgree(e.target.checked)}
            className="agree"
          />
          <label>Согласен с обработкой данных</label>
        </div>
        <div className='buttons'>
          <button type="submit" className='btn'>Зарегистрироваться</button>
          <Link to="authorization" className='btn'>Авторизация</Link>
        </div>
      </form>
    </div>
  );
};

export default Registration;