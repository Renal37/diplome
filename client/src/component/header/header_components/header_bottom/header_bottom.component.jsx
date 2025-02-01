import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './header_bottom.component.css';
import Cookies from 'js-cookie';

const Header_bottom = ({ isAuthenticated, setIsAuthenticated }) => {
    const navigate = useNavigate();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);

    useEffect(() => {
        if (isAuthenticated) {
            // Получаем имя пользователя из localStorage или из другого источника
            const storedUsername = localStorage.getItem('username');
            if (storedUsername) {
                setUsername(storedUsername);
            }
        }
    }, [isAuthenticated]);

    const handleLogout = async () => {
        try {
            const response = await fetch('http://localhost:5000/logout', {
                method: 'POST',
                credentials: 'include',
            });

            if (response.ok) {
                Cookies.remove("token", { path: "/" });
                localStorage.clear();
                setIsAuthenticated(false);
                navigate('/');
            } else {
                console.error('Ошибка при выходе из системы');
            }
        } catch (error) {
            console.error('Ошибка при выходе из системы:', error);
        }
    };

    const toggleDropdown = () => {
        setIsDropdownOpen(!isDropdownOpen);
    };

    return (
        <div className="header_bottom">
            <nav>
                <ul>
                    <li><Link to="/professional" className="a">Профессиональная переподготовка</Link></li>
                    <li><Link to="/promotion" className="a">Повышение квалификации</Link></li>
                    {isAuthenticated ? (
                        <li className="dropdown">
                            <Link onClick={toggleDropdown} id="drop" className="a">
                                Мои данные
                            </Link>
                            {isDropdownOpen && (
                                <div className="dropdown-menu">
                                    <Link to="/auth/profile" className="dropdown-item">Профиль</Link>
                                    <Link to="/auth/edit_profile" className="dropdown-item">Редактировать профиль</Link>
                                    <button onClick={handleLogout} className="dropdown-item">Выйти</button>
                                </div>
                            )}
                        </li>
                    ) : (
                        <li><Link to="/auth" className="a">Войти</Link></li>
                    )}
                </ul>
            </nav>
        </div>
    );
};

export default Header_bottom;