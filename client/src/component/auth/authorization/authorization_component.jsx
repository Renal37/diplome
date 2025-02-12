import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './authorization_component.css';
import privates from "../../../assets/private.svg";
import look from "../../../assets/eye-look-icon.svg";

const AuthorizationComponent = ({ setIsAuthenticated }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        const response = await fetch('http://localhost:5000/login', {
            method: 'POST',
            credentials: "include",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        if (response.ok) {
            alert('Авторизация успешна!');
            navigate('/auth/profile');
        } else {
            alert('Ошибка авторизации');
        }
    };

    return (
        <div className="authorization-component">
            <h1>Авторизация</h1>
            <form onSubmit={handleSubmit} className="authorization-form">
                <div className="input">
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
                        <button type="button" onClick={() => setShowPassword(!showPassword)}>
                            {showPassword ? (
                                <img src={privates} alt="Скрыть" />
                            ) : (
                                <img src={look} alt="Показать" />
                            )}
                        </button>
                    </div>
                </div>
                <div className="btns">
                    <Link to="/auth" className='btn'>Регистрация</Link>
                    <button type="submit" className='btn'>Войти</button>
                </div>
            </form>
        </div>
    );
};

export default AuthorizationComponent;