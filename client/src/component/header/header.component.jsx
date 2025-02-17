import { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import Header_top from './header_components/header_top/header_top.component';
import Header_center from './header_components/header_center/header_center.component';
import Header_bottom from './header_components/header_bottom/header_bottom.component';
import { Fragment } from 'react';
import './header.component.css';

const Header = () => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    const checkToken = async () => {
        try {
            const response = await fetch('http://localhost:5000/check-token', {
                method: 'POST',
                credentials: 'include', // Куки отправляются автоматически
            });
            if (response.ok) {
                setIsAuthenticated(true);
            } else {
                setIsAuthenticated(false);
            }
        } catch (error) {
            console.error('Ошибка при проверке токена:', error);
            setIsAuthenticated(false);
        }
    };

    useEffect(() => {
        checkToken();
        const interval = setInterval(() => {
            checkToken();
        }, 1 * 1000); // Проверяем токен каждую минуту
        return () => clearInterval(interval);
    }, []);

    return (
        <Fragment>
            <header>
                <div className="header container">
                    <Header_top />
                    <Header_center />
                </div>
                <Header_bottom
                    isAuthenticated={isAuthenticated}
                    setIsAuthenticated={setIsAuthenticated}
                    checkToken={checkToken} // Передаем функцию checkToken
                />
            </header>
            <Outlet />
        </Fragment>
    );
};

export default Header;