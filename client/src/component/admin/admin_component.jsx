import { Route, Routes, Link, useLocation, Navigate } from 'react-router-dom';
import './admin_component.css';
import AdminCheckProfile from '../admin_components/admin_check_profile/admin_check_profile';
import AdminAcceptPage from '../admin_components/admin_accept_page/admin_accept_page';
import AdminCourseManagement from '../admin_components/admin_course_management/admin_course_management';
import AdminCoursesManagement from '../admin_components/admin_check_manage/admin_check_manage';
import AdminGroupManagement from '../admin_components/admin_group_management/admin_group_management';
import { useEffect, useState } from 'react';

const AdminComponent = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Проверяем текущий путь
  const location = useLocation();

  // Функция для проверки, активна ли ссылка
  const isActive = (path) => {
    return location.pathname === path;
  };

  // Загрузка профиля пользователя
  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch("http://localhost:5000/profile", {
          method: "GET",
          credentials: "include",
        });

        if (!response.ok) {
          throw new Error("Ошибка при загрузке профиля");
        }

        const data = await response.json();
        setProfile(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  // Защита маршрута: Проверяем, является ли пользователь администратором
  if (loading) {
    return <p>Загрузка...</p>;
  }

  if (error) {
    return <p>{error}</p>;
  }

  if (!profile || profile.username !== 'admin') {
    return <Navigate to="/" replace />; // Перенаправление на страницу входа
  }

  return (
    <div className="container admin-component">
      <div className="admin-component-top">
        <div>
          <h1 className="admin-title">Административная панель</h1>
        </div>
        <nav>
          <ul className="nav-links">
            <li>
              <Link
                to="/admin"
                className={`nav-link ${isActive('/admin') ? 'active' : ''}`}
              >
                Управление курсами
              </Link>
            </li>
            <li>
              <Link
                to="/admin/profile"
                className={`nav-link ${isActive('/admin/profile') ? 'active' : ''}`}
              >
                Просмотр всех профилей
              </Link>
            </li>
            <li>
              <Link
                to="/admin/courses"
                className={`nav-link ${isActive('/admin/courses') ? 'active' : ''}`}
              >
                Управление заявками
              </Link>
            </li>
            <li>
              <Link
                to="/admin/groups"
                className={`nav-link ${isActive('/admin/groups') ? 'active' : ''}`}
              >
                Управление группами
              </Link>
            </li>
          </ul>
        </nav>
      </div>
      <div className="admin_components">
        <Routes>
          <Route index element={<AdminCourseManagement />} />
          <Route path="profile" element={<AdminCheckProfile />} />
          <Route path="courses" element={<AdminCoursesManagement />} />
          <Route path="accept" element={<AdminAcceptPage />} />
          <Route path="groups" element={<AdminGroupManagement />} />
        </Routes>
      </div>
    </div>
  );
};

export default AdminComponent;