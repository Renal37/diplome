import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './registration_component.css';

const Registration = () => {
  const [fullName, setFullName] = useState('');
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [birthDate, setBirthDate] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch('http://localhost:5000/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        fullName,
        username,
        email,
        password,
        birthDate,
      }),
    });

    if (response.ok) {
      alert('Регистрация успешна!');
      navigate('/profile');
    } else {
      alert('Ошибка при регистрации');
    }
  };

  return (
    <div className="registration-component">
      <h1>Регистрация</h1>
      <form onSubmit={handleSubmit} className="registration-form">
        <input
          type="text"
          placeholder="ФИО"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Логин"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="email"
          placeholder="Почта"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <input
          type="date"
          placeholder="Дата рождения"
          value={birthDate}
          onChange={(e) => setBirthDate(e.target.value)}
        />
        <button type="submit">Зарегистрироваться</button>
      </form>
    </div>
  );
};

export default Registration;