import { Route, Routes, Link, useLocation } from 'react-router-dom';
import './admin_component.css';
import AdminCheckProfile from '../admin_components/admin_check_profile/admin_check_profile';
import AdminAcceptPage from '../admin_components/admin_accept_page/admin_accept_page';
import AdminCourseManagement from '../admin_components/admin_course_management/admin_course_management';
import AdminCoursesManagement from '../admin_components/admin_check_manage/admin_check_manage';
import AdminGroupManagement from '../admin_components/admin_group_management/admin_group_management';

const AdminComponent = () => {
  const location = useLocation(); // Получаем текущий путь

  // Функция для проверки, активна ли ссылка
  const isActive = (path) => {
    return location.pathname === path;
  };

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

            {/* <li>
              <Link to="/admin/accept" className="nav-link">Выдача докумета</Link>
            </li> */}
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