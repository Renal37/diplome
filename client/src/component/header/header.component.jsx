import { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import Header_top from './header_components/header_top/header_top.component';
import Header_center from './header_components/header_center/header_center.component';
import Header_bottom from './header_components/header_bottom/header_bottom.component';
import { Fragment } from 'react';
import './header.component.css';

const Header = () => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            checkToken(token);
        }

        // Добавляем слушатель для изменения токена в localStorage
        const storageListener = () => {
            const token = localStorage.getItem('token');
            if (token) {
                checkToken(token);
            } else {
                setIsAuthenticated(false);
            }
        };

        window.addEventListener('storage', storageListener);

        return () => {
            window.removeEventListener('storage', storageListener);
        };
    }, []);

    const checkToken = async (token) => {
        const response = await fetch('http://localhost:5000/check-token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            credentials: 'include',
        });

        if (response.ok) {
            console.log('Токен действителен');
            setIsAuthenticated(true);
        } else {
            console.log('Неверный токен');
            setIsAuthenticated(false);
        }
    };

    return (
        <Fragment>
            <header>
                <div className="header container">
                    <Header_top />
                    <Header_center />
                </div>
                <Header_bottom isAuthenticated={isAuthenticated} setIsAuthenticated={setIsAuthenticated} />
            </header>
            <Outlet />
        </Fragment>
    );
};

export default Header;
