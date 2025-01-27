import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './authorization_component.css';

const Authorization = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();
    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            navigate('/auth/profile');
        }
    }, [navigate]);
    const handleSubmit = async (e) => {
        e.preventDefault();
        const response = await fetch('http://localhost:5000/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include', // Добавьте это для отправки куки
            body: JSON.stringify({
                username,
                password,
            }),
        });

        if (response.ok) {
            const data = await response.json();
            localStorage.setItem('token', data.token);
            alert('Авторизация успешна!');
            navigate('/auth/profile');
        } else {
            alert('Ошибка авторизации');
        }
    };

    return (
        <div className="authorization-component">
            <h1>Авторизация</h1>
            <Link to="/auth">Регистрация</Link>
            <form onSubmit={handleSubmit} className="authorization-form">
                <input
                    type="text"
                    placeholder="Логин"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Пароль"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button type="submit">Войти</button>
            </form>
        </div>
    );
};

export default Authorization;