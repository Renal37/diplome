import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './authorization_component.css';

const AuthorizationComponent = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        const response = await fetch('http://localhost:5000/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
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
                <div className="password-input">
                    <input
                        type={showPassword ? "text" : "password"}
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                    >
                        {showPassword ? "Скрыть" : "Показать"}
                    </button>
                </div>
                <button type="submit">Войти</button>
            </form>
        </div>
    );
};

export default AuthorizationComponent;