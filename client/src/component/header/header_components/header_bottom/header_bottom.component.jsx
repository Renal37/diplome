import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './header_bottom.component.css';
import Cookies from 'js-cookie';

const Header_bottom = ({ isAuthenticated, setIsAuthenticated }) => {
    const navigate = useNavigate();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [isScrolled, setIsScrolled] = useState(false); // Новое состояние для отслеживания скролла

    useEffect(() => {
        if (isAuthenticated) {
            const storedUsername = localStorage.getItem('username');
            if (storedUsername) {
                setUsername(storedUsername);
            }
        }

        // Функция для отслеживания скролла
        const handleScroll = () => {
            if (window.scrollY > 220) { // Если прокручено больше 100px
                setIsScrolled(true);
            } else {
                setIsScrolled(false);
            }
        };

        // Добавляем слушатель события scroll
        window.addEventListener('scroll', handleScroll);

        // Удаляем слушатель при размонтировании компонента
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
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
        <div className={`header_bottom ${isScrolled ? 'scrolled' : ''}`}>
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
            <div className='request_course'>
                <div className="request_card">
                    <div className="request_item">
                        <div className="request_text">
                            <h1>Записаться на интересный курс: </h1>
                            <p>Педагог профессионального обучения</p>
                        </div>

                        <Link to="/courses/register/678df19a0a96fa5b989aeaa5" className="request_button">Записаться</Link>
                    </div>

                </div>

            </div>
        </div>
    );
};

export default Header_bottom;