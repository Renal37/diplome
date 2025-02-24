import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import './header_bottom.component.css';
import Cookies from 'js-cookie';

const Header_bottom = ({ isAuthenticated, setIsAuthenticated }) => {
    const navigate = useNavigate();
    const location = useLocation(); // Получаем текущий путь
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
            if (window.scrollY > 350) {
                setIsScrolled(true);
                // Добавляем отступ для контента под header_bottom
                document.querySelector('.main_top').style.paddingTop = '170px'; // Замените на актуальную высоту
            } else {
                setIsScrolled(false);
                // Убираем отступ, если header_bottom не фиксированный
                document.querySelector('.main_top').style.paddingTop = '0';
            }
        };

        // Добавляем слушатель события scroll
        window.addEventListener('scroll', handleScroll);

        // Удаляем слушатель при размонтировании компонента
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [isAuthenticated]);

    useEffect(() => {
        // Устанавливаем атрибут data-path для body
        document.body.setAttribute('data-path', location.pathname);
    }, [location.pathname]);

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

    // Определяем класс для header_bottom с учетом текущего пути
    const headerBottomClass = location.pathname.startsWith('/admin') || location.pathname.startsWith('/auth')
        ? 'header_bottom' // Без класса "scrolled" на страницах /admin и /auth
        : `header_bottom ${isScrolled ? 'scrolled' : ''}`;

    return (
        <div className={headerBottomClass}>
            <div className="header_bottom_bg">
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
            {/* Проверяем, чтобы путь не начинался с /admin или /auth */}
            {!(location.pathname.startsWith('/admin') || location.pathname.startsWith('/auth')) && (
                <div className='request_course'>
                    <div className="request_card">
                        <div className="request_item">
                            <div className="request_text">
                                <h1>Записаться на интересный курс:</h1>
                                <p>Педагог профессионального обучения</p>
                            </div>
                            <Link to="/courses/register/678df19a0a96fa5b989aeaa5" className="request_button">Записаться</Link>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Header_bottom;